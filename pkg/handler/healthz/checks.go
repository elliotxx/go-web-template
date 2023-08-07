package healthz

import (
	"github.com/elliotxx/healthcheck/checks"
	"gorm.io/gorm"
)

// gormDBCheck is a check that returns true if the database is
// available.
type gormDBCheck struct {
	db *gorm.DB
}

func NewGormDBCheck(db *gorm.DB) checks.Check {
	return &gormDBCheck{
		db: db,
	}
}

func (c *gormDBCheck) Name() string {
	return "Database"
}

func (c *gormDBCheck) Pass() bool {
	if c.db == nil {
		return false
	}

	sqldb, err := c.db.DB()
	if err != nil {
		return false
	}

	sqlCheck := checks.NewSQLCheck(sqldb)
	return sqlCheck.Pass()
}
