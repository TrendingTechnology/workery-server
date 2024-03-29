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
	AssociateAwayLogRepo              models.AssociateAwayLogRepository
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
	LiteAssociateAwayLogRepo          models.LiteAssociateAwayLogRepository
	LiteAssociateRepo                 models.LiteAssociateRepository
	LiteBulletinBoardItemRepo         models.LiteBulletinBoardItemRepository
	LiteCustomerRepo                  models.LiteCustomerRepository
	LiteDeactivatedCustomerRepo       models.LiteDeactivatedCustomerRepository
	LiteFinancialRepo                 models.LiteFinancialRepository
	LiteHowHearAboutUsItemRepo        models.LiteHowHearAboutUsItemRepository
	LiteInsuranceRequirementRepo      models.LiteInsuranceRequirementRepository
	LitePartnerRepo                   models.LitePartnerRepository
	LiteSkillSetRepo                  models.LiteSkillSetRepository
	LiteStaffRepo                     models.LiteStaffRepository
	LiteTagRepo                       models.LiteTagRepository
	LiteTaskItemRepo                  models.LiteTaskItemRepository
	LiteTenantRepo                    models.LiteTenantRepository
	LiteVehicleTypeRepo               models.LiteVehicleTypeRepository
	LiteWorkOrderRepo                 models.LiteWorkOrderRepository
	LiteWorkOrderServiceFeeRepo       models.LiteWorkOrderServiceFeeRepository
	LiteOngoingWorkOrderRepo          models.LiteOngoingWorkOrderRepository
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
	// --- TENANTS ---
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

	// --- GATEWAY & PROFILE & DASHBOARD ---
	case n == 2 && p[0] == "v1" && p[1] == "login" && r.Method == http.MethodPost:
		h.loginEndpoint(w, r)
	case n == 2 && p[0] == "v1" && p[1] == "profile" && r.Method == http.MethodGet:
		h.profileEndpoint(w, r)
	case n == 2 && p[0] == "v1" && p[1] == "dashboard" && r.Method == http.MethodGet:
		h.dashboardEndpoint(w, r)
	case n == 2 && p[0] == "v1" && p[1] == "navigation" && r.Method == http.MethodGet:
		h.navigationEndpoint(w, r)

	// --- CUSTOMERS ---
	case n == 2 && p[0] == "v1" && p[1] == "customers" && r.Method == http.MethodGet:
		h.customersListEndpoint(w, r)
	case n == 3 && p[0] == "v1" && p[1] == "customer" && r.Method == http.MethodGet:
		h.customerGetEndpoint(w, r, p[2])

	// --- WORK ORDERS ---
	case n == 2 && p[0] == "v1" && p[1] == "orders" && r.Method == http.MethodGet:
		h.workOrdersListEndpoint(w, r)

	// --- ASSOCIATES ---
	case n == 2 && p[0] == "v1" && p[1] == "associates" && r.Method == http.MethodGet:
		h.associatesListEndpoint(w, r)

	// --- TASKS ---
	case n == 2 && p[0] == "v1" && p[1] == "tasks" && r.Method == http.MethodGet:
		h.taskItemsListEndpoint(w, r)

	// --- ONGOING WORK ORDERS ---
	case n == 2 && p[0] == "v1" && p[1] == "ongoing-orders" && r.Method == http.MethodGet:
		h.ongoingWorkOrdersListEndpoint(w, r)

		// --- PARTNERS ---
	case n == 2 && p[0] == "v1" && p[1] == "partners" && r.Method == http.MethodGet:
		h.partnersListEndpoint(w, r)

		// --- STAFF ---
	case n == 2 && p[0] == "v1" && p[1] == "staff" && r.Method == http.MethodGet:
		h.staffListEndpoint(w, r)

		// --- FINANCIALS ---
	case n == 2 && p[0] == "v1" && p[1] == "financials" && r.Method == http.MethodGet:
		h.financialsListEndpoint(w, r)

		// --- BULLETIN BOARD ITEMS ---
	case n == 2 && p[0] == "v1" && p[1] == "bulletin-board-items" && r.Method == http.MethodGet:
		h.bulletinBoardItemsListEndpoint(w, r)

		// --- SKILL SETS ---
	case n == 2 && p[0] == "v1" && p[1] == "skill-sets" && r.Method == http.MethodGet:
		h.skillSetsListEndpoint(w, r)

		// --- TAGS ---
	case n == 2 && p[0] == "v1" && p[1] == "tags" && r.Method == http.MethodGet:
		h.tagsListEndpoint(w, r)

	// --- ASSOCIATE AWAY LOGS ---
	case n == 2 && p[0] == "v1" && p[1] == "associate-away-logs" && r.Method == http.MethodGet:
		h.associateAwayLogsListEndpoint(w, r)

	// --- INSURANCE REQUIREMENTS ---
	case n == 2 && p[0] == "v1" && p[1] == "insurance-requirements" && r.Method == http.MethodGet:
		h.insuranceRequirementsListEndpoint(w, r)

	// --- WORK ORDER SERVICE FEES ---
	case n == 2 && p[0] == "v1" && p[1] == "order-service-fees" && r.Method == http.MethodGet:
		h.workOrderServiceFeesListEndpoint(w, r)

	// --- DEACTIVATED CUSTOMER ---
	case n == 2 && p[0] == "v1" && p[1] == "deactivated-customers" && r.Method == http.MethodGet:
		h.deactivatedCustomersListEndpoint(w, r)

	// --- DEACTIVATED CUSTOMER ---
	case n == 2 && p[0] == "v1" && p[1] == "vehicle-types" && r.Method == http.MethodGet:
		h.vehicleTypesListEndpoint(w, r)

		// --- HOW HEAR ABOUT US ITEM ---
	case n == 2 && p[0] == "v1" && p[1] == "how-hears" && r.Method == http.MethodGet:
		h.howHearAboutUsItemsListEndpoint(w, r)

	// --- CATCH ALL: D.N.E. ---
	default:
		http.NotFound(w, r)
	}
}

/*
# Away logs.
url(r'^api/away-logs$', AwayLogListCreateAPIView.as_view(), name='workery_away_log_list_create_api_endpoint'),
url(r'^api/away-log/(?P<pk>[^/.]+)/$', AwayLogRetrieveUpdateDestroyAPIView.as_view(), name='workery_away_log_retrieve_update_destroy_api_endpoint'),

# Associates
url(r'^api/associates$', AssociateListCreateAPIView.as_view(), name='workery_associate_list_create_api_endpoint'),
url(r'^api/associates/validate$', AssociateCreateValidationAPIView.as_view(), name='workery_associate_create_validate_api_endpoint'),
url(r'^api/associate-files$', AssociateFileUploadListCreateAPIView.as_view(), name='workery_associate_file_upload_api_endpoint'),
url(r'^api/associate-file/(?P<pk>[^/.]+)/$', AssociateFileUploadArchiveAPIView.as_view(), name='workery_associate_file_upload_archive_api_endpoint'),
url(r'^api/associate/(?P<pk>[^/.]+)/contact$', AssociateContactUpdateAPIView.as_view(), name='workery_associate_contact_update_api_endpoint'),
url(r'^api/associate/(?P<pk>[^/.]+)/address$', AssociateAddressUpdateAPIView.as_view(), name='workery_associate_address_update_api_endpoint'),
url(r'^api/associate/(?P<pk>[^/.]+)/account$', AssociateAccountUpdateAPIView.as_view(), name='workery_associate_account_update_api_endpoint'),
url(r'^api/associate/(?P<pk>[^/.]+)/metrics$', AssociateMetricsUpdateAPIView.as_view(), name='workery_associate_metrics_update_api_endpoint'),
url(r'^api/associate/(?P<pk>[^/.]+)/$', AssociateRetrieveUpdateDestroyAPIView.as_view(), name='workery_associate_retrieve_update_destroy_api_endpoint'),
url(r'^api/associate-comments$', AssociateCommentListCreateAPIView.as_view(), name='workery_associate_comment_list_create_api_endpoint'),
url(r'^api/associates/operation/avatar$', AssociateAvatarCreateOrUpdateOperationCreateAPIView.as_view(), name='workery_associate_avatar_operation_api_endpoint'),
url(r'^api/associates/operation/balance$', AssociateBalanceOperationAPIView.as_view(), name='workery_associate_balance_operation_api_endpoint'),
url(r'^api/associates/operation/password$', AssociateChangePasswordOperationAPIView.as_view(), name='workery_associate_password_operation_api_endpoint'),
url(r'^api/associates/operation/upgrade$', AssociateUprageOperationAPIView.as_view(), name='workery_associate_upgrade_operation_api_endpoint'),
url(r'^api/associates/operation/downgrade$', AssociateDowngradeOperationAPIView.as_view(), name='workery_associate_downgrade_operation_api_endpoint'),
url(r'^api/associates/operation/archive$', AssociateArchiveOperationCreateAPIView.as_view(), name='workery_associate_archive_operation_create_api_endpoint'),

# Customers
url(r'^api/customers$', CustomerListCreateAPIView.as_view(), name='workery_customer_list_create_api_endpoint'),
url(r'^api/v2/customers$', CustomerListCreateV2APIView.as_view(), name='workery_customer_list_create_v2_api_endpoint'),
url(r'^api/customers/validate$', CustomerCreateValidationAPIView.as_view(), name='workery_customer_create_validate_api_endpoint'),
url(r'^api/customer-files$', CustomerFileUploadListCreateAPIView.as_view(), name='workery_customer_file_upload_api_endpoint'),
url(r'^api/customer-file/(?P<pk>[^/.]+)/$', CustomerFileUploadArchiveAPIView.as_view(), name='workery_customer_file_upload_archive_api_endpoint'),
url(r'^api/customer/(?P<pk>[^/.]+)/$', CustomerRetrieveUpdateDestroyAPIView.as_view(), name='workery_customer_retrieve_update_destroy_api_endpoint'),
url(r'^api/v2/customer/(?P<pk>[^/.]+)/$', CustomerRetrieveUpdateDestroyV2APIView.as_view(), name='workery_customer_retrieve_update_destroy_v2_api_endpoint'),
url(r'^api/customer/(?P<pk>[^/.]+)/contact$', CustomerContactUpdateAPIView.as_view(), name='workery_customer_contact_update_api_endpoint'),
url(r'^api/customer/(?P<pk>[^/.]+)/address$', CustomerAddressUpdateAPIView.as_view(), name='workery_customer_address_update_api_endpoint'),
url(r'^api/customer/(?P<pk>[^/.]+)/metrics$', CustomerMetricsUpdateAPIView.as_view(), name='workery_customer_metrics_update_api_endpoint'),
url(r'^api/customer-comments$', CustomerCommentListCreateAPIView.as_view(), name='workery_customer_comment_list_create_api_endpoint'),
url(r'^api/deactivated-customers$', DeactivatedCustomerListAPIView.as_view(), name='workery_deactivated_customer_list_api_endpoint'),

# Customers - Operations
url(r'^api/customers/operation/archive$', CustomerArchiveOperationCreateAPIView.as_view(), name='workery_archive_customer_operation_create_api_endpoint'),
url(r'^api/customers/operation/upgrade-residential$', ResidentialCustomerUpgradeOperationCreateAPIView.as_view(), name='workery_residential_customer_upgrade_operation_api_endpoint'),
url(r'^api/customers/operation/avatar$', CustomerAvatarCreateOrUpdateOperationCreateAPIView.as_view(), name='workery_avatar_customer_operation_create_or_update_api_endpoint'),

# Insurance Requirements
url(r'^api/insurance_requirements$', InsuranceRequirementListCreateAPIView.as_view(), name='workery_insurance_requirement_list_create_api_endpoint'),
url(r'^api/insurance_requirement/(?P<pk>[^/.]+)/$', InsuranceRequirementRetrieveUpdateDestroyAPIView.as_view(), name='workery_insurance_requirement_retrieve_update_destroy_api_endpoint'),

# Public Image Uploads.
url(r'^api/public-image-uploads$', PublicImageUploadListCreateAPIView.as_view(), name='workery_public_image_upload_list_create_api_endpoint'),

# WorkOrders
url(r'^api/orders$', WorkOrderListCreateAPIView.as_view(), name='workery_order_list_create_api_endpoint'),
url(r'^api/order/(?P<pk>[^/.]+)/$', WorkOrderRetrieveUpdateDestroyAPIView.as_view(), name='workery_order_retrieve_update_destroy_api_endpoint'),
url(r'^api/order/(?P<pk>[^/.]+)/lite$', WorkOrderLiteUpdateAPIView.as_view(), name='workery_order_lite_update_api_endpoint'),
url(r'^api/order/(?P<pk>[^/.]+)/financial$', WorkOrderFinancialUpdateAPIView.as_view(), name='workery_order_financial_update_api_endpoint'),
url(r'^api/order/(?P<pk>[^/.]+)/invoice$', WorkOrderInvoiceRetrieveAPIView.as_view(), name='workery_order_invoice_retrieve_api_endpoint'),
# url(r'^api/order/(?P<pk>[^/.]+)/invoice/first-section$', WorkOrderInvoiceFirstSectionUpdateAPIView.as_view(), name='workery_order_invoice_first_section_update_api_endpoint'), # DEPRECATED
url(r'^api/order/(?P<pk>[^/.]+)/invoice/second-section$', WorkOrderInvoiceSecondSectionUpdateAPIView.as_view(), name='workery_order_invoice_second_section_update_api_endpoint'),
url(r'^api/order/(?P<pk>[^/.]+)/invoice/third-section$', WorkOrderInvoiceThirdSectionUpdateAPIView.as_view(), name='workery_order_invoice_third_section_update_api_endpoint'),
url(r'^api/order/(?P<pk>[^/.]+)/download-invoice-pdf$', WorkOrderInvoiceDownloadPDFAPIView.as_view(), name='workery_order_invoice_pdf_download_api_endpoint'),
url(r'^api/order/(?P<pk>[^/.]+)/deposits$', WorkOrderDepositListCreateAPIView.as_view(), name='workery_order_deposit_list_create_api_endpoint'),
url(r'^api/order/(?P<order_pk>[^/.]+)/deposit/(?P<payment_pk>[^/.]+)$', WorkOrderDepositDeleteAPIView.as_view(), name='workery_order_deposit_delete_api_endpoint'),
url(r'^api/order-comments$', WorkOrderCommentListCreateAPIView.as_view(), name='workery_job_comment_list_create_api_endpoint'),
url(r'^api/order-files$', WorkOrderFileUploadListCreateAPIView.as_view(), name='workery_work_order_file_upload_api_endpoint'),
url(r'^api/order-file/(?P<pk>[^/.]+)/$', WorkOrderFileUploadArchiveAPIView.as_view(), name='workery_work_order_file_upload_archive_api_endpoint'),
url(r'^api/my-orders$',MyWorkOrderListAPIView.as_view(), name='workery_my_orders_list_api_endpoint'),
url(r'^api/my-order/(?P<pk>[^/.]+)/$', MyWorkOrderRetrieveAPIView.as_view(), name='workery_my_order_retrieve_api_endpoint'),

# WorkOrder - Operations
url(r'^api/orders/operation/unassign$', WorkOrderUnassignOperationCreateAPIView.as_view(), name='workery_order_unassign_operation_api_endpoint'),
url(r'^api/orders/operation/clone$', WorkOrderCloneOperationCreateAPIView.as_view(), name='workery_order_clone_operation_api_endpoint'),
url(r'^api/orders/operation/close$', WorkOrderCloseOperationCreateAPIView.as_view(), name='workery_order_close_operation_api_endpoint'),
url(r'^api/orders/operation/postpone$', WorkOrderPostponeOperationCreateAPIView.as_view(), name='workery_order_postpone_operation_api_endpoint'),
url(r'^api/orders/operation/reopen$', WorkOrderReopenOperationCreateAPIView.as_view(), name='workery_order_reopen_operation_api_endpoint'),                           #TODO: DELETE
url(r'^api/orders/operation/transfer$', TransferWorkerOrderOperationAPIView.as_view(), name='workery_transfer_order_operation_api_endpoint'),                           #TODO: DELETE
url(r'^api/orders/operation/invoice$', WorkOrderInvoiceCreateOrUpdateOperationAPIView.as_view(), name='workery_order_invoice_create_or_update_operation_api_endpoint'),                           #TODO: DELETE


# Work Order Service Fees
url(r'^api/order_service_fees$', WorkOrderServiceFeeListCreateAPIView.as_view(), name='workery_order_service_fee_list_create_api_endpoint'),
url(r'^api/order_service_fee/(?P<pk>[^/.]+)/$', WorkOrderServiceFeeRetrieveUpdateDestroyAPIView.as_view(), name='workery_order_service_fee_retrieve_update_destroy_api_endpoint'),

# Partners
url(r'^api/partners$', PartnerListCreateAPIView.as_view(), name='workery_partner_list_create_api_endpoint'),
url(r'^api/v2/partners$', PartnerListCreateV2APIView.as_view(), name='workery_partner_list_create_api_v2_endpoint'),
url(r'^api/partners/validate$', PartnerCreateValidationAPIView.as_view(), name='workery_partner_create_validate_api_endpoint'),
url(r'^api/partner/(?P<pk>[^/.]+)/$', PartnerRetrieveUpdateDestroyAPIView.as_view(), name='workery_partner_retrieve_update_destroy_api_endpoint'),
url(r'^api/v2/partner/(?P<pk>[^/.]+)/$', PartnerRetrieveUpdateDestroyV2APIView.as_view(), name='workery_partner_retrieve_update_destroy_api_v2_endpoint'),
url(r'^api/partner/(?P<pk>[^/.]+)/contact$', PartnerContactUpdateAPIView.as_view(), name='workery_partner_contact_update_api_endpoint'),
url(r'^api/partner/(?P<pk>[^/.]+)/address$', PartnerAddressUpdateAPIView.as_view(), name='workery_partner_address_update_api_endpoint'),
url(r'^api/partner/(?P<pk>[^/.]+)/metrics$', PartnerMetricsUpdateAPIView.as_view(), name='workery_partner_metrics_update_api_endpoint'),
url(r'^api/partner-comments$', PartnerCommentListCreateAPIView.as_view(), name='workery_partner_comment_list_create_api_endpoint'),
url(r'^api/partner-files$', PartnerFileUploadListCreateAPIView.as_view(), name='workery_partner_file_upload_api_endpoint'),
url(r'^api/partner-file/(?P<pk>[^/.]+)/$', PartnerFileUploadArchiveAPIView.as_view(), name='workery_partner_file_upload_archive_api_endpoint'),
url(r'^api/partners/operation/avatar$', PartnerAvatarOperationAPIView.as_view(), name='workery_partner_avatar_operation_api_endpoint'),

# Skill Sets
url(r'^api/skill_sets$', SkillSetListCreateAPIView.as_view(), name='workery_skill_set_list_create_api_endpoint'),
url(r'^api/skill_set/(?P<pk>[^/.]+)/$', SkillSetRetrieveUpdateDestroyAPIView.as_view(), name='workery_skill_set_retrieve_update_destroy_api_endpoint'),

# Staff
url(r'^api/staves$', StaffListCreateAPIView.as_view(), name='workery_staff_list_create_api_endpoint'),
url(r'^api/staves/validate$', StaffCreateValidationAPIView.as_view(), name='workery_staff_create_validate_api_endpoint'),
url(r'^api/staff/(?P<pk>[^/.]+)/$', StaffRetrieveUpdateDestroyAPIView.as_view(), name='workery_staff_retrieve_update_destroy_api_endpoint'),
url(r'^api/staff/(?P<pk>[^/.]+)/contact$', StaffContactUpdateAPIView.as_view(), name='workery_staff_contact_update_api_endpoint'),
url(r'^api/staff/(?P<pk>[^/.]+)/address$', StaffAddressUpdateAPIView.as_view(), name='workery_staff_address_update_api_endpoint'),
url(r'^api/staff/(?P<pk>[^/.]+)/account$', StaffAccountUpdateAPIView.as_view(), name='workery_staff_account_update_api_endpoint'),
url(r'^api/staff/(?P<pk>[^/.]+)/metrics$', StaffMetricsUpdateAPIView.as_view(), name='workery_staff_metrics_update_api_endpoint'),
url(r'^api/staff-comments$', StaffCommentListCreateAPIView.as_view(), name='workery_staff_comment_list_create_api_endpoint'),
url(r'^api/staff-files$', StaffFileUploadListCreateAPIView.as_view(), name='workery_staff_file_upload_api_endpoint'),
url(r'^api/staff-file/(?P<pk>[^/.]+)/$', StaffFileUploadArchiveAPIView.as_view(), name='workery_staff_file_upload_archive_api_endpoint'),
url(r'^api/v2/staves$', StaffListCreateV2APIView.as_view(), name='workery_v2_staff_list_create_api_endpoint'),
url(r'^api/v2/staff/(?P<pk>[^/.]+)/$', StaffRetrieveAPIView.as_view(), name='workery_v2_staff_retrieve_api_endpoint'),

# Staff Operations
url(r'^api/staff/(?P<pk>[^/.]+)/archive$', StaffArchiveAPIView.as_view(), name='workery_staff_archive_api_endpoint'),
url(r'^api/staff/(?P<pk>[^/.]+)/change-role$', StaffChangeRoleOperationAPIView.as_view(), name='workery_staff_change_role_operation_api_endpoint'),
url(r'^api/staff/(?P<pk>[^/.]+)/change-password$', StaffChangePasswordOperationAPIView.as_view(), name='workery_staff_change_password_operation_api_endpoint'),
url(r'^api/staff/operation/avatar$', StaffAvatarOperationCreateAPIView.as_view(), name='workery_staff_avatar_operation_api_endpoint'),

# Tags
url(r'^api/tags$', TagListCreateAPIView.as_view(), name='workery_tag_list_create_api_endpoint'),
url(r'^api/tag/(?P<pk>[^/.]+)/$', TagRetrieveUpdateDestroyAPIView.as_view(), name='workery_tag_retrieve_update_destroy_api_endpoint'),

# Vehicle Type
url(r'^api/vehicle_types$', VehicleTypeListCreateAPIView.as_view(), name='workery_vehicle_type_list_create_api_endpoint'),
url(r'^api/vehicle_type/(?P<pk>[^/.]+)/$', VehicleTypeRetrieveUpdateDestroyAPIView.as_view(), name='workery_vehicle_type_retrieve_update_destroy_api_endpoint'),

# HowHearAboutUsItems
url(r'^api/how_hears$', HowHearAboutUsItemListCreateAPIView.as_view(), name='workery_how_hear_list_create_api_endpoint'),
url(r'^api/how_hear/(?P<pk>[^/.]+)/$', HowHearAboutUsItemRetrieveUpdateDestroyAPIView.as_view(), name='workery_how_hear_retrieve_update_destroy_api_endpoint'),

# Utility
url(r'^api/utility/find-customer-matching$', FindCustomerMatchingAPIView.as_view(), name='workery_find_customer_matching_api_endpoint'),

# Ongoing Work Order
url(r'^api/ongoing-orders$', OngoingWorkOrderListCreateAPIView.as_view(), name='workery_ongoing_order_list_create_api_endpoint'),
url(r'^api/ongoing-order/(?P<pk>[^/.]+)/$', OngoingWorkOrderRetrieveUpdateDestroyAPIView.as_view(), name='workery_ongoing_order_retrieve_update_destroy_api_endpoint'),
url(r'^api/ongoing-order-comments$', OngoingWorkOrderCommentListCreateAPIView.as_view(), name='workery_ongoing_job_comment_list_create_api_endpoint'),
url(r'^api/v2/ongoing-orders$', OngoingWorkOrderListCreateV2APIView.as_view(), name='workery_ongoing_order_list_create_v2_api_endpoint'),
url(r'^api/v2/ongoing-order/(?P<pk>[^/.]+)/$', OngoingWorkOrderRetrieveUpdateDestroyV2APIView.as_view(), name='workery_ongoing_order_retrieve_update_destroy_v2_api_endpoint'),

# Bulletin Board Items
url(r'^api/bulletin_board_items$', BulletinBoardItemListCreateAPIView.as_view(), name='workery_bulletin_board_item_list_create_api_endpoint'),
url(r'^api/bulletin_board_item/(?P<pk>[^/.]+)/$', BulletinBoardItemRetrieveUpdateDestroyAPIView.as_view(), name='workery_bulletin_board_item_retrieve_update_destroy_api_endpoint'),

# Tasks
url(r'^api/tasks$', TaskItemListPIView.as_view(), name='workery_task_item_list_api_endpoint'),
url(r'^api/task/(?P<pk>[^/.]+)/$', TaskItemRetrieveAPIView.as_view(), name='workery_task_item_retrieve_api_endpoint'),
url(r'^api/task/(?P<pk>[^/.]+)/available-associates$', TaskItemAvailableAssociateListCreateAPIView.as_view(), name='workery_task_item_available_associate_list_create_api_endpoint'),

# Tasks - Operation
url(r'^api/task/operation/assign-associate$', AssignAssociateTaskOperationAPIView.as_view(), name='workery_order_task_operation_assign_associate_api_endpoint'),
url(r'^api/task/operation/follow-up$', FollowUpTaskOperationAPIView.as_view(), name='workery_task_operation_follow_up_create_api_endpoint'),
url(r'^api/orders/complete$', FollowUpTaskOperationAPIView.as_view(), name='workery_order_order_complete_create_api_endpoint'),
url(r'^api/task/operation/follow-up-pending$', FollowUpPendingTaskOperationAPIView.as_view(), name='workery_order_task_operation_follow_up_pending_api_endpoint'),
url(r'^api/task/operation/close$', CloseTaskOperationAPIView.as_view(), name='workery_task_operation_close_api_endpoint'), #TODO: Integrate `CloseTaskOperationAPIView` with current close out view.
url(r'^api/v2/task/operation/follow-up$', FollowUpTaskOperationV2APIView.as_view(), name='workery_task_operation_follow_up_create_v2_api_endpoint'),
url(r'^api/v2/task/operation/follow-up-pending$', FollowUpPendingTaskOperationV2APIView.as_view(), name='workery_order_task_operation_follow_up_pending_v2_api_endpoint'),
url(r'^api/task/operation/order-completion$', OrderCompletionTaskOperationAPIView.as_view(), name='workery_task_operation_order_completion_api_endpoint'),
url(r'^api/task/operation/survey$', SurveyTaskOperationAPIView.as_view(), name='workery_task_operation_survey_api_endpoint'),

# ActivitySheetItem
url(r'^api/activity-sheets$', ActivitySheetItemListCreateAPIView.as_view(), name='workery_activity_sheet_list_create_api_endpoint'),
url(r'^api/activity-sheet/(?P<pk>[^/.]+)/$', ActivitySheetItemRetrieveUpdateDestroyAPIView.as_view(), name='workery_activity_sheet_retrieve_update_destroy_api_endpoint'),

# Search
url(r'^api/v1/search$', UnifiedSearchItemListAPIView.as_view(), name='workery_unified_search_list_api_endpoint'),

*/
