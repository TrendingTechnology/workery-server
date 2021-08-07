package idos

import (
	//null "gopkg.in/guregu/null.v4"

	"github.com/over55/workery-server/internal/models"
)

type DashboardIDO struct {
    CustomerCount          uint64                      `json:"customer_count"`
	JobCount               uint64                      `json:"job_count"`
	MemberCount            uint64                      `json:"member_count"`
	TasksCount             uint64                      `json:"tasks_count"`
	BulletinBoardItems     []*models.BulletinBoardItem `json:"bulletin_board_items"`
	LastModifiedJobsByUser []*models.LiteWorkOrder         `json:"last_modified_jobs_by_user"`
	// AwayLog                []*models.AwayLog           `json:"away_log"` //TODO: IMPL.
	LastModifiedJobsByTeam []*models.LiteWorkOrder         `json:"last_modified_jobs_by_team"`
	PastFewDayComments     []*models.Comment           `json:"past_few_day_comments"`
}
