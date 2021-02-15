package models

import (
    "context"
    "time"
)

type Tenant struct {
    Id uint64
    Uuid string
    Name string
    State int8
    Timezone string
    CreatedTime time.Time
    ModifiedTime time.Time
}

type TenantRepository interface {
    Insert(ctx context.Context, u *Tenant) error
    Update(ctx context.Context, u *Tenant) error
    GetById(ctx context.Context, id uint64) (*Tenant, error)
    GetByName(ctx context.Context, name string) (*Tenant, error)
    CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
    CheckIfExistsByName(ctx context.Context, name string) (bool, error)
    InsertOrUpdate(ctx context.Context, u *Tenant) error
}
