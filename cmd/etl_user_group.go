package cmd

import (
	"bufio"
	"context"
	"fmt"
	"encoding/csv"
	"os"
	"log"
	"io"
	"strconv"

	"github.com/spf13/cobra"
	// "github.com/google/uuid"

	// "github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	userGroupFilePath string
)

func init() {
	userGroupETLCmd.Flags().StringVarP(&userGroupFilePath, "filepath", "f", "", "Path to the workery user group csv file.")
	userGroupETLCmd.MarkFlagRequired("filepath")
	rootCmd.AddCommand(userGroupETLCmd)
}

var userGroupETLCmd = &cobra.Command{
	Use:   "etl_user_group",
	Short: "Import the user group data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportUserGroup()
	},
}

func doRunImportUserGroup() {
	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName)
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repositories.NewUserRepo(db)

	f, err := os.Open(userGroupFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// defer the closing of our `f` so that we can parse it later on
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))

	for {
		// Read line by line until no more lines left.
		line, error := reader.Read()
		if error == io.EOF {
			break
        } else if error != nil {
			log.Fatal(error)
		}

		saveUserGroupRowInDb(r, line)
	}
}

func saveUserGroupRowInDb(r *repositories.UserRepo, col []string) {
	// For debugging purposes only.
	// log.Println(col)

	// // Extract the row.
	// idString := col[0]
	userIdString := col[1]
	roleIdString := col[2]

	userId, _ := strconv.ParseUint(userIdString, 10, 64)

	roleId, err := strconv.ParseUint(roleIdString, 10, 64)
	if err != nil {
		panic(err)
	}

    ctx := context.Background()
	user, err := r.GetById(ctx, userId)
	if err != nil {
		panic(err)
	}
	if user != nil {
		user.Role = int8(roleId)
		r.UpdateById(ctx, user)
		fmt.Println("Processed UserId #", userIdString)
	} else {
		fmt.Println("Skipped UserId #", userIdString)
	}
}
