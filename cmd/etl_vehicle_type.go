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
	vehicleTypeETLSchemaName string
)

func init() {
	vehicleTypeETLCmd.Flags().StringVarP(&vehicleTypeETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	vehicleTypeETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(vehicleTypeETLCmd)
}

var vehicleTypeETLCmd = &cobra.Command{
	Use:   "etl_vehicle_type",
	Short: "Import the vehicleType data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportVehicleType()
	},
}

func doRunImportVehicleType() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, vehicleTypeETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

    // Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	irr := repositories.NewVehicleTypeRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, vehicleTypeETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runVehicleTypeETL(ctx, tenant.Id, irr, oldDb)
}

type OldUVehicleType struct {
	Id                uint64    `json:"id"`
	Text              string    `json:"text"`
	Description       string    `json:"description"`
	IsArchived        bool      `json:"is_archived"`
}

func ListAllVehicleTypes(db *sql.DB) ([]*OldUVehicleType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, text, description, is_archived
	FROM
	    workery_vehicle_types
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUVehicleType
	defer rows.Close()
	for rows.Next() {
		m := new(OldUVehicleType)
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

func runVehicleTypeETL(ctx context.Context, tenantId uint64, irr *repositories.VehicleTypeRepo, oldDb *sql.DB) {
	vehicleTypes, err := ListAllVehicleTypes(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oir := range vehicleTypes {
		insertVehicleTypeETL(ctx, tenantId, irr, oir)
	}
}

func insertVehicleTypeETL(ctx context.Context, tid uint64, irr *repositories.VehicleTypeRepo, oir *OldUVehicleType) {
	var state int8 = 1
	if oir.IsArchived == true {
		state = 0
	}

	m := &models.VehicleType{
		OldId: oir.Id,
		TenantId: tid,
		Uuid: uuid.NewString(),
		Text: oir.Text,
		Description: oir.Description,
		State: state,
	}
	err := irr.Insert(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", oir.Id)
}
