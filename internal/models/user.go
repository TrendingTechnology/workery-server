package models

import (
	"context"
	"time"
)

type User struct {
	Id           uint64
	Uuid         string
	TenantId     uint64
	Email        string
	FirstName    string
	LastName     string
	PasswordHash string
	State        int8
	Role         int8
	Timezone     string
	CreatedTime  time.Time
	ModifiedTime time.Time
}

type UserRepository interface {
	Insert(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	GetById(ctx context.Context, id uint64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	CheckIfExistsByEmail(ctx context.Context, email string) (bool, error)
	InsertOrUpdate(ctx context.Context, u *User) error
}
