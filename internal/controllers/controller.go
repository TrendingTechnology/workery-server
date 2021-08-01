package controllers

import (
	"net/http"

	// "github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/session"
)

type Controller struct {
	SecretSigningKeyBin               []byte
	ActivitySheetItemRepo             models.ActivitySheetItemRepository
	AssociateCommentRepo              models.AssociateCommentRepository
	AssociateInsuranceRequirementRepo models.AssociateInsuranceRequirementRepository
	AssociateSkillSetRepo             models.AssociateSkillSetRepository
	AssociateTagRepo                  models.AssociateTagRepository
	AssociateVehicleTypeRepo          models.AssociateVehicleTypeRepository
	AssociateRepo                     models.AssociateRepository
	BulletinBoardItemRepo             models.BulletinBoardItemRepository
	CommentRepo                       models.CommentRepository
	CustomerCommentRepo               models.CustomerCommentRepository
	CustomerTagRepo                   models.CustomerTagRepository
	CustomerRepo                      models.CustomerRepository
	HowHearAboutUsItemRepo            models.HowHearAboutUsItemRepository
	InsuranceRequirementRepo          models.InsuranceRequirementRepository
	LiteTenantRepo                    models.LiteTenantRepository
	OngoingWorkOrderRepo              models.OngoingWorkOrderRepository
	PartnerCommentRepo                models.PartnerCommentRepository
	PartnerRepo                       models.PartnerRepository
	PrivateFileRepo                   models.PrivateFileRepository
	PublicImageUploadRepo             models.PublicImageUploadRepository
	SkillSetInsuranceRequirementRepo  models.SkillSetInsuranceRequirementRepository
	SkillSetRepo                      models.SkillSetRepository
	StaffCommentRepo                  models.StaffCommentRepository
	StaffTagRepo                      models.StaffTagRepository
	StaffRepo                         models.StaffRepository
	TagRepo                           models.TagRepository
	TaskItemRepo                      models.TaskItemRepository
	TenantRepo                        models.TenantRepository
	UserRepo                          models.UserRepository
	VehicleTypeRepo                   models.VehicleTypeRepository
	WorkOrderCommentRepo              models.WorkOrderCommentRepository
	WorkOrderDepositRepo              models.WorkOrderDepositRepository
	WorkOrderInvoiceRepo              models.WorkOrderInvoiceRepository
	WorkOrderServiceFeeRepo           models.WorkOrderServiceFeeRepository
	WorkOrderSkillSetRepo             models.WorkOrderSkillSetRepository
	WorkOrderTagRepo                  models.WorkOrderTagRepository
	WorkOrderRepo                     models.WorkOrderRepository
	SessionManager                    *session.SessionManager
}

func (h *Controller) HandleRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get our URL paths which are slash-seperated.
	ctx := r.Context()
	p := ctx.Value("url_split").([]string)
	n := len(p)

	switch {
	case n == 2 && p[0] == "v1" && p[1] == "tenants" && r.Method == http.MethodGet:
		h.liteTenantsListEndpoint(w, r)
	case n == 2 && p[0] == "v1" && p[1] == "franchises" && r.Method == http.MethodGet: // Same URL names.
		h.liteTenantsListEndpoint(w, r)
	// case n == 2 && p[0] == "v1" && p[1] == "tenants" && r.Method == http.MethodPost:
	// 	h.postCreateTenant(w, r)
	case n == 3 && p[0] == "v1" && p[1] == "franchise" && r.Method == http.MethodGet:
		h.tenantGetEndpoint(w, r, p[2])
	case n == 3 && p[0] == "v1" && p[1] == "franchise" && r.Method == http.MethodPut:
		h.tenantUpdateEndpoint(w, r, p[2])
	// case n == 3 && p[0] == "v1" && p[1] == "tenant" && r.Method == http.MethodDelete:
	// 	h.deleteTenantById(w, r, p[2])
	case n == 2 && p[0] == "v1" && p[1] == "login" && r.Method == http.MethodPost:
		h.loginEndpoint(w, r)
	case n == 2 && p[0] == "v1" && p[1] == "profile" && r.Method == http.MethodGet:
		h.profileEndpoint(w, r)
	default:
		http.NotFound(w, r)
	}
}
