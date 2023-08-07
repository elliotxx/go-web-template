package persistence

import (
	"context"

	"github.com/elliotxx/go-web-template/pkg/domain/entity"
	"github.com/elliotxx/go-web-template/pkg/domain/repository"
	"github.com/elliotxx/errors"
	"gorm.io/gorm"
)

// The systemConfigRepository type implements the repository.SystemConfigRepository interface.
// If the systemConfigRepository type does not implement all the methods of the interface,
// the compiler will produce an error.
var _ repository.SystemConfigRepository = &systemConfigRepository{}

// systemConfigRepository is a repository that stores systemConfigs in a gorm database.
type systemConfigRepository struct {
	// db is the underlying gorm database where systemConfigs are stored.
	db *gorm.DB
}

// NewSystemConfigRepository creates a new systemConfig repository.
func NewSystemConfigRepository(db *gorm.DB) repository.SystemConfigRepository {
	return &systemConfigRepository{db: db}
}

// Create saves a system config to the repository.
func (r *systemConfigRepository) Create(ctx context.Context, dataEntity *entity.SystemConfig) error {
	err := dataEntity.Validate()
	if err != nil {
		return err
	}

	// Map the data from Entity to DO
	var dataModel SystemConfigModel
	err = dataModel.FromEntity(dataEntity)
	if err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create new record in the store
		err = tx.WithContext(ctx).Create(&dataModel).Error
		if err != nil {
			return err
		}

		// Map fresh record's data into Entity
		newEntity, err := dataModel.ToEntity()
		if err != nil {
			return err
		}
		*dataEntity = *newEntity

		return nil
	})
}

// Delete removes a system config from the repository.
func (r *systemConfigRepository) Delete(ctx context.Context, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var dataModel SystemConfigModel
		err := tx.WithContext(ctx).First(&dataModel, id).Error
		if err != nil {
			return err
		}

		return tx.WithContext(ctx).Delete(&dataModel).Error
	})
}

// Update updates an existing system config in the repository.
func (r *systemConfigRepository) Update(ctx context.Context, dataEntity *entity.SystemConfig) error {
	// Map the data from Entity to DO
	var dataModel SystemConfigModel
	err := dataModel.FromEntity(dataEntity)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Updates(&dataModel).Error
}

// Find retrieves a system config by its ID.
func (r *systemConfigRepository) Get(ctx context.Context, id uint) (*entity.SystemConfig, error) {
	var dataModel SystemConfigModel
	err := r.db.WithContext(ctx).First(&dataModel, id).Error
	if err != nil {
		return nil, err
	}

	return dataModel.ToEntity()
}

// Find returns a list of specified system configs in the repository.
func (r *systemConfigRepository) Find(ctx context.Context, query repository.Query) ([]*entity.SystemConfig, error) {
	var systemConfigModels []*SystemConfigModel
	if err := r.db.WithContext(ctx).
		Where("config LIKE ?", "%"+query.Keyword+"%").
		Limit(query.Limit).
		Offset(query.Offset).
		Find(&systemConfigModels).Error; err != nil {
		return nil, err
	}

	systemConfigEntities := make([]*entity.SystemConfig, 0, len(systemConfigModels))
	for _, model := range systemConfigModels {
		newEntity, err := model.ToEntity()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert db model (ID: %d) to entity", model.ID)
		}

		systemConfigEntities = append(systemConfigEntities, newEntity)
	}
	return systemConfigEntities, nil
}

// Count returns the total of system configs.
func (r *systemConfigRepository) Count(ctx context.Context) (int, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&SystemConfigModel{}).Count(&total).Error
	if err != nil {
		return 0, err
	}

	return int(total), nil
}
