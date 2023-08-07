package repository

import (
	"context"

	"github.com/elliotxx/go-web-template/pkg/domain/entity"
)

// SystemConfigRepository is an interface that defines the repository
// operations for system config.
// It follows the principles of domain-driven design (DDD).
type SystemConfigRepository interface {
	// Create creates a new system config.
	Create(ctx context.Context, systemConfig *entity.SystemConfig) error
	// Delete deletes a system config by its ID.
	Delete(ctx context.Context, id uint) error
	// Update updates an existing system config.
	Update(ctx context.Context, systemConfig *entity.SystemConfig) error
	// Get retrieves a system config by its ID.
	Get(ctx context.Context, id uint) (*entity.SystemConfig, error)
	// Find returns a list of specified system config.
	Find(ctx context.Context, query Query) ([]*entity.SystemConfig, error)
	// Count returns the total of system configs.
	Count(ctx context.Context) (int, error)
}
