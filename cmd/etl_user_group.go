package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	// "github.com/google/uuid"

	// "github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

func init() {
	rootCmd.AddCommand(userGroupETLCmd)
}

var userGroupETLCmd = &cobra.Command{
	Use:   "etl_user_group",
	Short: "Import the user groups from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportUserGroup()
	},
}

func doRunImportUserGroup() {
	// Load up our new database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repositories.NewUserRepo(db)

	// Load up our old database.
	oldDBHost := os.Getenv("WORKERY_OLD_DB_HOST")
	oldDBPort := os.Getenv("WORKERY_OLD_DB_PORT")
	oldDBUser := os.Getenv("WORKERY_OLD_DB_USER")
	oldDBPassword := os.Getenv("WORKERY_OLD_DB_PASSWORD")
	oldDBName := os.Getenv("WORKERY_OLD_DB_NAME")
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, "public")
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Begin the operation.
	runUserGroupETL(r, oldDb)
}

type OldUserGroup struct {
	Id      uint64 `json:"id"`
	UserId  uint64 `json:"shareduser_id"`
	GroupId uint64 `json:"group_id"`
}

// Function returns a paginated list of all type element items.
func ListAllUserGroups(db *sql.DB) ([]*OldUserGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, shareduser_id, group_id
	FROM
	    workery_users_groups
	ORDER BY
		id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUserGroup
	defer rows.Close()
	for rows.Next() {
		m := new(OldUserGroup)
		err = rows.Scan(
			&m.Id,
			&m.UserId,
			&m.GroupId,
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

func runUserGroupETL(r *repositories.UserRepo, oldDb *sql.DB) {
	userGroups, err := ListAllUserGroups(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range userGroups {
		runUserGroupInsert(v, r)
	}
}

func runUserGroupInsert(ot *OldUserGroup, r *repositories.UserRepo) {
	ctx := context.Background()
	user, err := r.GetByOldId(ctx, ot.UserId)
	if err != nil {
		panic(err)
	}
	if user != nil {
		user.Role = int8(ot.GroupId)
		r.UpdateById(ctx, user)
		fmt.Println("Processed UserId #", user.Id)
	} else {
		fmt.Println("Skipped UserId #", user.Id)
	}
}
