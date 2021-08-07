package controllers

import (
	// "encoding/json"
	"log"
	"net/http"
	"encoding/json"
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

    customerCountCh := make(chan uint64)
	go func() {
		f := &models.LiteCustomerFilter{
			TenantId: tenantId,
			States: []int8{models.CustomerActiveState},
		}
		count, err := h.LiteCustomerRepo.CountByFilter(ctx, f)
		if err != nil {
			log.Println("WARNING | dashboardEndpoint | h.LiteCustomerRepo.CountByFilter | err", err)
			return
		}
		customerCountCh <- count
	}()

	//
	// Find "job_count".
	//

    jobCountCh := make(chan uint64)
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
			return
		}
		jobCountCh <- count
	}()

	//
	// Find "member_count".
	//

    memberCountCh := make(chan uint64)
	go func() {
		f := &models.LiteAssociateFilter{
			TenantId: tenantId,
			States: []int8{models.AssociateActiveState},
		}
		count, err := h.LiteAssociateRepo.CountByFilter(ctx, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		memberCountCh <- count
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasksCountCh <- count
	}()

	//
	// Find "bulletin_board_items".
	//

	bulletinBoardItemsCh := make(chan []*models.BulletinBoardItem)
	go func() {
		f := &models.BulletinBoardItemFilter{
			TenantId: tenantId,
			States: []int8{models.BulletinBoardItemActiveState},
		}
		arr, err := h.BulletinBoardItemRepo.ListByFilter(ctx, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bulletinBoardItemsCh <- arr[:]
	}()

	//
	// Find "last_modified_jobs_by_user".
	//

	lmbuCh := make(chan []*models.LiteWorkOrder)
	go func() {
		f := &models.LiteWorkOrderFilter{
			TenantId: tenantId,
			LastModifiedById: null.IntFrom(int64(userId)),
			Limit: 5,
		}
		arr, err := h.LiteWorkOrderRepo.ListByFilter(ctx, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lmbuCh <- arr[:]
	}()

	//
	// Find "away_log".
	//
	//TODO: IMPL.

	//
	// Find "last_modified_jobs_by_team".
	//

	lmbtCh := make(chan []*models.LiteWorkOrder)
	go func() {
		f := &models.LiteWorkOrderFilter{
			TenantId: tenantId,
			Limit: 10,
		}
		arr, err := h.LiteWorkOrderRepo.ListByFilter(ctx, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lmbtCh <- arr[:]
	}()

	//
	// Find "past_few_day_comments".
	//



	//
	// Generate our response
	//

	customerCount, jobCount, memberCount, taskCount, bulletinBoardItems, lmbu, lmbt := <- customerCountCh, <- jobCountCh, <- memberCountCh, <- tasksCountCh, <- bulletinBoardItemsCh, <- lmbuCh, <- lmbtCh

	ido := &idos.DashboardIDO{
		CustomerCount: customerCount,
		JobCount: jobCount,
		MemberCount: memberCount,
		TasksCount: taskCount,
		BulletinBoardItems: bulletinBoardItems,
		LastModifiedJobsByUser: lmbu,
		LastModifiedJobsByTeam: lmbt,
	}
	log.Println(ido)
	if err := json.NewEncoder(w).Encode(&ido); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*

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
