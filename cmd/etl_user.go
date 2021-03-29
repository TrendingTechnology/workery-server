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

func init() {
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
	// Load up our new database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	ur := repositories.NewUserRepo(db)

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
	runUserETL(tr, ur, oldDb)
}

type OldUser struct {
	Id       uint64        `json:"id"`
	TenantId sql.NullInt64 `json:"franchise_id"`
	// password character varying(128) COLLATE pg_catalog."default" NOT NULL,
	// last_login timestamp with time zone,
	// is_superuser boolean NOT NULL,
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	DateJoined time.Time `json:"date_joined"`
	IsActive   bool      `json:"is_active"`
	// avatar character varying(100) COLLATE pg_catalog."default",
	LastModified time.Time `json:"last_modified"`
	// salt character varying(127) COLLATE pg_catalog."default",
	WasEmailActivated bool `json:"was_email_activated"`
	// pr_access_code character varying(127) COLLATE pg_catalog."default" NOT NULL,
	// pr_expiry_date timestamp with time zone NOT NULL,

	// IsArchived              bool   `json:"is_archived"`
}

func ListAllUsers(db *sql.DB) ([]*OldUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, email, first_name, last_name, date_joined, is_active, last_modified, was_email_activated, franchise_id
	FROM
	    workery_users
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUser
	defer rows.Close()
	for rows.Next() {
		m := new(OldUser)
		err = rows.Scan(
			&m.Id,
			&m.Email,
			&m.FirstName,
			&m.LastName,
			&m.DateJoined,
			&m.IsActive,
			&m.LastModified,
			&m.WasEmailActivated,
			&m.TenantId,
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

func runUserETL(tr *repositories.TenantRepo, ur *repositories.UserRepo, oldDb *sql.DB) {
	users, err := ListAllUsers(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range users {
		runUserInsert(v, tr, ur)
	}
}

func runUserInsert(ou *OldUser, tr *repositories.TenantRepo, ur *repositories.UserRepo) {
	log.Println(ou, ur, "\n")

	var state int8 = 0
	if ou.IsActive == true {
		state = 1
	}

	tenantId := sql.NullInt64{Int64: 1, Valid: true}
	if ou.TenantId.Valid == true {
		tenantId = sql.NullInt64{Int64: ou.TenantId.Int64, Valid: true}
	}

	ctx := context.Background()
	tenant, err := tr.GetByOldId(ctx, uint64(tenantId.Int64))
	if err != nil {
		log.Fatal(err)
	}

	m := &models.User{
		OldId:             ou.Id,
		Uuid:              uuid.NewString(),
		FirstName:         ou.FirstName,
		LastName:          ou.LastName,
		Email:             ou.Email,
		JoinedTime:        ou.DateJoined,
		State:             state,
		Timezone:          "America/Toronto",
		CreatedTime:       ou.DateJoined,
		ModifiedTime:      ou.LastModified,
		Salt:              "",
		WasEmailActivated: ou.WasEmailActivated,
		PrAccessCode:      "",
		PrExpiryTime:      time.Now(),
		TenantId:          tenant.Id,
	}
	err = ur.InsertOrUpdateByEmail(ctx, m)
	if err != nil {
		log.Println("TenantId", m.TenantId)
		log.Panic(err)
	}
	fmt.Println("Imported ID#", ou.Id)
}
