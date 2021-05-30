package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	customerCommentETLSchemaName string
)

func init() {
	customerCommentETLCmd.Flags().StringVarP(&customerCommentETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	customerCommentETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(customerCommentETLCmd)
}

var customerCommentETLCmd = &cobra.Command{
	Use:   "etl_customer_comment",
	Short: "Import the customerComment data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportCustomerComment()
	},
}

func doRunImportCustomerComment() {
	// Load up our NEW database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load up our OLD database.
	oldDBHost := os.Getenv("WORKERY_OLD_DB_HOST")
	oldDBPort := os.Getenv("WORKERY_OLD_DB_PORT")
	oldDBUser := os.Getenv("WORKERY_OLD_DB_USER")
	oldDBPassword := os.Getenv("WORKERY_OLD_DB_PASSWORD")
	oldDBName := os.Getenv("WORKERY_OLD_DB_NAME")
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, customerCommentETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	irr := repositories.NewCustomerCommentRepo(db)
	om := repositories.NewCustomerRepo(db)
	commentr := repositories.NewCommentRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, customerCommentETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runCustomerCommentETL(ctx, tenant.Id, irr, om, commentr, oldDb)
}

type OldUCustomerComment struct {
	Id         uint64    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	CustomerId uint64    `json:"customer_id"`
	CommentId  uint64    `json:"comment_id"`
}

func ListAllCustomerComments(db *sql.DB) ([]*OldUCustomerComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, created_at, about_id, comment_id
	FROM
	    workery_customer_comments
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUCustomerComment
	defer rows.Close()
	for rows.Next() {
		m := new(OldUCustomerComment)
		err = rows.Scan(
			&m.Id,
			&m.CreatedAt,
			&m.CustomerId,
			&m.CommentId,
		)
		if err != nil {
			log.Panic("ListAllCustomerComments | Next | err", err)
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return arr, err
}

func runCustomerCommentETL(ctx context.Context, tenantId uint64, irr *repositories.CustomerCommentRepo, om *repositories.CustomerRepo, tr *repositories.CommentRepo, oldDb *sql.DB) {
	customerComments, err := ListAllCustomerComments(oldDb)
	if err != nil {
		log.Fatal("runCustomerCommentETL | ListAllCustomerComments | err", err)
	}
	for _, oir := range customerComments {

		customerId, err := om.GetIdByOldId(ctx, tenantId, oir.CustomerId)
		if err != nil {
			log.Panic("runCustomerCommentETL | om.GetIdByOldId | err", err)
		}

		commentId, err := tr.GetIdByOldId(ctx, tenantId, oir.CommentId)
		if err != nil {
			log.Panic("runCustomerCommentETL | tr.GetIdByOldId | err", err)
		}

		insertCustomerCommentETL(ctx, tenantId, oir.Id, customerId, commentId, irr)
	}
}

func insertCustomerCommentETL(ctx context.Context, tenantId uint64, oldId uint64, customerId uint64, commentId uint64, irr *repositories.CustomerCommentRepo) {
	fmt.Println("Pre-Imported Customer Comment ID#", oldId)
	m := &models.CustomerComment{
		OldId:      oldId,
		Uuid:       uuid.NewString(),
		TenantId:   tenantId,
		CustomerId: customerId,
		CommentId:  commentId,
	}
	err := irr.Insert(ctx, m)
	if err != nil {
		log.Panic("insertCustomerCommentETL | Insert | err:", err)
	}
	fmt.Println("Post-Imported Customer Comment ID#", oldId)
}
