package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
	"github.com/spf13/cobra"

	"github.com/over55/workery-server/internal/controllers"
	repo "github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/session"
	"github.com/over55/workery-server/internal/utils"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the JSON API over HTTP",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		runServeCmd()
	},
}

func doRunServe() {
	fmt.Println("Server started")
}

func runServeCmd() {
	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	ur := repo.NewUserRepo(db)

	// Open up our session handler, powered by redis and let's save the user
	// account with our ID
	sm := session.New()

	c := controllers.NewBaseHandler([]byte(applicationSigningKey), ur, sm)

	mux := http.NewServeMux()
	mux.HandleFunc("/", c.AttachMiddleware(c.HandleRequests))

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation via `https://github.com/rs/cors` for more options.
	handler := cors.AllowAll().Handler(mux)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "localhost", "5000"),
		Handler: handler,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go runMainRuntimeLoop(srv)

	log.Print("Server Started")

	// Run the main loop blocking code.
	<-done

	stopMainRuntimeLoop(srv)
}

func runMainRuntimeLoop(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func stopMainRuntimeLoop(srv *http.Server) {
	log.Printf("Starting graceful shutdown now...")

	// Execute the graceful shutdown sub-routine which will terminate any
	// active connections and reject any new connections.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Printf("Graceful shutdown finished.")
	log.Print("Server Exited")
}
