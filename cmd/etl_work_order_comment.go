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
	workOrderCommentSchemaName string
)

func init() {
	workOrderCommentCmd.Flags().StringVarP(&workOrderCommentSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	workOrderCommentCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(workOrderCommentCmd)
}

var workOrderCommentCmd = &cobra.Command{
	Use:   "etl_work_order_comment",
	Short: "Import the work order comments from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportWorkOrderComment()
	},
}

func doRunImportWorkOrderComment() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, workOrderCommentSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	wotp := repositories.NewWorkOrderCommentRepo(db)
	ar := repositories.NewWorkOrderRepo(db)
	vtr := repositories.NewCommentRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, workOrderCommentSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant != nil {
		runWorkOrderCommentETL(ctx, uint64(tenant.Id), wotp, ar, vtr, oldDb)
	}
}

func runWorkOrderCommentETL(
	ctx context.Context,
	tenantId uint64,
	wotp *repositories.WorkOrderCommentRepo,
	ar *repositories.WorkOrderRepo,
	vtr *repositories.CommentRepo,
	oldDb *sql.DB,
) {
	aats, err := ListAllWorkOrderComments(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oss := range aats {
		insertWorkOrderCommentETL(ctx, tenantId, wotp, ar, vtr, oss)
	}
}

type OldWorkOrderComment struct {
	Id          uint64    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	WorkOrderId uint64    `json:"about_id"`
	CommentId   uint64    `json:"comment_id"`
}

func ListAllWorkOrderComments(db *sql.DB) ([]*OldWorkOrderComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, created_at, about_id, comment_id
	FROM
        workery_work_order_comments
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldWorkOrderComment
	defer rows.Close()
	for rows.Next() {
		m := new(OldWorkOrderComment)
		err = rows.Scan(
			&m.Id,
			&m.CreatedAt,
			&m.WorkOrderId,
			&m.CommentId,
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

func insertWorkOrderCommentETL(
	ctx context.Context,
	tid uint64,
	wotp *repositories.WorkOrderCommentRepo,
	ar *repositories.WorkOrderRepo,
	vtr *repositories.CommentRepo,
	oss *OldWorkOrderComment,
) {
	//
	// OrderId
	//

	orderId, err := ar.GetIdByOldId(ctx, tid, oss.WorkOrderId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	//
	// CommentId
	//

	commentId, err := vtr.GetIdByOldId(ctx, tid, oss.CommentId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	//
	// Insert into database.
	//

	m := &models.WorkOrderComment{
		OldId:       oss.Id,
		TenantId:    tid,
		Uuid:        uuid.NewString(),
		CreatedTime: oss.CreatedAt,
		OrderId:     orderId,
		CommentId:   commentId,
	}
	err = wotp.Insert(ctx, m)
	if err != nil {
		log.Panic("wotp.Insert | err", err)
	}
	fmt.Println("Imported ID#", oss.Id)
}
