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
	assETLSchemaName string
	assETLTenantId   int
)

func init() {
	assETLCmd.Flags().StringVarP(&assETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	assETLCmd.MarkFlagRequired("schema_name")
	assETLCmd.Flags().IntVarP(&assETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	assETLCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(assETLCmd)
}

var assETLCmd = &cobra.Command{
	Use:   "etl_associate_skill_set",
	Short: "Import the associate vehicle types from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportAssociateSkillSet()
	},
}

func doRunImportAssociateSkillSet() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, assETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	assr := repositories.NewAssociateSkillSetRepo(db)
	ar := repositories.NewAssociateRepo(db)
	ssr := repositories.NewSkillSetRepo(db)

	runAssociateSkillSetETL(ctx, uint64(assETLTenantId), assr, ar, ssr, oldDb)
}

func runAssociateSkillSetETL(
	ctx context.Context,
	tenantId uint64,
	assr *repositories.AssociateSkillSetRepo,
	ar *repositories.AssociateRepo,
	ssr *repositories.SkillSetRepo,
	oldDb *sql.DB,
) {
	asss, err := ListAllAssociateSkillSets(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oss := range asss {
		insertAssociateSkillSetETL(ctx, tenantId, assr, ar, ssr, oss)
	}
}

type OldAssociateSkillSet struct {
	Id          uint64 `json:"id"`
	AssociateId uint64 `json:"associate_id"`
	SkillSetId  uint64 `json:"skillset_id"`
}

func ListAllAssociateSkillSets(db *sql.DB) ([]*OldAssociateSkillSet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, associate_id, skillset_id
	FROM
        workery_associates_skill_sets
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldAssociateSkillSet
	defer rows.Close()
	for rows.Next() {
		m := new(OldAssociateSkillSet)
		err = rows.Scan(
			&m.Id,
			&m.AssociateId,
			&m.SkillSetId,
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

func insertAssociateSkillSetETL(
	ctx context.Context,
	tid uint64,
	assr *repositories.AssociateSkillSetRepo,
	ar *repositories.AssociateRepo,
	ssr *repositories.SkillSetRepo,
	oss *OldAssociateSkillSet,
) {
	associateId, err := ar.GetIdByOldId(ctx, tid, oss.AssociateId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	skillSetId, err := ssr.GetIdByOldId(ctx, tid, oss.SkillSetId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	m := &models.AssociateSkillSet{
		OldId:       oss.Id,
		TenantId:    tid,
		Uuid:        uuid.NewString(),
		AssociateId: associateId,
		SkillSetId:  skillSetId,
	}
	err = assr.Insert(ctx, m)
	if err != nil {
		log.Panic("assr.Insert | err", err)
	}
	fmt.Println("Imported ID#", oss.Id)
}
