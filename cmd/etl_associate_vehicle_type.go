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
	avtETLSchemaName string
	avtETLTenantId   int
)

func init() {
	avtETLCmd.Flags().StringVarP(&avtETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	avtETLCmd.MarkFlagRequired("schema_name")
	avtETLCmd.Flags().IntVarP(&avtETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	avtETLCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(avtETLCmd)
}

var avtETLCmd = &cobra.Command{
	Use:   "etl_associate_vehicle_type",
	Short: "Import the associate vehicle types from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportAssociateVehicleType()
	},
}

func doRunImportAssociateVehicleType() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, avtETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	avtr := repositories.NewAssociateVehicleTypeRepo(db)
	ar := repositories.NewAssociateRepo(db)
	vtr := repositories.NewVehicleTypeRepo(db)

	runAssociateVehicleTypeETL(ctx, uint64(avtETLTenantId), avtr, ar, vtr, oldDb)
}

func runAssociateVehicleTypeETL(
	ctx context.Context,
	tenantId uint64,
	avtr *repositories.AssociateVehicleTypeRepo,
	ar *repositories.AssociateRepo,
	vtr *repositories.VehicleTypeRepo,
	oldDb *sql.DB,
) {
	avts, err := ListAllAssociateVehicleTypes(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oss := range avts {
		insertAssociateVehicleTypeETL(ctx, tenantId, avtr, ar, vtr, oss)
	}
}

type OldAssociateVehicleType struct {
	Id            uint64 `json:"id"`
	AssociateId   uint64 `json:"associate_id"`
	VehicleTypeId uint64 `json:"vehicletype_id"`
}

func ListAllAssociateVehicleTypes(db *sql.DB) ([]*OldAssociateVehicleType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, associate_id, vehicletype_id
	FROM
        workery_associates_vehicle_types
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldAssociateVehicleType
	defer rows.Close()
	for rows.Next() {
		m := new(OldAssociateVehicleType)
		err = rows.Scan(
			&m.Id,
			&m.AssociateId,
			&m.VehicleTypeId,
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

func insertAssociateVehicleTypeETL(
	ctx context.Context,
	tid uint64,
	avtr *repositories.AssociateVehicleTypeRepo,
	ar *repositories.AssociateRepo,
	vtr *repositories.VehicleTypeRepo,
	oss *OldAssociateVehicleType,
) {
	associateId, err := ar.GetIdByOldId(ctx, tid, oss.AssociateId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	vehicleTypeId, err := vtr.GetIdByOldId(ctx, tid, oss.VehicleTypeId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	m := &models.AssociateVehicleType{
		OldId:         oss.Id,
		TenantId:      tid,
		Uuid:          uuid.NewString(),
		AssociateId:   associateId,
		VehicleTypeId: vehicleTypeId,
	}
	err = avtr.Insert(ctx, m)
	if err != nil {
		log.Panic("avtr.Insert | err", err)
	}
	fmt.Println("Imported ID#", oss.Id)
}
