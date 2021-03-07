package cmd

import (
	"bufio"
	"context"
	"database/sql"
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
	commentETLTenantId int
	commentETLFilePath string
)

func init() {
	commentETLCmd.Flags().IntVarP(&commentETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	commentETLCmd.MarkFlagRequired("tenant_id")
	commentETLCmd.Flags().StringVarP(&commentETLFilePath, "filepath", "f", "", "Path to the workery comment csv file.")
	commentETLCmd.MarkFlagRequired("filepath")
	rootCmd.AddCommand(commentETLCmd)
}

var commentETLCmd = &cobra.Command{
	Use:   "etl_comment",
	Short: "Import the comment data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportComment()
	},
}

func doRunImportComment() {
	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName)
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repositories.NewCommentRepo(db)

	f, err := os.Open(commentETLFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// defer the closing of our `f` so that we can parse it later on
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))

    // Fixes `extraneous or missing " in quoted-field` error via https://stackoverflow.com/a/51687729
	reader.LazyQuotes = true
	// reader.Comma = ';'

    var previousLine []string
    var previousId uint64
	for {
		// Read line by line until no more lines left.
		line, error := reader.Read()
		if error == io.EOF {
			break
        } else if error != nil {
			fmt.Println("\n--------------\nPrevious ID #1", previousId, "\n\npreviousLine:", previousLine, "\n\ncurrentLine" ,line, "\n--------------\n")
			// log.Panic(error)
			continue // Skip this loop iteration
		}

		previousId = saveCommentRowInDb(r, line)
		previousLine = line
	}
}

func saveCommentRowInDb(r *repositories.CommentRepo, col []string) uint64 {
	// For debugging purposes only.
	// log.Println(col)

	// Extract the row.
	idStr := col[0]
	createdTimeStr := col[1]
	// created_from := col[2]
	// created_from_is_public := col[3]
	lastModifiedTimeStr := col[4]
	// last_modified_from := col[5]
	// last_modified_from_is_public := col[6]
	text := col[7]
	stateStr := col[8]
	createdByIdStr := col[9]
	lastModifiedByIdStr := col[10]

    id, _ := strconv.ParseUint(idStr, 10, 64)

	var cbId sql.NullInt64
	createdById, err := strconv.ParseInt(createdByIdStr, 10, 64)
	if err != nil {
		cbId = sql.NullInt64{Int64: 0, Valid: false}
	} else {
		cbId = sql.NullInt64{Int64: createdById, Valid: true}
	}

	var lmId sql.NullInt64
	lastModifiedById, err := strconv.ParseInt(lastModifiedByIdStr, 10, 64)
	if err != nil {
		lmId = sql.NullInt64{Int64: 0, Valid: true}
	} else {
		lmId = sql.NullInt64{Int64: lastModifiedById, Valid: true}
	}

	createdTime, _ := utils.ConvertPGAdminTimeStringToTime(createdTimeStr)
	lastModifiedTime, _ := utils.ConvertPGAdminTimeStringToTime(lastModifiedTimeStr)

	var state int8 = 1
	if stateStr == "t" {
		state = 0
	}

	if id != 0 {
		m := &models.Comment{
			OldId: id,
			TenantId: uint64(commentETLTenantId),
			Uuid: uuid.NewString(),
			CreatedTime: createdTime,
			CreatedById: cbId,
			LastModifiedTime: lastModifiedTime,
			LastModifiedById: lmId,
			Text: text,
			State: state,
		}
		ctx := context.Background()
		err := r.Insert(ctx, m)
		if err != nil {
			log.Println(col)
			log.Panic(err)
		}
		// fmt.Println("Imported ID#", id)
	}

	return id
}
