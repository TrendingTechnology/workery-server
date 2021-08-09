package controllers

import (
	// "encoding/json"
	"encoding/json"
	"log"
	"net/http"
	// "time"

	// "github.com/google/uuid"
	null "gopkg.in/guregu/null.v4"

	"github.com/over55/workery-server/internal/models"
	// "github.com/over55/workery-server/internal/idos"
	"github.com/over55/workery-server/internal/idos"
)

func (h *Controller) dashboardEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantId := uint64(ctx.Value("user_tenant_id").(uint64))
	userId := uint64(ctx.Value("user_id").(uint64))

	//
	// Find "customer_count".
	//

	ccCh := make(chan uint64)
	go func() {
		f := &models.LiteCustomerFilter{
			TenantId: tenantId,
			States:   []int8{models.CustomerActiveState},
		}
		count, err := h.LiteCustomerRepo.CountByFilter(ctx, f)
		if err != nil {
			log.Println("WARNING | dashboardEndpoint | h.LiteCustomerRepo.CountByFilter | err", err)
			ccCh <- 0
		} else {
			ccCh <- count
		}

	}()

	//
	// Find "job_count".
	//

	jcCh := make(chan uint64)
	go func() {
		f := &models.LiteWorkOrderFilter{
			TenantId: tenantId,
			States: []int8{
				models.WorkOrderNewState,
				models.WorkOrderPendingState,
				models.WorkOrderOngoingState,
				models.WorkOrderInProgressState,
			},
		}
		count, err := h.LiteWorkOrderRepo.CountByFilter(ctx, f)
		if err != nil {
			log.Println("WARNING | dashboardEndpoint | h.LiteWorkOrderRepo.CountByFilter | err", err)
			jcCh <- 0
		} else {
			jcCh <- count
		}
	}()

	//
	// Find "member_count".
	//

	mcCh := make(chan uint64)
	go func() {
		f := &models.LiteAssociateFilter{
			TenantId: tenantId,
			States:   []int8{models.AssociateActiveState},
		}
		count, err := h.LiteAssociateRepo.CountByFilter(ctx, f)
		if err != nil {
			log.Println("WARNING | dashboardEndpoint | LiteAssociateRepo.CountByFilter|err:", err)
			mcCh <- 0
		} else {
			mcCh <- count
		}
	}()

	//
	// Find "tasks_count".
	//

	tasksCountCh := make(chan uint64)
	go func() {
		f := &models.LiteTaskFilter{
			TenantId: tenantId,
			IsClosed: null.BoolFrom(false),
		}
		count, err := h.LiteTaskRepo.CountByFilter(ctx, f)
		if err != nil {
			log.Println("WARNING | dashboardEndpoint | LiteTaskRepo.CountByFilter|err:", err)
			tasksCountCh <- 0
		} else {
			tasksCountCh <- count
		}
	}()

	//
	// Find "bulletin_board_items".
	//

	bbiCh := make(chan []*models.BulletinBoardItem)
	go func() {
		f := &models.BulletinBoardItemFilter{
			TenantId: tenantId,
			States:   []int8{models.BulletinBoardItemActiveState},
		}
		arr, err := h.BulletinBoardItemRepo.ListByFilter(ctx, f)
		if err != nil {
			log.Println("WARNING | dashboardEndpoint | BulletinBoardItemRepo.ListByFilter|err:", err)
			bbiCh <- []*models.BulletinBoardItem{}
		} else {
			bbiCh <- arr[:]
		}
	}()

	//
	// Find "last_modified_jobs_by_user".
	//

	lmbuCh := make(chan []*models.LiteWorkOrder)
	go func() {
		f := &models.LiteWorkOrderFilter{
			TenantId:         tenantId,
			LastModifiedById: null.IntFrom(int64(userId)),
			Limit:            5,
		}
		arr, err := h.LiteWorkOrderRepo.ListByFilter(ctx, f)
		if err != nil {
			log.Println("WARNING | dashboardEndpoint | LiteWorkOrderRepo.ListByFilter|err:", err)
			lmbuCh <- []*models.LiteWorkOrder{}
		} else {
			lmbuCh <- arr[:]
		}
	}()

	//
	// Find "away_log".
	//

	alCh := make(chan []*models.AssociateAwayLog)
	go func() {
		f := &models.AssociateAwayLogFilter{
			TenantId: tenantId,
			States:   []int8{models.AssociateAwayLogActiveState},
		}
		arr, err := h.AssociateAwayLogRepo.ListByFilter(ctx, f)
		if err != nil {
			log.Println("WARNING | dashboardEndpoint | AssociateAwayLogRepo.ListByFilter|err:", err)
			alCh <- []*models.AssociateAwayLog{}
		} else {
			alCh <- arr[:]
		}
	}()

	//
	// Find "last_modified_jobs_by_team".
	//

	lmbtCh := make(chan []*models.LiteWorkOrder)
	go func() {
		f := &models.LiteWorkOrderFilter{
			TenantId: tenantId,
			Limit:    10,
		}
		arr, err := h.LiteWorkOrderRepo.ListByFilter(ctx, f)
		if err != nil {
			log.Println("WARNING | dashboardEndpoint | LiteWorkOrderRepo.ListByFilter|err:", err)
			lmbtCh <- []*models.LiteWorkOrder{}
		} else {
			lmbtCh <- arr[:]
		}
	}()

	//
	// Find "past_few_day_comments".
	//

	wocCh := make(chan []*models.WorkOrderComment)
	go func() {
		// sevenDaysAgoTime := null.TimeFrom(time.Now().Add(-7*24*time.Hour)) // 7 days ago //TODO: UNCOMMENT WHEN READY!
		f := &models.WorkOrderCommentFilter{
			TenantId:    tenantId,
			// CreatedTime: null.TimeFrom(sevenDaysAgoTime), //TODO: UNCOMMENT WHEN READY!
			Limit:       10,
		}
		arr, err := h.WorkOrderCommentRepo.ListByFilter(ctx, f)
		if err != nil {
			log.Println("WARNING | dashboardEndpoint | WorkOrderCommentRepo.ListByFilter|err:", err)
			wocCh <- []*models.WorkOrderComment{}
		} else {
			wocCh <- arr[:]
		}
	}()

	//
	// Block this function until all the `goroutines` finish before proceeding further.
	//

	cc, jc, mc, tc, bbi, lmbu, lmbt, al, woc := <-ccCh, <-jcCh, <-mcCh, <-tasksCountCh, <-bbiCh, <-lmbuCh, <-lmbtCh, <-alCh, <-wocCh

	//
	// Generate our response
	//

	ido := &idos.DashboardIDO{
		CustomerCount:          cc,
		JobCount:               jc,
		MemberCount:            mc,
		TasksCount:             tc,
		BulletinBoardItems:     bbi,
		LastModifiedJobsByUser: lmbu,
		LastModifiedJobsByTeam: lmbt,
		AwayLog:                al,
		PastFewDayComments:     woc,
	}
	log.Println("PastFewDayComments:", woc)
	if err := json.NewEncoder(w).Encode(&ido); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*


    def to_associate_representation(self, user):
        associate = Associate.objects.get(owner=user)
        return {
            'balance_owing_amount': str(associate.balance_owing_amount.amount),
        }

    def to_representation(self, user):
        if user.is_associate():
            return self.to_associate_representation(user)
        else:
            return self.to_staff_representation(user)

*/

//------------//
// NAVIGATION //
//------------//

// class NavigationSerializer(serializers.Serializer):
//     def to_representation(self, user):
//         tasks_count = 0
//         if user.is_associate():
//             tasks_count = TaskItem.objects.filter(
//                 is_closed=False,
//                 job__associate__owner=user,
//             ).count()
//         else:
//             tasks_count = TaskItem.objects.filter(is_closed=False).count()
//         return {
//             "tasks_count": tasks_count,
//         }
