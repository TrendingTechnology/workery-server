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
	asETLSchemaName string
	asETLTenantId   int
)

func init() {
	asETLCmd.Flags().StringVarP(&asETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	asETLCmd.MarkFlagRequired("schema_name")
	asETLCmd.Flags().IntVarP(&asETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	asETLCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(asETLCmd)
}

var asETLCmd = &cobra.Command{
	Use:   "etl_associate_comment",
	Short: "Import the associate vehicle types from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportAssociateComment()
	},
}

func doRunImportAssociateComment() {
	// Load up our NEW database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
		log.Fatal("utils.ConnectDB", err)
	}
	defer db.Close()

	// Load up our OLD database.
	oldDBHost := os.Getenv("WORKERY_OLD_DB_HOST")
	oldDBPort := os.Getenv("WORKERY_OLD_DB_PORT")
	oldDBUser := os.Getenv("WORKERY_OLD_DB_USER")
	oldDBPassword := os.Getenv("WORKERY_OLD_DB_PASSWORD")
	oldDBName := os.Getenv("WORKERY_OLD_DB_NAME")
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, asETLSchemaName)
	if err != nil {
		log.Fatal("utils.ConnectDB", err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	asr := repositories.NewAssociateCommentRepo(db)
	ar := repositories.NewAssociateRepo(db)
	vtr := repositories.NewCommentRepo(db)

	runAssociateCommentETL(ctx, uint64(asETLTenantId), asr, ar, vtr, oldDb)
}

func runAssociateCommentETL(
	ctx context.Context,
	tenantId uint64,
	asr *repositories.AssociateCommentRepo,
	ar *repositories.AssociateRepo,
	vtr *repositories.CommentRepo,
	oldDb *sql.DB,
) {
	ass, err := ListAllAssociateComments(oldDb)
	if err != nil {
		log.Fatal("ListAllAssociateComments", err)
	}
	for _, oss := range ass {
		insertAssociateCommentETL(ctx, tenantId, asr, ar, vtr, oss)
	}
}

type OldAssociateComment struct {
	Id          uint64 `json:"id"`
	AssociateId uint64 `json:"about_id"`
	CommentId   uint64 `json:"comment_id"`
}

func ListAllAssociateComments(db *sql.DB) ([]*OldAssociateComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, about_id, comment_id
	FROM
        workery_associate_comments
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldAssociateComment
	defer rows.Close()
	for rows.Next() {
		m := new(OldAssociateComment)
		err = rows.Scan(
			&m.Id,
			&m.AssociateId,
			&m.CommentId,
		)
		if err != nil {
			log.Fatal("rows.Scan", err)
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("rows.Err", err)
	}
	return arr, err
}

func insertAssociateCommentETL(
	ctx context.Context,
	tid uint64,
	asr *repositories.AssociateCommentRepo,
	ar *repositories.AssociateRepo,
	vtr *repositories.CommentRepo,
	oss *OldAssociateComment,
) {
	associateId, err := ar.GetIdByOldId(ctx, tid, oss.AssociateId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	commentId, err := vtr.GetIdByOldId(ctx, tid, oss.CommentId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	if associateId != 0 && commentId != 0 {
		m := &models.AssociateComment{
			OldId:       oss.Id,
			TenantId:    tid,
			Uuid:        uuid.NewString(),
			AssociateId: associateId,
			CommentId:   commentId,
		}
		err = asr.Insert(ctx, m)
		if err != nil {
			log.Panic("asr.Insert | err", err, "\n\n", m, oss)
		}
		fmt.Println("Imported ID#", oss.Id)
	} else {
		fmt.Println("-------------------\nSkipped ID#", oss.Id, "\n-------------------\nassociateId #", associateId, "\ncommentId #", commentId, "\n\noss", oss, "\n\n")
	}
}
