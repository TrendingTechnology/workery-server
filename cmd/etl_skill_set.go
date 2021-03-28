package cmd

import (
	"context"
	"fmt"
	"os"
	"log"
	"database/sql"
	"time"

	"github.com/spf13/cobra"
	"github.com/google/uuid"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	skillSetETLSchemaName string
	skillSetETLTenantId int
)

func init() {
	skillSetETLCmd.Flags().StringVarP(&skillSetETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	skillSetETLCmd.MarkFlagRequired("schema_name")
	skillSetETLCmd.Flags().IntVarP(&skillSetETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	skillSetETLCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(skillSetETLCmd)
}

var skillSetETLCmd = &cobra.Command{
	Use:   "etl_skill_set",
	Short: "Import the insurance requirement data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportSkillSet()
	},
}

func doRunImportSkillSet() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, skillSetETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

    // Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	ssr := repositories.NewSkillSetRepo(db)

	runSkillSetETL(ctx, uint64(skillSetETLTenantId), ssr, oldDb)
}

func runSkillSetETL(ctx context.Context, tenantId uint64, ssr *repositories.SkillSetRepo, oldDb *sql.DB) {
	skillSets, err := ListAllSkillSets(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oss := range skillSets {
		insertSkillSetETL(ctx, tenantId, ssr, oss)
	}
}

type OldSkillSet struct {
	Id                      uint64 `json:"id"`
	Category          string    `json:"category"`
	SubCategory       string    `json:"sub_category"`
	Description       string    `json:"description"`
	IsArchived        bool      `json:"is_archived"`
    OldId             uint64    `json:"old_id"`
}

func ListAllSkillSets(db *sql.DB) ([]*OldSkillSet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, category, sub_category, description, is_archived
	FROM
        workery_skill_sets
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldSkillSet
	defer rows.Close()
	for rows.Next() {
		m := new(OldSkillSet)
		err = rows.Scan(
			&m.Id,
			&m.Category,
			&m.SubCategory,
			&m.Description,
			&m.IsArchived,
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

func insertSkillSetETL(ctx context.Context, tid uint64, ssr *repositories.SkillSetRepo, oss *OldSkillSet) {
	var state int8 = 1
	if oss.IsArchived == true {
		state = 0
	}

	m := &models.SkillSet{
		OldId: oss.Id,
		TenantId: tid,
		Uuid: uuid.NewString(),
		Category: oss.Category,
		SubCategory: oss.SubCategory,
		Description: oss.Description,
		State: state,
	}
	err := ssr.Insert(ctx, m)
	if err != nil {
		log.Panic("ssr.Insert | err", err)
	}
	fmt.Println("Imported ID#", oss.Id)
}
