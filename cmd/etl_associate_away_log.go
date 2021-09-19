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
	awlETLSchemaName string
	awlETLTenantId   int
)

func init() {
	awlETLCmd.Flags().StringVarP(&awlETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	awlETLCmd.MarkFlagRequired("schema_name")
	awlETLCmd.Flags().IntVarP(&awlETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	awlETLCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(awlETLCmd)
}

var awlETLCmd = &cobra.Command{
	Use:   "etl_associate_away_log",
	Short: "Import the associate away logs from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportAssociateAwayLog()
	},
}

func doRunImportAssociateAwayLog() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, awlETLSchemaName)
	if err != nil {
		log.Fatal("utils.ConnectDB", err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	asr := repositories.NewAssociateAwayLogRepo(db)
	ar := repositories.NewAssociateRepo(db)
	vtr := repositories.NewAssociateAwayLogRepo(db)
	ur := repositories.NewUserRepo(db)

	runAssociateAwayLogETL(ctx, uint64(awlETLTenantId), asr, ar, vtr, ur, oldDb)
}

func runAssociateAwayLogETL(
	ctx context.Context,
	tenantId uint64,
	asr *repositories.AssociateAwayLogRepo,
	ar *repositories.AssociateRepo,
	vtr *repositories.AssociateAwayLogRepo,
	ur *repositories.UserRepo,
	oldDb *sql.DB,
) {
	ass, err := ListAllAssociateAwayLogs(oldDb)
	if err != nil {
		log.Fatal("ListAllAssociateAwayLogs", err)
	}
	for _, oss := range ass {
		insertAssociateAwayLogETL(ctx, tenantId, asr, ar, vtr, ur, oss)
	}
}

type OldAssociateAwayLog struct {
	Id                 uint64      `json:"id"`
	AssociateId        uint64      `json:"associate_id"`
	Reason             int8        `json:"reason"`
	ReasonOther        null.String `json:"reason_other"`
	UntilFurtherNotice bool        `json:"until_further_notice"`
	UntilDate          null.Time   `json:"until_date"`
	StartDate          null.Time   `json:"start_date"`
	WasDeleted         bool        `json:"was_deleted"`
	CreatedTime        time.Time   `json:"created"`
	CreatedById        null.Int    `json:"created_by_id"`
	LastModifiedTime   time.Time   `json:"last_modified"`
	LastModifiedById   null.Int    `json:"last_modified_by_id"`
}

func ListAllAssociateAwayLogs(db *sql.DB) ([]*OldAssociateAwayLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, associate_id, reason, reason_other, until_further_notice, until_date,
		start_date, was_deleted, created, created_by_id,
		last_modified, last_modified_by_id
	FROM
        workery_away_logs
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldAssociateAwayLog
	defer rows.Close()
	for rows.Next() {
		m := new(OldAssociateAwayLog)
		err = rows.Scan(
			&m.Id,
			&m.AssociateId,
			&m.Reason,
			&m.ReasonOther,
			&m.UntilFurtherNotice,
			&m.UntilDate,
			&m.StartDate,
			&m.WasDeleted,
			&m.CreatedTime,
			&m.CreatedById,
			&m.LastModifiedTime,
			&m.LastModifiedById,
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

func insertAssociateAwayLogETL(
	ctx context.Context,
	tid uint64,
	asr *repositories.AssociateAwayLogRepo,
	ar *repositories.AssociateRepo,
	vtr *repositories.AssociateAwayLogRepo,
	ur *repositories.UserRepo,
	oss *OldAssociateAwayLog,
) {
	//
	// Get by the `Associate`.
	//

	associateId, err := ar.GetIdByOldId(ctx, tid, oss.AssociateId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	associate, err := ar.GetById(ctx, associateId)
	if err != nil {
		log.Panic("ar.GetById | err", err)
	}

	//
	// Set the `state`.
	//

	var state int8 = 1
	if oss.WasDeleted == true {
		state = 0
	}

	//
	// Get `createdById` and `createdByName` values.
	//

	var createdById null.Int
	var createdByName null.String
	if oss.CreatedById.ValueOrZero() > 0 {
		userId, err := ur.GetIdByOldId(ctx, tid, uint64(oss.CreatedById.ValueOrZero()))

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
	if oss.LastModifiedById.ValueOrZero() > 0 {
		userId, err := ur.GetIdByOldId(ctx, tid, uint64(oss.LastModifiedById.ValueOrZero()))
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
	// Insert the `AssociateAwayLog` if it exis.
	//

	if associateId != 0 {
		m := &models.AssociateAwayLog{
			OldId:                oss.Id,
			TenantId:             tid,
			Uuid:                 uuid.NewString(),
			AssociateId:          associateId,
			AssociateName:        associate.Name,
			AssociateLexicalName: associate.LexicalName,
			Reason:               oss.Reason,
			ReasonOther:          oss.ReasonOther,
			UntilFurtherNotice:   oss.UntilFurtherNotice,
			UntilDate:            oss.UntilDate,
			StartDate:            oss.StartDate,
			State:                state,
			CreatedTime:          oss.CreatedTime,
			CreatedById:          createdById,
			CreatedByName:        createdByName,
			LastModifiedTime:     oss.LastModifiedTime,
			LastModifiedById:     lastModifiedById,
			LastModifiedByName:   lastModifiedByName,
		}
		err = asr.Insert(ctx, m)
		if err != nil {
			log.Panic("asr.Insert | err", err, "\n\n", m, oss)
		}
		fmt.Println("Imported ID#", oss.Id)
	} else {
		fmt.Println("-------------------\nSkipped ID#", oss.Id)
	}
}
