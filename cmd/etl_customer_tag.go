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
	customerTagETLSchemaName string
)

func init() {
	customerTagETLCmd.Flags().StringVarP(&customerTagETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	customerTagETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(customerTagETLCmd)
}

var customerTagETLCmd = &cobra.Command{
	Use:   "etl_customer_tag",
	Short: "Import the customerTag data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportCustomerTag()
	},
}

func doRunImportCustomerTag() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, customerTagETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	irr := repositories.NewCustomerTagRepo(db)
	om := repositories.NewCustomerRepo(db)
	tagr := repositories.NewTagRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, customerTagETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runCustomerTagETL(ctx, tenant.Id, irr, om, tagr, oldDb)
}

type OldUCustomerTag struct {
	Id          uint64 `json:"id"`
	CustomerId  uint64 `json:"customer_id"`
	TagId  uint64 `json:"tag_id"`
}

func ListAllCustomerTags(db *sql.DB) ([]*OldUCustomerTag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, customer_id, tag_id
	FROM
	    workery_customers_tags
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUCustomerTag
	defer rows.Close()
	for rows.Next() {
		m := new(OldUCustomerTag)
		err = rows.Scan(
			&m.Id,
			&m.CustomerId,
			&m.TagId,
		)
		if err != nil {
			log.Panic("ListAllCustomerTags | Next | err", err)
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return arr, err
}

func runCustomerTagETL(ctx context.Context, tenantId uint64, irr *repositories.CustomerTagRepo, om *repositories.CustomerRepo, tr *repositories.TagRepo, oldDb *sql.DB) {
	customerTags, err := ListAllCustomerTags(oldDb)
	if err != nil {
		log.Fatal("runCustomerTagETL | ListAllCustomerTags | err", err)
	}
	for _, oir := range customerTags {

		customerId, err := om.GetIdByOldId(ctx, tenantId, oir.CustomerId)
		if err != nil {
			log.Panic("runCustomerTagETL | om.GetIdByOldId | err", err)
		}

		tagId, err := tr.GetIdByOldId(ctx, tenantId, oir.TagId)
		if err != nil {
			log.Panic("runCustomerTagETL | tr.GetIdByOldId | err", err)
		}

		insertCustomerTagETL(ctx, tenantId, oir.Id, customerId, tagId, irr)
	}
}

func insertCustomerTagETL(ctx context.Context, tenantId uint64, oldId uint64, customerId uint64, tagId uint64, irr *repositories.CustomerTagRepo) {
    fmt.Println("Pre-Imported ID#", oldId)
	m := &models.CustomerTag{
		OldId:       oldId,
		Uuid:        uuid.NewString(),
		TenantId:    tenantId,
		CustomerId:  customerId,
		TagId:       tagId,
	}
	err := irr.Insert(ctx, m)
	if err != nil {
		log.Panic("insertCustomerTagETL | Insert | err:", err)
	}
	fmt.Println("Post-Imported ID#", oldId)
}
