package persistence

import (
	"github.com/elliotxx/errors"
	"github.com/elliotxx/go-web-template/pkg/domain/entity"
	"gorm.io/gorm"
)

// SystemConfigModel is a DO used to map the entity to the database.
type SystemConfigModel struct {
	gorm.Model
	Tenant      string
	Env         string
	Type        string
	Config      string
	Description string
	Creator     string
	Modifier    string
}

// The TableName method returns the name of the database table that the struct is mapped to.
func (m *SystemConfigModel) TableName() string {
	return "system_config"
}

// ToEntity converts the DO to an entity.
func (m *SystemConfigModel) ToEntity() (*entity.SystemConfig, error) {
	if m == nil {
		return nil, ErrSystemConfigModelNil
	}

	env, err := entity.ParseEnv(m.Env)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}

	return &entity.SystemConfig{
		ID:          m.ID,
		Tenant:      m.Tenant,
		Env:         env,
		Type:        m.Type,
		Config:      m.Config,
		Description: m.Description,
		Creator:     m.Creator,
		Modifier:    m.Modifier,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}

// FromEntity converts an entity to a DO.
func (m *SystemConfigModel) FromEntity(e *entity.SystemConfig) error {
	if m == nil {
		return ErrSystemConfigModelNil
	}

	m.ID = e.ID
	m.Tenant = e.Tenant
	m.Env = string(e.Env)
	m.Type = e.Type
	m.Config = e.Config
	m.Description = e.Description
	m.Creator = e.Creator
	m.Modifier = e.Modifier
	m.CreatedAt = e.CreatedAt
	m.UpdatedAt = e.UpdatedAt

	return nil
}
