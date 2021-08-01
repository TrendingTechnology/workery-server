package controllers

import (
	// "encoding/json"
	// "log"
	"net/http"
	"encoding/json"
	// "time"
    //
	// "github.com/google/uuid"
    //
	// "github.com/over55/workery-server/internal/models"
	// "github.com/over55/workery-server/internal/idos"
	"github.com/over55/workery-server/internal/idos"
)

func (h *Controller) dashboardEndpoint(w http.ResponseWriter, req *http.Request) {
	ido := &idos.DashboardIDO{}
	if err := json.NewEncoder(w).Encode(&ido); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
# --- COUNTING ---
customer_count = Customer.objects.filter(
	state=Customer.CUSTOMER_STATE.ACTIVE
).count()

job_count = WorkOrder.objects.filter(
	Q(state=WORK_ORDER_STATE.NEW) |
	Q(state=WORK_ORDER_STATE.PENDING) |
	Q(state=WORK_ORDER_STATE.ONGOING) |
	Q(state=WORK_ORDER_STATE.IN_PROGRESS)
).count()

member_count = Associate.objects.filter(
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
	"customer_count": customer_count,
	"job_count": job_count,
	"member_count": member_count,
	"tasks_count": tasks_count,
	"bulletin_board_items": bbi_s.data,
	"last_modified_jobs_by_user": lmjbu_s.data,
	"away_log": away_log_s.data,
	"last_modified_jobs_by_team": lmjbt_s.data,
	"past_few_day_comments": c_s.data,
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
