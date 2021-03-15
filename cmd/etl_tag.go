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
	tagETLTenantId int
	tagETLFilePath string
)

func init() {
	tagETLCmd.Flags().IntVarP(&tagETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	tagETLCmd.MarkFlagRequired("tenant_id")
	tagETLCmd.Flags().StringVarP(&tagETLFilePath, "filepath", "f", "", "Path to the workery tag csv file.")
	tagETLCmd.MarkFlagRequired("filepath")
	rootCmd.AddCommand(tagETLCmd)
}

var tagETLCmd = &cobra.Command{
	Use:   "etl_tag",
	Short: "Import the tag data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportTag()
	},
}

func doRunImportTag() {
	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repositories.NewTagRepo(db)

	f, err := os.Open(tagETLFilePath)
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

		saveTagRowInDb(r, line)
	}
}

func saveTagRowInDb(r *repositories.TagRepo, col []string) {
	// For debugging purposes only.
	// log.Println(col)

	// Extract the row.
	idString := col[0]
	text := col[1]
	description := col[2]
	stateString := col[3]

	var state int8 = 1
	if stateString == "t" {
		state = 0
	}

	id, _ := strconv.ParseUint(idString, 10, 64)
	if id != 0 {
		m := &models.Tag{
			OldId: id,
			TenantId: uint64(tagETLTenantId),
			Uuid: uuid.NewString(),
			Text: text,
			Description: description,
			State: state,
		}
		ctx := context.Background()
		err := r.Insert(ctx, m)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("Imported ID#", id)
	}
}
