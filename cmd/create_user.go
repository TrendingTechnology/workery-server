package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/over55/workery-server/internal/models"
	repo "github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	cuTenantId  int
	cuFirstName string
	cuLastName  string
	cuEmail     string
	cuPassword  string
	cuState     int
	cuRoleId    int
)

func init() {
	createUserCmd.Flags().IntVarP(&cuTenantId, "tid", "t", 0, "Tenant Id of the user account")
	createUserCmd.MarkFlagRequired("tenant_id")
	createUserCmd.Flags().StringVarP(&cuFirstName, "fname", "f", "", "First name of the user account")
	createUserCmd.MarkFlagRequired("fname")
	createUserCmd.Flags().StringVarP(&cuLastName, "lname", "l", "", "Last name of the user account")
	createUserCmd.MarkFlagRequired("lname")
	createUserCmd.Flags().StringVarP(&cuEmail, "email", "e", "", "Email of the user account")
	createUserCmd.MarkFlagRequired("email")
	createUserCmd.Flags().StringVarP(&cuPassword, "password", "p", "", "Password of the user account")
	createUserCmd.MarkFlagRequired("password")
	createUserCmd.Flags().IntVarP(&cuState, "state", "s", 0, "State of the user account")
	createUserCmd.MarkFlagRequired("state")
	createUserCmd.Flags().IntVarP(&cuRoleId, "role_id", "r", 0, "RoleId of the user account")
	createUserCmd.MarkFlagRequired("role_id")
	rootCmd.AddCommand(createUserCmd)
}

var createUserCmd = &cobra.Command{
	Use:   "create_user",
	Short: "Create a user account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		runAddUser()
	},
}

func runAddUser() {
	ctx := context.Background()

	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repo.NewUserRepo(db)

	// Check to see the user account already exists.
	userFound, _ := r.CheckIfExistsByEmail(ctx, cuEmail)
	if userFound {
		log.Fatal("Email already exists.")
	}

	passwordHash, err := utils.HashPassword(cuPassword)
	if err != nil {
		log.Fatal(err)
	}

	m := &models.User{
		Uuid:              uuid.NewString(),
		TenantId:          uint64(cuTenantId),
		FirstName:         cuFirstName,
		LastName:          cuLastName,
		Email:             cuEmail,
		PasswordAlgorithm: utils.HashPasswordAlgorithm(),
		PasswordHash:      passwordHash,
		State:             int8(cuState),
		RoleId:            int8(cuRoleId),
		Timezone:          "utc",
		CreatedTime:       time.Now(),
		ModifiedTime:      time.Now(),
	}

	err = r.InsertOrUpdateById(ctx, m)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("\033[H\033[2J")
	fmt.Println("User created with UUID:", m.Uuid)
}
