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
	"github.com/google/uuid"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	userFilePath string
)

func init() {
	userETLCmd.Flags().StringVarP(&userFilePath, "filepath", "f", "", "Path to the workery user csv file.")
	userETLCmd.MarkFlagRequired("filepath")
	rootCmd.AddCommand(userETLCmd)
}

var userETLCmd = &cobra.Command{
	Use:   "etl_user",
	Short: "Import the user data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportUser()
	},
}

func doRunImportUser() {
	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName)
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repositories.NewUserRepo(db)

	f, err := os.Open(userFilePath)
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

		saveUserRowInDb(r, line)
	}
}

func saveUserRowInDb(r *repositories.UserRepo, col []string) {
	// For debugging purposes only.
	// log.Println(col)

	// // Extract the row.
	idString := col[0]
	// col[1] // PasswordHash
	// col[2] // LastLoginTime
	// col[3] // IsSuperUser
	email := col[4]
	firstName := col[5]
	lastName := col[6]
	joinedTimeString := col[7]
	isActiveString := col[8]
	// col[9] // ???
	lastModifiedTimeString := col[10]
	salt := col[11]
	timezone := "America/Toronto"
	wasEmailActivatedString := col[12]
	prAccessCodeString := col[13]
	prExpiryTimeString := col[14]
	tenantIdString := col[15]

	joinedTime, _ := utils.ConvertPGAdminTimeStringToTime(joinedTimeString)
	lastModifiedTime, _ := utils.ConvertPGAdminTimeStringToTime(lastModifiedTimeString)
	prExpiryTime, _ := utils.ConvertPGAdminTimeStringToTime(prExpiryTimeString)

	var isActive int8
	if isActiveString == "t" {
		isActive = 1
	} else {
        isActive = 0
	}

	var wasEmailActivated bool
	if wasEmailActivatedString == "t" {
		wasEmailActivated = true
	} else {
		wasEmailActivated = false
	}

	id, _ := strconv.ParseUint(idString, 10, 64)

	tenantId, err := strconv.ParseUint(tenantIdString, 10, 64)
	if err != nil {
		tenantId = 1
	}

	if id != 0 {
		m := &models.User{
			Id: id,
			Uuid: uuid.NewString(),
			FirstName: firstName,
			LastName: lastName,
			Email: email,
			JoinedTime: joinedTime,
			State: isActive,
			Timezone: timezone,
			CreatedTime: joinedTime,
			ModifiedTime: lastModifiedTime,
			Salt: salt,
			WasEmailActivated: wasEmailActivated,
			PrAccessCode: prAccessCodeString,
			PrExpiryTime: prExpiryTime,
			TenantId: tenantId,
		}
		ctx := context.Background()
		err := r.InsertOrUpdateByEmail(ctx, m)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("Imported ID#", id)
	}
}
