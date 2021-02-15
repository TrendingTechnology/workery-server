package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/google/uuid"

	"github.com/over55/workery-server/internal/models"
	repo "github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	ctName string
	ctState int
)

func init() {
	createTenantCmd.Flags().StringVarP(&ctName, "name", "n", "", "Name of the tenant")
	createTenantCmd.MarkFlagRequired("name")
	createTenantCmd.Flags().IntVarP(&ctState, "state", "s", 0, "State of the tenant")
	createTenantCmd.MarkFlagRequired("state")
	rootCmd.AddCommand(createTenantCmd)
}

var createTenantCmd = &cobra.Command{
	Use:   "create_tenant",
	Short: "Create a tenant",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		runCreateTenant()
	},
}

func runCreateTenant() {
	ctx := context.Background()

	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName)
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repo.NewTenantRepo(db)

	// Check to see the tenant account already exists.
	tenantFound, _ := r.CheckIfExistsByName(ctx, ctName)
	if tenantFound {
		log.Fatal("Email already exists.")
	}

	m := &models.Tenant{
		Uuid: uuid.NewString(),
		Name: ctName,
		State: int8(ctState),
		Timezone: "utc",
		CreatedTime: time.Now(),
		ModifiedTime: time.Now(),
	}

	err = r.InsertOrUpdate(ctx, m)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("\033[H\033[2J")
	fmt.Println("Tenant created.")
}
