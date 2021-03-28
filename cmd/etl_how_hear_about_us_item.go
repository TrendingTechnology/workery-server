package cmd

import (
	"context"
	"fmt"
	// "encoding/csv"
	"database/sql"
	"os"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/google/uuid"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	etlhhauiTenantSchema string
	etlhhauiTenantId int
)

func init() {
	howHearAboutUsItemETLCmd.Flags().StringVarP(&etlhhauiTenantSchema, "schema_name", "s", "public", "Schema name of the old workery")
	howHearAboutUsItemETLCmd.MarkFlagRequired("schema_name")
	howHearAboutUsItemETLCmd.Flags().IntVarP(&etlhhauiTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	howHearAboutUsItemETLCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(howHearAboutUsItemETLCmd)
}

var howHearAboutUsItemETLCmd = &cobra.Command{
	Use:   "etl_how_hear_about_us_item",
	Short: "Import the how_hear_about_us_item data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportHowHearAboutUsItem(uint64(etlhhauiTenantId))
	},
}

func doRunImportHowHearAboutUsItem(tid uint64) {
	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repositories.NewHowHearAboutUsItemRepo(db)

	// Load up our old database.
	oldDBHost := os.Getenv("WORKERY_OLD_DB_HOST")
	oldDBPort := os.Getenv("WORKERY_OLD_DB_PORT")
	oldDBUser := os.Getenv("WORKERY_OLD_DB_USER")
	oldDBPassword := os.Getenv("WORKERY_OLD_DB_PASSWORD")
	oldDBName := os.Getenv("WORKERY_OLD_DB_NAME")
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, etlhhauiTenantSchema)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

    // Begin the operation.
	runHowHearAboutUsItemETL(tid, r, oldDb)
}

func runHowHearAboutUsItemETL(tid uint64, r *repositories.HowHearAboutUsItemRepo, oldDb *sql.DB) {
	items, err := ListAllHowHearAboutUsItems(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range items {
		runHowHearAboutUsItemInsert(tid, v, r)
	}
}

type OldHowHearAboutUsItem struct {
	Id                uint64    `json:"id"`
	Uuid              string    `json:"uuid"`
	TenantId          uint64    `json:"tenant_id"`
	Text              string    `json:"text"`
	SortNumber        int8      `json:"sort_number"`
	IsForAssociate    bool      `json:"is_for_associate"`
	IsForCustomer     bool      `json:"is_for_customer"`
	IsForStaff        bool      `json:"is_for_staff"`
	IsForPartner      bool      `json:"is_for_partner"`
	IsArchived        bool      `json:"is_archived"`
}

// Function returns a paginated list of all type element items.
func ListAllHowHearAboutUsItems(oldDb *sql.DB) ([]*OldHowHearAboutUsItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, text, sort_number, is_for_associate, is_for_customer,
		is_for_staff, is_for_partner, is_archived
	FROM
        workery_how_hear_about_us_items
	`
	rows, err := oldDb.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldHowHearAboutUsItem
	defer rows.Close()
	for rows.Next() {
		m := new(OldHowHearAboutUsItem)
		err = rows.Scan(
			&m.Id, &m.Text, &m.SortNumber, &m.IsForAssociate, &m.IsForCustomer,
			&m.IsForStaff, &m.IsForPartner, &m.IsArchived,
		)
		if err != nil {
			panic(err)
		}
		m.Uuid = uuid.NewString()
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return arr, err
}

func runHowHearAboutUsItemInsert(tid uint64, ot *OldHowHearAboutUsItem, r *repositories.HowHearAboutUsItemRepo) {
	var state int8 = 1
	if ot.IsArchived == true {
		state = 0
	}

	m := &models.HowHearAboutUsItem{
		OldId: ot.Id,
		TenantId: tid,
		Uuid: uuid.NewString(),
		Text: ot.Text,
		IsForAssociate: ot.IsForAssociate,
		IsForCustomer: ot.IsForCustomer,
		IsForStaff: ot.IsForStaff,
		IsForPartner: ot.IsForPartner,
		SortNumber: 1,
		State: state,
	}
	ctx := context.Background()
	err := r.Insert(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", ot.Id)
}
