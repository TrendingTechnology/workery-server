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
	psETLSchemaName string
	psETLTenantId   int
)

func init() {
	psETLCmd.Flags().StringVarP(&psETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	psETLCmd.MarkFlagRequired("schema_name")
	psETLCmd.Flags().IntVarP(&psETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	psETLCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(psETLCmd)
}

var psETLCmd = &cobra.Command{
	Use:   "etl_partner_comment",
	Short: "Import the partner vehicle types from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportPartnerComment()
	},
}

func doRunImportPartnerComment() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, psETLSchemaName)
	if err != nil {
		log.Fatal("utils.ConnectDB", err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	asr := repositories.NewPartnerCommentRepo(db)
	ar := repositories.NewPartnerRepo(db)
	vtr := repositories.NewCommentRepo(db)

	runPartnerCommentETL(ctx, uint64(psETLTenantId), asr, ar, vtr, oldDb)
}

func runPartnerCommentETL(
	ctx context.Context,
	tenantId uint64,
	asr *repositories.PartnerCommentRepo,
	ar *repositories.PartnerRepo,
	vtr *repositories.CommentRepo,
	oldDb *sql.DB,
) {
	ass, err := ListAllPartnerComments(oldDb)
	if err != nil {
		log.Fatal("ListAllPartnerComments", err)
	}
	for _, oss := range ass {
		insertPartnerCommentETL(ctx, tenantId, asr, ar, vtr, oss)
	}
}

type OldPartnerComment struct {
	Id        uint64 `json:"id"`
	PartnerId uint64 `json:"about_id"`
	CommentId uint64 `json:"comment_id"`
}

func ListAllPartnerComments(db *sql.DB) ([]*OldPartnerComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, about_id, comment_id
	FROM
        workery_partner_comments
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldPartnerComment
	defer rows.Close()
	for rows.Next() {
		m := new(OldPartnerComment)
		err = rows.Scan(
			&m.Id,
			&m.PartnerId,
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

func insertPartnerCommentETL(
	ctx context.Context,
	tid uint64,
	asr *repositories.PartnerCommentRepo,
	ar *repositories.PartnerRepo,
	vtr *repositories.CommentRepo,
	oss *OldPartnerComment,
) {
	partnerId, err := ar.GetIdByOldId(ctx, tid, oss.PartnerId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	commentId, err := vtr.GetIdByOldId(ctx, tid, oss.CommentId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	if partnerId != 0 && commentId != 0 {
		m := &models.PartnerComment{
			OldId:     oss.Id,
			TenantId:  tid,
			Uuid:      uuid.NewString(),
			PartnerId: partnerId,
			CommentId: commentId,
		}
		err = asr.Insert(ctx, m)
		if err != nil {
			log.Panic("asr.Insert | err", err, "\n\n", m, oss)
		}
		fmt.Println("Imported ID#", oss.Id)
	} else {
		fmt.Println("-------------------\nSkipped ID#", oss.Id, "\n-------------------\npartnerId #", partnerId, "\ncommentId #", commentId, "\n\noss", oss)
	}
}
