package cmd

import (
	"bytes"
	"context"
	"database/sql"
	// "fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	// "github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	// null "gopkg.in/guregu/null.v4"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	privateFileUploadToTMPDIRETLSchemaName string
)

func init() {
	privateFileUploadToTMPDIRETLCmd.Flags().StringVarP(&privateFileUploadToTMPDIRETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	privateFileUploadToTMPDIRETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(privateFileUploadToTMPDIRETLCmd)
}

var privateFileUploadToTMPDIRETLCmd = &cobra.Command{
	Use:   "etl_private_file_upload_from_tmp_dir",
	Short: "Upload private files from the temporary directory to the new S3 bucket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunUploadPrivateFileFromTmpDir()
	},
}

// Special thanks via https://docs.digitalocean.com/products/spaces/resources/s3-sdk-examples/
func getS3ClientInstance() (*s3.S3, string) {
	key := os.Getenv("WORKERY_AWS_S3_ACCESS_KEY")
	secret := os.Getenv("WORKERY_AWS_S3_SECRET_KEY")
	endpoint := os.Getenv("WORKERY_AWS_S3_ENDPOINT")
	region := os.Getenv("WORKERY_AWS_S3_REGION")
	bucketName := os.Getenv("WORKERY_AWS_S3_BUCKET_NAME")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)
	return s3Client, bucketName
}

func ListAllPrivateFilesByTenantId(db *sql.DB, tenantId uint64) ([]*models.PrivateFile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, uuid, tenant_id, s3_key, title, description, indexed_text,
		created_time, created_by_id, created_from_ip, last_modified_time, last_modified_by_id,
	    last_modified_from_ip, associate_id, customer_id, partner_id, staff_id,
	    work_order_id, state, old_id
	FROM
	    private_files
	WHERE
	    tenant_id = $1 AND state = $2
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query, tenantId, 2) // 2=our unique case
	if err != nil {
		return nil, err
	}

	var arr []*models.PrivateFile
	defer rows.Close()
	for rows.Next() {
		m := new(models.PrivateFile)
		err = rows.Scan(
			&m.Id, &m.Uuid, &m.TenantId, &m.S3Key, &m.Title, &m.Description, &m.IndexedText,
			&m.CreatedTime, &m.CreatedById, &m.CreatedFromIP, &m.LastModifiedTime, &m.LastModifiedById,
			&m.LastModifiedFromIP, &m.AssociateId, &m.CustomerId, &m.PartnerId, &m.StaffId,
			&m.WorkOrderId, &m.State, &m.OldId,
		)
		if err != nil {
			panic(err)
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return arr, err
}

func doRunUploadPrivateFileFromTmpDir() {
	// Load up our NEW database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	pfr := repositories.NewPrivateFileRepo(db)

	// Load up our S3 instances
	s3Client, bucketName := getS3ClientInstance()

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, privateFileUploadToTMPDIRETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	privateFiles, err := ListAllPrivateFilesByTenantId(db, tenant.Id)
	if err != nil {
		log.Fatal(err)
	}
	for _, privateFile := range privateFiles {
		uploadPrivateFileToS3(pfr, s3Client, bucketName, privateFile)
	}
}

func uploadPrivateFileToS3(pfr *repositories.PrivateFileRepo, s3 *s3.S3, bucketName string, privateFile *models.PrivateFile) {
	// Generate a new key path for the S3 bucket to save the private file.
	// It's important that the following happens:
	// - Tenant ID must exists in key
	// - Private indiciation must exists
	// - Remove the UUID code from the filepath
	tenantIdStr := strconv.FormatUint(privateFile.TenantId, 10)
	filename := strings.ReplaceAll(privateFile.S3Key, "/tmp/", "")
	filename = strings.ReplaceAll(filename, privateFile.Uuid+"-", "")
	newS3Key := "tenant/" + tenantIdStr + "/private/uploads/" + filename

	// When we remove the local key and write it to S3, we need to check if
	// a file with a similar s3 does not exist, if it does then we need to
	// generate a new key.
	doesExist, err := pfr.CheckIfExistsByS3Key(context.Background(), newS3Key)
	if err != nil {
		log.Fatal("pfr.CheckIfExistsByS3Key:", err)
	}
	if doesExist {
		log.Println("Duplicate found! Appending UUID to file for private file ID:", privateFile.Id)
		newS3Key = "tenant/" + tenantIdStr + "/private/uploads/" + privateFile.Uuid + "-" + filename
	}

	// Open the file and read the content.
	f, err := os.Open(privateFile.S3Key)
	if err != nil {
		log.Fatal("os.Open:", err)
	}
	defer f.Close()

	// Read the contents of the file in byte[] format and convert to string.
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	contents := buf.String()

	// Upload the content to S3.
	err = utils.UploadBinToS3(s3, bucketName, newS3Key, contents, "private")
	if err != nil {
		log.Fatal("UploadBinToS3", err)
	}

	// Update the private file in the database.
	privateFile.State = 1
	privateFile.S3Key = newS3Key
	err = pfr.UpdateById(context.Background(), privateFile)
	if err != nil {
		log.Fatal("pfr.UpdateById:", err)
	}

	log.Println("Imported ID#", privateFile.Id)

	// log.Fatal("PROGRAMMER HALT") // For debugging purposes only.
}
