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
	aatETLSchemaName string
	aatETLTenantId   int
)

func init() {
	aatETLCmd.Flags().StringVarP(&aatETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	aatETLCmd.MarkFlagRequired("schema_name")
	aatETLCmd.Flags().IntVarP(&aatETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	aatETLCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(aatETLCmd)
}

var aatETLCmd = &cobra.Command{
	Use:   "etl_associate_tag",
	Short: "Import the associate vehicle types from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportAssociateTag()
	},
}

func doRunImportAssociateTag() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, aatETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	aatr := repositories.NewAssociateTagRepo(db)
	ar := repositories.NewAssociateRepo(db)
	vtr := repositories.NewTagRepo(db)

	runAssociateTagETL(ctx, uint64(aatETLTenantId), aatr, ar, vtr, oldDb)
}

func runAssociateTagETL(
	ctx context.Context,
	tenantId uint64,
	aatr *repositories.AssociateTagRepo,
	ar *repositories.AssociateRepo,
	vtr *repositories.TagRepo,
	oldDb *sql.DB,
) {
	aats, err := ListAllAssociateTags(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oss := range aats {
		insertAssociateTagETL(ctx, tenantId, aatr, ar, vtr, oss)
	}
}

type OldAssociateTag struct {
	Id          uint64 `json:"id"`
	AssociateId uint64 `json:"associate_id"`
	TagId       uint64 `json:"vehicletype_id"`
}

func ListAllAssociateTags(db *sql.DB) ([]*OldAssociateTag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, associate_id, vehicletype_id
	FROM
        workery_associates_vehicle_types
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldAssociateTag
	defer rows.Close()
	for rows.Next() {
		m := new(OldAssociateTag)
		err = rows.Scan(
			&m.Id,
			&m.AssociateId,
			&m.TagId,
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

func insertAssociateTagETL(
	ctx context.Context,
	tid uint64,
	aatr *repositories.AssociateTagRepo,
	ar *repositories.AssociateRepo,
	vtr *repositories.TagRepo,
	oss *OldAssociateTag,
) {
	associateId, err := ar.GetIdByOldId(ctx, tid, oss.AssociateId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	tagId, err := vtr.GetIdByOldId(ctx, tid, oss.TagId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	m := &models.AssociateTag{
		OldId:       oss.Id,
		TenantId:    tid,
		Uuid:        uuid.NewString(),
		AssociateId: associateId,
		TagId:       tagId,
	}
	err = aatr.Insert(ctx, m)
	if err != nil {
		log.Panic("aatr.Insert | err", err)
	}
	fmt.Println("Imported ID#", oss.Id)
}
