package models

import (
	"context"
	"time"
)

type User struct {
	Id                uint64    `json:"id"`
	Uuid              string    `json:"uuid"`
	TenantId          uint64    `json:"tenant_id"`
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	PasswordAlgorithm string    `json:"password_algorithm"`
	PasswordHash      string    `json:"password_hash"`
	State             int8      `json:"state"`
	Role              int8      `json:"role"`
	Timezone          string    `json:"timezone"`
	CreatedTime       time.Time `json:"created_time"`
	ModifiedTime      time.Time `json:"modified_time"`
	JoinedTime        time.Time `json:"joined_time"`
	Salt              string    `json:"salt"`
	WasEmailActivated bool      `json:"was_email_activated"`
	PrAccessCode      string    `json:"pr_access_code"`
	PrExpiryTime      time.Time `json:"pr_expiry_time"`
}

type UserRepository interface {
	Insert(ctx context.Context, u *User) error
	UpdateById(ctx context.Context, u *User) error
	UpdateByEmail(ctx context.Context, u *User) error
	GetById(ctx context.Context, id uint64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	CheckIfExistsByEmail(ctx context.Context, email string) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *User) error
	InsertOrUpdateByEmail(ctx context.Context, u *User) error
}
