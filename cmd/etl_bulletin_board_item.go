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
	null "gopkg.in/guregu/null.v4"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	bulletinBoardItemETLSchemaName string
)

func init() {
	bulletinBoardItemETLCmd.Flags().StringVarP(&bulletinBoardItemETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	bulletinBoardItemETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(bulletinBoardItemETLCmd)
}

var bulletinBoardItemETLCmd = &cobra.Command{
	Use:   "etl_bulletin_board_item",
	Short: "Import the bulletin_board_item data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportBulletinBoardItem()
	},
}

func doRunImportBulletinBoardItem() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, bulletinBoardItemETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	irr := repositories.NewBulletinBoardItemRepo(db)
	ur := repositories.NewUserRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, bulletinBoardItemETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runBulletinBoardItemETL(ctx, tenant.Id, irr, ur, oldDb)
}

type OldBulletinBoardItem struct {
	Id               uint64    `json:"id"`
	Text             string    `json:"text"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedById      null.Int  `json:"created_by_id,omitempty"`
	CreatedFrom      string    `json:"created_from"`
	LastModifiedAt   time.Time `json:"last_modified_at"`
	LastModifiedById null.Int  `json:"last_modified_by_id,omitempty"`
	LastModifiedFrom string    `json:"last_modified_from"`
	IsArchived       bool      `json:"is_archived"`
}

func ListAllBulletinBoardItems(db *sql.DB) ([]*OldBulletinBoardItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, text, created_at, created_by_id, created_from, last_modified_at, last_modified_by_id, last_modified_from, is_archived
	FROM
	    workery_bulletin_board_items
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldBulletinBoardItem
	defer rows.Close()
	for rows.Next() {
		m := new(OldBulletinBoardItem)
		err = rows.Scan(
			&m.Id,
			&m.Text,
			&m.CreatedAt,
			&m.CreatedById,
			&m.CreatedFrom,
			&m.LastModifiedAt,
			&m.LastModifiedById,
			&m.LastModifiedFrom,
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

func runBulletinBoardItemETL(ctx context.Context, tenantId uint64, irr *repositories.BulletinBoardItemRepo, ur *repositories.UserRepo, oldDb *sql.DB) {
	bulletin_board_item, err := ListAllBulletinBoardItems(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oir := range bulletin_board_item {
		insertBulletinBoardItemETL(ctx, tenantId, irr, ur, oir)
	}
}

func insertBulletinBoardItemETL(ctx context.Context, tid uint64, irr *repositories.BulletinBoardItemRepo, ur *repositories.UserRepo, oir *OldBulletinBoardItem) {
	//
	// Set the `state`.
	//

	var state int8 = 1
	if oir.IsArchived == true {
		state = 0
	}

	//
	// Get `createdById` and `createdByName` values.
	//

	var createdById null.Int
	var createdByName null.String
	if oir.CreatedById.ValueOrZero() > 0 {
		userId, err := ur.GetIdByOldId(ctx, tid, uint64(oir.CreatedById.ValueOrZero()))

		if err != nil {
			log.Panic("ur.GetIdByOldId", err)
		}
		user, err := ur.GetById(ctx, userId)
		if err != nil {
			log.Panic("ur.GetById", err)
		}

		if user != nil {
			createdById = null.IntFrom(int64(userId))
			createdByName = null.StringFrom(user.Name)
		} else {
			log.Println("WARNING: D.N.E.")
		}

		// // For debugging purposes only.
		// log.Println("createdById:", createdById)
		// log.Println("createdByName:", createdByName)
	}

	//
	// Get `lastModifiedById` and `lastModifiedByName` values.
	//

	var lastModifiedById null.Int
	var lastModifiedByName null.String
	if oir.LastModifiedById.ValueOrZero() > 0 {
		userId, err := ur.GetIdByOldId(ctx, tid, uint64(oir.LastModifiedById.ValueOrZero()))
		if err != nil {
			log.Panic("ur.GetIdByOldId", err)
		}
		user, err := ur.GetById(ctx, userId)
		if err != nil {
			log.Panic("ur.GetById", err)
		}

		if user != nil {
			lastModifiedById = null.IntFrom(int64(userId))
			lastModifiedByName = null.StringFrom(user.Name)
		} else {
			log.Println("WARNING: D.N.E.")
		}

		// // For debugging purposes only.
		// log.Println("lastModifiedById:", lastModifiedById)
		// log.Println("lastModifiedByName:", lastModifiedByName)
	}

	//
	// Insert the `BulletinBoardItem`.
	//

	m := &models.BulletinBoardItem{
		OldId:              oir.Id,
		TenantId:           tid,
		Uuid:               uuid.NewString(),
		Text:               oir.Text,
		CreatedTime:        oir.CreatedAt,
		CreatedById:        createdById,
		CreatedByName:      createdByName,
		CreatedFromIP:      oir.CreatedFrom,
		LastModifiedTime:   oir.LastModifiedAt,
		LastModifiedById:   lastModifiedById,
		LastModifiedByName: lastModifiedByName,
		LastModifiedFromIP: oir.LastModifiedFrom,
		State:              state,
	}
	err := irr.Insert(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", oir.Id)
}
