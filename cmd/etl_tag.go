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
	tagETLSchemaName string
)

func init() {
	tagETLCmd.Flags().StringVarP(&tagETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	tagETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(tagETLCmd)
}

var tagETLCmd = &cobra.Command{
	Use:   "etl_tag",
	Short: "Import the tag data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportTag()
	},
}

func doRunImportTag() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, tagETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	irr := repositories.NewTagRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, tagETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runTagETL(ctx, tenant.Id, irr, oldDb)
}

type OldUTag struct {
	Id          uint64 `json:"id"`
	Text        string `json:"text"`
	Description string `json:"description"`
	IsArchived  bool   `json:"is_archived"`
}

func ListAllTags(db *sql.DB) ([]*OldUTag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, text, description, is_archived
	FROM
	    workery_tags
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUTag
	defer rows.Close()
	for rows.Next() {
		m := new(OldUTag)
		err = rows.Scan(
			&m.Id,
			&m.Text,
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

func runTagETL(ctx context.Context, tenantId uint64, irr *repositories.TagRepo, oldDb *sql.DB) {
	tags, err := ListAllTags(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oir := range tags {
		insertTagETL(ctx, tenantId, irr, oir)
	}
}

func insertTagETL(ctx context.Context, tid uint64, irr *repositories.TagRepo, oir *OldUTag) {
	var state int8 = 1
	if oir.IsArchived == true {
		state = 0
	}

	m := &models.Tag{
		OldId:       oir.Id,
		TenantId:    tid,
		Uuid:        uuid.NewString(),
		Text:        oir.Text,
		Description: oir.Description,
		State:       state,
	}
	err := irr.Insert(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", oir.Id)
}
