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
	commentETLSchemaName string
)

func init() {
	commentETLCmd.Flags().StringVarP(&commentETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	commentETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(commentETLCmd)
}

var commentETLCmd = &cobra.Command{
	Use:   "etl_comment",
	Short: "Import the comment data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportComment()
	},
}

func doRunImportComment() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, commentETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	irr := repositories.NewCommentRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, commentETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runCommentETL(ctx, tenant.Id, irr, oldDb)
}

type OldComment struct {
	Id               uint64         `json:"id"`
	CreatedAt        time.Time      `json:"created_time"`
	CreatedById      sql.NullInt64  `json:"created_by_id,omitempty"`
	CreatedFrom      sql.NullString `json:"created_from"`
	LastModifiedAt   time.Time      `json:"last_modified_time"`
	LastModifiedById sql.NullInt64  `json:"last_modified_by_id,omitempty"`
	LastModifiedFrom sql.NullString `json:"last_modified_from"`
	Text             string         `json:"text"`
	IsArchived       bool           `json:"is_archived"`
}

func ListAllComments(db *sql.DB) ([]*OldComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
		id, created_at, created_by_id, created_from, last_modified_at, last_modified_by_id, last_modified_from, text, is_archived
	FROM
	    workery_comments
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldComment
	defer rows.Close()
	for rows.Next() {
		m := new(OldComment)
		err = rows.Scan(
			&m.Id,
			&m.CreatedAt,
			&m.CreatedById,
			&m.CreatedFrom,
			&m.LastModifiedAt,
			&m.LastModifiedById,
			&m.LastModifiedFrom,
			&m.Text,
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

func runCommentETL(ctx context.Context, tenantId uint64, irr *repositories.CommentRepo, oldDb *sql.DB) {
	comments, err := ListAllComments(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oir := range comments {
		insertCommentETL(ctx, tenantId, irr, oir)
	}
}

func insertCommentETL(ctx context.Context, tid uint64, irr *repositories.CommentRepo, oir *OldComment) {
	var state int8 = 1
	if oir.IsArchived == true {
		state = 0
	}
	m := &models.Comment{
		OldId:              oir.Id,
		TenantId:           tid,
		Uuid:               uuid.NewString(),
		CreatedTime:        oir.CreatedAt,
		CreatedById:        oir.CreatedById,
		CreatedFromIP:      oir.CreatedFrom.String,
		LastModifiedTime:   oir.LastModifiedAt,
		LastModifiedById:   oir.LastModifiedById,
		LastModifiedFromIP: oir.LastModifiedFrom.String,
		Text:               oir.Text,
		State:              state,
	}
	err := irr.Insert(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", oir.Id)
}
