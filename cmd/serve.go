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
	db, err := utils.ConnectDB(
		databaseHost,
		databasePort,
		databaseUser,
		databasePassword,
		databaseName,
		"public",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	asir := repo.NewActivitySheetItemRepo(db)
	aalr := repo.NewAssociateAwayLogRepo(db)
	acr := repo.NewAssociateCommentRepo(db)
	airr := repo.NewAssociateInsuranceRequirementRepo(db)
	assr := repo.NewAssociateSkillSetRepo(db)
	atr := repo.NewAssociateTagRepo(db)
	avtr := repo.NewAssociateVehicleTypeRepo(db)
	ar := repo.NewAssociateRepo(db)
	bbir := repo.NewBulletinBoardItemRepo(db)
	comr := repo.NewCommentRepo(db)
	ccr := repo.NewCustomerCommentRepo(db)
	ctr := repo.NewCustomerTagRepo(db)
	cr := repo.NewCustomerRepo(db)
	hhauir := repo.NewHowHearAboutUsItemRepo(db)
	irr := repo.NewInsuranceRequirementRepo(db)
	lar := repo.NewLiteAssociateRepo(db)
	lcr := repo.NewLiteCustomerRepo(db)
	ltar := repo.NewLiteTaskItemRepo(db)
	ltr := repo.NewLiteTenantRepo(db)
	lwor := repo.NewLiteWorkOrderRepo(db)
	lpr := repo.NewLitePartnerRepo(db)
	lowor := repo.NewLiteOngoingWorkOrderRepo(db)
	owor := repo.NewOngoingWorkOrderRepo(db)
	pcr := repo.NewPartnerCommentRepo(db)
	pr := repo.NewPartnerRepo(db)
	pfr := repo.NewPrivateFileRepo(db)
	// piur := repo.NewPublicImageUploadRepo(db)
	skillsirr := repo.NewSkillSetInsuranceRequirementRepo(db)
	skillsr := repo.NewSkillSetRepo(db)
	staffcr := repo.NewStaffCommentRepo(db)
	staffTagr := repo.NewStaffTagRepo(db)
	staffr := repo.NewStaffRepo(db)
	tagr := repo.NewTagRepo(db)
	tir := repo.NewTaskItemRepo(db)
	tr := repo.NewTenantRepo(db)
	ur := repo.NewUserRepo(db)
	vtr := repo.NewVehicleTypeRepo(db)
	wocr := repo.NewWorkOrderCommentRepo(db)
	wodr := repo.NewWorkOrderDepositRepo(db)
	woir := repo.NewWorkOrderInvoiceRepo(db)
	wosfr := repo.NewWorkOrderServiceFeeRepo(db)
	wossr := repo.NewWorkOrderSkillSetRepo(db)
	wotr := repo.NewWorkOrderTagRepo(db)
	wor := repo.NewWorkOrderRepo(db)

	// Open up our session handler, powered by redis and let's save the user
	// account with our ID
	sm := session.New()

	// Instead of using a `New` sort of function, we will populate our structure
	// so we can use it.
	c := &controllers.Controller{
		SecretSigningKeyBin:               []byte(applicationSigningKey),
		ActivitySheetItemRepo:             asir,
		AssociateAwayLogRepo:              aalr,
		AssociateCommentRepo:              acr,
		AssociateInsuranceRequirementRepo: airr,
		AssociateSkillSetRepo:             assr,
		AssociateTagRepo:                  atr,
		AssociateVehicleTypeRepo:          avtr,
		AssociateRepo:                     ar,
		BulletinBoardItemRepo:             bbir,
		CommentRepo:                       comr,
		CustomerCommentRepo:               ccr,
		CustomerTagRepo:                   ctr,
		CustomerRepo:                      cr,
		HowHearAboutUsItemRepo:            hhauir,
		InsuranceRequirementRepo:          irr,
		LiteAssociateRepo:                 lar,
		LiteCustomerRepo:                  lcr,
		LiteTaskItemRepo:                  ltar,
		LiteTenantRepo:                    ltr,
		LiteWorkOrderRepo:                 lwor,
		LitePartnerRepo:                   lpr,
		LiteOngoingWorkOrderRepo:          lowor,
		OngoingWorkOrderRepo:              owor,
		PartnerCommentRepo:                pcr,
		PartnerRepo:                       pr,
		PrivateFileRepo:                   pfr,
		// PublicImageUploadRepo:          piur,
		SkillSetInsuranceRequirementRepo: skillsirr,
		SkillSetRepo:                     skillsr,
		StaffCommentRepo:                 staffcr,
		StaffTagRepo:                     staffTagr,
		StaffRepo:                        staffr,
		TagRepo:                          tagr,
		TaskItemRepo:                     tir,
		TenantRepo:                       tr,
		UserRepo:                         ur,
		VehicleTypeRepo:                  vtr,
		WorkOrderCommentRepo:             wocr,
		WorkOrderDepositRepo:             wodr,
		WorkOrderInvoiceRepo:             woir,
		WorkOrderServiceFeeRepo:          wosfr,
		WorkOrderSkillSetRepo:            wossr,
		WorkOrderTagRepo:                 wotr,
		WorkOrderRepo:                    wor,
		SessionManager:                   sm,
	}

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
