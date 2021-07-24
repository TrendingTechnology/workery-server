package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	repo "github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

// ex:
// $ go run main.go change_password --email="b@b.com" --password="123"

var (
	changePassEmail    string
	changePassPassword string
)

func init() {
	changePasswordCmd.Flags().StringVarP(&changePassEmail, "email", "e", "", "Email of the user account")
	changePasswordCmd.MarkFlagRequired("email")
	changePasswordCmd.Flags().StringVarP(&changePassPassword, "password", "p", "", "Password of the user account")
	changePasswordCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(changePasswordCmd)
}

var changePasswordCmd = &cobra.Command{
	Use:   "change_password",
	Short: "Change user password",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		runChangePassword()
	},
}

func runChangePassword() {
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
	user, _ := r.GetByEmail(ctx, changePassEmail)
	if user == nil {
		log.Fatal("User D.N.E.")
	}

	passwordHash, err := utils.HashPassword(changePassPassword)
	if err != nil {
		log.Fatal("HashPassword:", err)
	}
	user.PasswordHash = passwordHash

	err = r.UpdateByEmail(ctx, user)
	if err != nil {
		log.Fatal("UpdateByEmail:", err)
	}

	fmt.Print("\033[H\033[2J")
	fmt.Println("Password successfully changed")
}
