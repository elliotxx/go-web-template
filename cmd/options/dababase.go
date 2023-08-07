package options

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/elliotxx/go-web-template/cmd/options/types"
	"github.com/elliotxx/go-web-template/pkg/server"
	gomysql "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrDBHostNotSpecified               = errors.New("--db-host must be specified")
	ErrDBNameNotSpecified               = errors.New("--db-name must be specified")
	ErrDBUserNotSpecified               = errors.New("--db-user must be specified")
	ErrDBPortNotSpecified               = errors.New("--db-port must be specified")
	_                     types.Options = &DatabaseOptions{}
)

// DatabaseOptions is a Database options struct
type DatabaseOptions struct {
	DBName     string `json:"dbName,omitempty" yaml:"dbName,omitempty"`
	DBUser     string `json:"dbUser,omitempty" yaml:"dbUser,omitempty"`
	DBPassword string `json:"dbPassword,omitempty" yaml:"dbPassword,omitempty"`
	DBHost     string `json:"dbHost,omitempty" yaml:"dbHost,omitempty"`
	DBPort     int    `json:"dbPort,omitempty" yaml:"dbPort,omitempty"`
	// AutoMigrate will attempt to automatically migrate all tables
	AutoMigrate bool   `json:"autoMigrate,omitempty" yaml:"autoMigrate,omitempty"`
	MigrateFile string `json:"migrateFile,omitempty" yaml:"migrateFile,omitempty"`
}

// NewDatabaseOptions returns a DatabaseOptions instance with the default values
func NewDatabaseOptions() *DatabaseOptions {
	return &DatabaseOptions{
		DBHost:      "127.0.0.1",
		DBPort:      3306,
		AutoMigrate: false,
	}
}

// InstallDB uses the run options to generate and open a db session.
func (o *DatabaseOptions) InstallDB() (*gorm.DB, error) {
	// Generate go-sql-driver.mysql config to format DSN
	config := gomysql.NewConfig()
	config.User = o.DBUser
	config.Passwd = o.DBPassword
	config.Addr = o.DBHost + ":" + strconv.Itoa(o.DBPort)
	config.DBName = o.DBName
	config.Net = "tcp"
	config.ParseTime = true
	config.InterpolateParams = true
	config.Params = map[string]string{
		"charset": "utf8",
		"loc":     "Asia/Shanghai",
	}
	dsn := config.FormatDSN()
	// silence log output
	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	return gorm.Open(mysql.Open(dsn), cfg) // todo: add db connection check to healthz check
}

// Validate checks DatabaseOptions and return a slice of found error(s)
func (o *DatabaseOptions) Validate() error {
	if o == nil {
		return errors.Errorf("options is nil")
	}

	if o.AutoMigrate && len(o.MigrateFile) == 0 {
		return errors.Errorf("when --auto-migrate is true, --migrate-file must be specified")
	}

	if len(o.DBHost) == 0 {
		return ErrDBHostNotSpecified
	}
	if len(o.DBName) == 0 {
		return ErrDBNameNotSpecified
	}
	if len(o.DBUser) == 0 {
		return ErrDBUserNotSpecified
	}
	if o.DBPort == 0 {
		return ErrDBPortNotSpecified
	}

	return nil
}

// ApplyTo apply database options to the server config
func (o *DatabaseOptions) ApplyTo(config *server.Config) {
	d, err := o.InstallDB()
	if err != nil {
		logrus.Fatalf("Failed to apply database options to server.Config as: %+v", err)
	}
	config.DB = d

	// AutoMigrate will attempt to automatically migrate all tables
	if o.AutoMigrate {
		logrus.Debugf("AutoMigrate will attempt to automatically migrate all tables from [%s]", o.MigrateFile)
		// Read all content by migrate file
		migrateSQL, err := os.ReadFile(o.MigrateFile)
		if err != nil {
			logrus.Fatalf("Failed to read migrate file: %+v", err)
		}

		// Split multiple SQL statements into individual statements
		stmts := strings.Split(string(migrateSQL), ";")

		// Iterate over all statements and execute them
		for _, stmt := range stmts {
			// Ignore empty statements
			if len(strings.TrimSpace(stmt)) == 0 {
				continue
			}

			// Use gorm.Exec() function to execute SQL statement
			if err = config.DB.Exec(stmt).Error; err != nil {
				logrus.Warnf("Failed to exec migrate sql: %+v", err)
			}
		}
	}
}

// AddFlags adds flags for a specific Option to the specified FlagSet
func (o *DatabaseOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.DBName, "db-name", o.DBName, "the database name")
	fs.StringVar(&o.DBUser, "db-user", o.DBUser, "the user name used to access database")
	fs.StringVar(&o.DBPassword, "db-pwd", o.DBPassword, "the user password used to access database")
	fs.StringVar(&o.DBHost, "db-host", o.DBHost, "database host")
	fs.IntVar(&o.DBPort, "db-port", o.DBPort, "database port")
	fs.BoolVar(&o.AutoMigrate, "auto-migrate", o.AutoMigrate, "Whether to enable automatic migration")
	fs.StringVar(&o.MigrateFile, "migrate-file", o.MigrateFile, "The migrate sql file")
}

// MarshalJSON is custom marshalling function for masking sensitive field values
func (o DatabaseOptions) MarshalJSON() ([]byte, error) {
	type tempOptions DatabaseOptions
	o2 := tempOptions(o)
	o2.DBPassword = types.MaskString
	return json.Marshal(&o2)
}
