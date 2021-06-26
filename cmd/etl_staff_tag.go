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
	staffTagETLSchemaName string
)

func init() {
	staffTagETLCmd.Flags().StringVarP(&staffTagETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	staffTagETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(staffTagETLCmd)
}

var staffTagETLCmd = &cobra.Command{
	Use:   "etl_staff_tag",
	Short: "Import the staffTag data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportStaffTag()
	},
}

func doRunImportStaffTag() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, staffTagETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	irr := repositories.NewStaffTagRepo(db)
	om := repositories.NewStaffRepo(db)
	tagr := repositories.NewTagRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, staffTagETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runStaffTagETL(ctx, tenant.Id, irr, om, tagr, oldDb)
}

type OldUStaffTag struct {
	Id      uint64 `json:"id"`
	StaffId uint64 `json:"staff_id"`
	TagId   uint64 `json:"tag_id"`
}

func ListAllStaffTags(db *sql.DB) ([]*OldUStaffTag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, staff_id, tag_id
	FROM
	    workery_staff_tags
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUStaffTag
	defer rows.Close()
	for rows.Next() {
		m := new(OldUStaffTag)
		err = rows.Scan(
			&m.Id,
			&m.StaffId,
			&m.TagId,
		)
		if err != nil {
			log.Panic("ListAllStaffTags | Next | err", err)
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return arr, err
}

func runStaffTagETL(ctx context.Context, tenantId uint64, irr *repositories.StaffTagRepo, om *repositories.StaffRepo, tr *repositories.TagRepo, oldDb *sql.DB) {
	staffTags, err := ListAllStaffTags(oldDb)
	if err != nil {
		log.Fatal("runStaffTagETL | ListAllStaffTags | err", err)
	}
	for _, oir := range staffTags {

		staffId, err := om.GetIdByOldId(ctx, tenantId, oir.StaffId)
		if err != nil {
			log.Panic("runStaffTagETL | om.GetIdByOldId | err", err)
		}

		tagId, err := tr.GetIdByOldId(ctx, tenantId, oir.TagId)
		if err != nil {
			log.Panic("runStaffTagETL | tr.GetIdByOldId | err", err)
		}

		insertStaffTagETL(ctx, tenantId, oir.Id, staffId, tagId, irr)
	}
}

func insertStaffTagETL(ctx context.Context, tenantId uint64, oldId uint64, staffId uint64, tagId uint64, irr *repositories.StaffTagRepo) {
	fmt.Println("Pre-Imported Staff Tag ID#", oldId)
	m := &models.StaffTag{
		OldId:    oldId,
		Uuid:     uuid.NewString(),
		TenantId: tenantId,
		StaffId:  staffId,
		TagId:    tagId,
	}
	err := irr.Insert(ctx, m)
	if err != nil {
		log.Panic("insertStaffTagETL | Insert | err:", err)
	}
	fmt.Println("Post-Imported Staff Tag ID#", oldId)
}
