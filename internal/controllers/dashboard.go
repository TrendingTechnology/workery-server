package controllers

import (
	// "encoding/json"
	"log"
	"net/http"
	"encoding/json"
	// "time"
    //
	// "github.com/google/uuid"
    //
	"github.com/over55/workery-server/internal/models"
	// "github.com/over55/workery-server/internal/idos"
	"github.com/over55/workery-server/internal/idos"
)

func (h *Controller) dashboardEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantId := uint64(ctx.Value("user_tenant_id").(uint64))

    //
	// Find "customer_count".
	//

    customerCountCh := make(chan uint64)
	go func() {
		lcf := &models.LiteCustomerFilter{
			TenantId: tenantId,
			States: []int8{models.CustomerActiveState},
		}
		customerCount, err := h.LiteCustomerRepo.CountByFilter(ctx, lcf)
		if err != nil {
			log.Println("WARNING | dashboardEndpoint | h.LiteCustomerRepo.CountByFilter | err", err)
			return
		}
		customerCountCh <- customerCount
	}()

	//
	// Find "job_count".
	//

    jobCountCh := make(chan uint64)
	go func() {
		lwof := &models.LiteWorkOrderFilter{
			TenantId: tenantId,
			States: []int8{
				models.WorkOrderNewState,
				models.WorkOrderPendingState,
				models.WorkOrderOngoingState,
				models.WorkOrderInProgressState,
			},
		}
		workOrderCount, err := h.LiteWorkOrderRepo.CountByFilter(ctx, lwof)
		if err != nil {
			log.Println("WARNING | dashboardEndpoint | h.LiteWorkOrderRepo.CountByFilter | err", err)
			return
		}
		jobCountCh <- workOrderCount
	}()

	//
	// Find "member_count".
	//

    memberCountCh := make(chan uint64)
	go func() {
		laf := &models.LiteAssociateFilter{
			TenantId: tenantId,
			States: []int8{models.AssociateActiveState},
		}
		memberCount, err := h.LiteAssociateRepo.CountByFilter(ctx, laf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		memberCountCh <- memberCount
	}()

	//
	// TODO
	//
	// "tasks_count": tasks_count,                 // TODO
	// "bulletin_board_items": bbi_s.data,         // TODO
	// "last_modified_jobs_by_user": lmjbu_s.data, // TODO
	// "away_log": away_log_s.data,                // TODO
	// "last_modified_jobs_by_team": lmjbt_s.data, // TODO
	// "past_few_day_comments": c_s.data,          // TODO

	//
	// Generate our response
	//

	customerCount, jobCount, memberCount := <- customerCountCh, <- jobCountCh, <- memberCountCh

	ido := &idos.DashboardIDO{
		CustomerCount: customerCount,
		JobCount: jobCount,
		MemberCount: memberCount,
	}
	log.Println(ido)
	if err := json.NewEncoder(w).Encode(&ido); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
# --- COUNTING ---
customer_count = Customer.objects.filter(  // DONE
	state=Customer.CUSTOMER_STATE.ACTIVE
).count()

job_count = WorkOrder.objects.filter( // DONE
	Q(state=WORK_ORDER_STATE.NEW) |
	Q(state=WORK_ORDER_STATE.PENDING) |
	Q(state=WORK_ORDER_STATE.ONGOING) |
	Q(state=WORK_ORDER_STATE.IN_PROGRESS)
).count()

member_count = Associate.objects.filter( // DONE
	owner__is_active=True
).count()

tasks_count = TaskItem.objects.filter(
	is_closed=False
).count()

# --- BULLETIN BOARD ITEMS ---
bulletin_board_items = BulletinBoardItem.objects.filter(
	is_archived=False
).order_by(
	'-created_at'
).prefetch_related(
	'created_by'
)
bbi_s = BulletinBoardItemListCreateSerializer(bulletin_board_items, many=True)

# --- LATEST JOBS BY USER ---
last_modified_jobs_by_user = WorkOrder.objects.filter(
	last_modified_by = user
).order_by(
	'-last_modified'
).prefetch_related(
	'associate',
	'customer'
)[0:5]
lmjbu_s = WorkOrderListCreateSerializer(last_modified_jobs_by_user, many=True)

# --- LATEST JOBS BY TEAM ---
last_modified_jobs_by_team = WorkOrder.objects.order_by(
	'-last_modified'
).prefetch_related(
	'associate',
	'customer'
)[0:10]
lmjbt_s = WorkOrderListCreateSerializer(last_modified_jobs_by_team, many=True)

# --- ASSOCIATE AWAY LOGS ---
away_log = AwayLog.objects.filter(
	was_deleted=False
).prefetch_related(
	'associate'
)
away_log_s = AwayLogListCreateSerializer(away_log, many=True)

# --- LATEST AWAY COMMENT ---
one_week_before_today = get_todays_date_minus_days(5)
past_few_day_comments = WorkOrderComment.objects.filter(
	created_at__gte=one_week_before_today
).order_by(
	'-created_at'
).prefetch_related(
	'about',
	'comment'
)
c_s = WorkOrderCommentListCreateSerializer(past_few_day_comments, many=True)

return {
	"customer_count": customer_count, // DONE
	"job_count": job_count,           // DONE
	"member_count": member_count,     // DONE
	"tasks_count": tasks_count,                 // TODO
	"bulletin_board_items": bbi_s.data,         // TODO
	"last_modified_jobs_by_user": lmjbu_s.data, // TODO
	"away_log": away_log_s.data,                // TODO
	"last_modified_jobs_by_team": lmjbt_s.data, // TODO
	"past_few_day_comments": c_s.data,          // TODO
}
/

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
