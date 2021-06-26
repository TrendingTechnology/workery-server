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
	staffCommentETLSchemaName string
)

func init() {
	staffCommentETLCmd.Flags().StringVarP(&staffCommentETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	staffCommentETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(staffCommentETLCmd)
}

var staffCommentETLCmd = &cobra.Command{
	Use:   "etl_staff_comment",
	Short: "Import the staffComment data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportStaffComment()
	},
}

func doRunImportStaffComment() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, staffCommentETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	irr := repositories.NewStaffCommentRepo(db)
	om := repositories.NewStaffRepo(db)
	commentr := repositories.NewCommentRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, staffCommentETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runStaffCommentETL(ctx, tenant.Id, irr, om, commentr, oldDb)
}

type OldUStaffComment struct {
	Id        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	StaffId   uint64    `json:"staff_id"`
	CommentId uint64    `json:"comment_id"`
}

func ListAllStaffComments(db *sql.DB) ([]*OldUStaffComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, created_at, about_id, comment_id
	FROM
	    workery_staff_comments
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUStaffComment
	defer rows.Close()
	for rows.Next() {
		m := new(OldUStaffComment)
		err = rows.Scan(
			&m.Id,
			&m.CreatedAt,
			&m.StaffId,
			&m.CommentId,
		)
		if err != nil {
			log.Panic("ListAllStaffComments | Next | err", err)
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return arr, err
}

func runStaffCommentETL(
	ctx context.Context,
	tenantId uint64,
	irr *repositories.StaffCommentRepo,
	om *repositories.StaffRepo,
	tr *repositories.CommentRepo,
	oldDb *sql.DB,
) {
	staffComments, err := ListAllStaffComments(oldDb)
	if err != nil {
		log.Fatal("runStaffCommentETL | ListAllStaffComments | err", err)
	}
	for _, oir := range staffComments {

		staffId, err := om.GetIdByOldId(ctx, tenantId, oir.StaffId)
		if err != nil {
			log.Panic("runStaffCommentETL | om.GetIdByOldId | err", err)
		}

		commentId, err := tr.GetIdByOldId(ctx, tenantId, oir.CommentId)
		if err != nil {
			log.Panic("runStaffCommentETL | tr.GetIdByOldId | err", err)
		}

		insertStaffCommentETL(ctx, tenantId, oir.Id, staffId, commentId, irr)
	}
}

func insertStaffCommentETL(ctx context.Context, tenantId uint64, oldId uint64, staffId uint64, commentId uint64, irr *repositories.StaffCommentRepo) {
	fmt.Println("Pre-Imported Staff Comment ID#", oldId)
	m := &models.StaffComment{
		OldId:     oldId,
		Uuid:      uuid.NewString(),
		TenantId:  tenantId,
		StaffId:   staffId,
		CommentId: commentId,
	}
	err := irr.Insert(ctx, m)
	if err != nil {
		log.Panic("insertStaffCommentETL | Insert | err:", err)
	}
	fmt.Println("Post-Imported Staff Comment ID#", oldId)
}
