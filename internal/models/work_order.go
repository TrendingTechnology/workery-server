package models

import (
	"context"
	"time"

	null "gopkg.in/guregu/null.v4"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

// '''
// description text COLLATE pg_catalog."default",
//     assignment_date date,
//     is_ongoing boolean NOT NULL,
//     is_home_support_service boolean NOT NULL,
//     start_date date NOT NULL,
//     completion_date date,
//     hours numeric(7,1),
//     type_of smallint NOT NULL,
//     indexed_text character varying(2047) COLLATE pg_catalog."default",
//     closing_reason smallint,
//     closing_reason_other character varying(1024) COLLATE pg_catalog."default",
//     state character varying(50) COLLATE pg_catalog."default" NOT NULL,
//     was_job_satisfactory boolean NOT NULL,
//     was_job_finished_on_time_and_on_budget boolean NOT NULL,
//     was_associate_punctual boolean NOT NULL,
//     was_associate_professional boolean NOT NULL,
//     would_customer_refer_our_organization boolean NOT NULL,
//     score smallint NOT NULL,
//     invoice_date date,
//     invoice_quote_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_quote_amount numeric(10,2) NOT NULL,
//     invoice_labour_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_labour_amount numeric(10,2) NOT NULL,
//     invoice_material_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_material_amount numeric(10,2) NOT NULL,
//     invoice_tax_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_tax_amount numeric(10,2) NOT NULL,
//     invoice_total_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_total_amount numeric(10,2) NOT NULL,
//     invoice_service_fee_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_service_fee_amount numeric(10,2) NOT NULL,
//     invoice_service_fee_payment_date date,
//     created timestamp with time zone NOT NULL,
//     created_from inet,
//     created_from_is_public boolean NOT NULL,
//     last_modified timestamp with time zone NOT NULL,
//     last_modified_from inet,
//     last_modified_from_is_public boolean NOT NULL,
//     associate_id bigint,
//     created_by_id integer,
//     customer_id bigint NOT NULL,
//     invoice_service_fee_id bigint,
//     last_modified_by_id integer,
//     latest_pending_task_id bigint,
//     ongoing_work_order_id bigint,
//     was_survey_conducted boolean NOT NULL,
//     was_there_financials_inputted boolean NOT NULL,
//     invoice_actual_service_fee_amount_paid numeric(10,2) NOT NULL,
//     invoice_actual_service_fee_amount_paid_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_balance_owing_amount numeric(10,2) NOT NULL,
//     invoice_balance_owing_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_quoted_labour_amount numeric(10,2) NOT NULL,
//     invoice_quoted_labour_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_quoted_material_amount numeric(10,2) NOT NULL,
//     invoice_quoted_material_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_total_quote_amount numeric(10,2) NOT NULL,
//     invoice_total_quote_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     visits smallint NOT NULL,
//     invoice_ids character varying(127) COLLATE pg_catalog."default",
//     no_survey_conducted_reason smallint,
//     no_survey_conducted_reason_other character varying(1024) COLLATE pg_catalog."default",
//     cloned_from_id bigint,
//     invoice_deposit_amount numeric(10,2) NOT NULL,
//     invoice_deposit_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_other_costs_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_quoted_other_costs_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_paid_to smallint,
//     invoice_amount_due numeric(10,2) NOT NULL,
//     invoice_amount_due_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_sub_total_amount numeric(10,2) NOT NULL,
//     invoice_sub_total_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_other_costs_amount numeric(10,2) NOT NULL,
//     invoice_quoted_other_costs_amount numeric(10,2) NOT NULL,
//     closing_reason_comment character varying(1024) COLLATE pg_catalog."default",
// '''

type WorkOrder struct {
	Id                 uint64      `json:"id"`
	Uuid               string      `json:"uuid"`
	TenantId           uint64      `json:"tenant_id"`
	CustomerId         uint64      `json:"customer_id"`
	AssociateId        null.Int    `json:"associate_id"`
	State              int8        `json:"state"`
	CreatedTime        time.Time   `json:"created_time"`
	CreatedById        null.Int    `json:"created_by_id"`
	CreatedFromIP      null.String `json:"created_from_ip"`
	LastModifiedTime   time.Time   `json:"last_modified_time"`
	LastModifiedById   null.Int    `json:"last_modified_by_id"`
	LastModifiedFromIP null.String `json:"last_modified_from_ip"`
	OldId              uint64      `json:"old_id"`
}

type WorkOrderRepository interface {
	Insert(ctx context.Context, u *WorkOrder) error
	UpdateById(ctx context.Context, u *WorkOrder) error
	GetById(ctx context.Context, id uint64) (*WorkOrder, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *WorkOrder) error
}
