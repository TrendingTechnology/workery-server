package models

import (
	"context"
	"time"
)

// RoleId
//---------------------
// 1 = Executive
// 2 = Management
// 3 = Frontline Staff
// 4 = Associates
// 5 = Customers

// State
//---------------------
// 1 = Active
// 0 = Inactive

type User struct {
	Id                uint64    `json:"id,omitempty"`
	Uuid              string    `json:"uuid,omitempty"`
	TenantId          uint64    `json:"tenant_id,omitempty"`
	Email             string    `json:"email,omitempty"`
	FirstName         string    `json:"first_name,omitempty"`
	LastName          string    `json:"last_name,omitempty"`
	Name              string    `json:"name,omitempty"`
	LexicalName       string    `json:"lexical_name,omitempty"`
	PasswordAlgorithm string    `json:"password_algorithm,omitempty"`
	PasswordHash      string    `json:"password_hash,omitempty"`
	State             int8      `json:"state,omitempty"`
	RoleId            int8      `json:"role_id,omitempty"`
	Timezone          string    `json:"timezone,omitempty"`
	CreatedTime       time.Time `json:"created_time,omitempty"`
	ModifiedTime      time.Time `json:"modified_time,omitempty"`
	JoinedTime        time.Time `json:"joined_time,omitempty"`
	Salt              string    `json:"salt,omitempty"`
	WasEmailActivated bool      `json:"was_email_activated,omitempty"`
	PrAccessCode      string    `json:"pr_access_code,omitempty"`
	PrExpiryTime      time.Time `json:"pr_expiry_time,omitempty"`
	OldId             uint64    `json:"old_id,omitempty"`
	AccessToken       string    `json:"access_token,omitempty"`
	RefreshToken      string    `json:"refresh_token,omitempty"`
}

type UserRepository interface {
	Insert(ctx context.Context, u *User) error
	UpdateById(ctx context.Context, u *User) error
	UpdateByEmail(ctx context.Context, u *User) error
	GetById(ctx context.Context, id uint64) (*User, error)
	GetByOldId(ctx context.Context, oldId uint64) (*User, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	CheckIfExistsByEmail(ctx context.Context, email string) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *User) error
	InsertOrUpdateByEmail(ctx context.Context, u *User) error
}
