package gorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Driver string

const (
	MySQL      Driver = "mysql"
	PostgreSQL Driver = "postgresql"
	SQLServer  Driver = "sqlserver"
	SQLite     Driver = "sqlite"
)

type ChannelConfig map[string]Config

type Config struct {
	SkipDefaultTransaction                   bool `mapstructure:"skip_default_transaction" json:"skip_default_transaction,omitempty"`
	PrepareStmt                              bool `mapstructure:"prepare_stmt" json:"prepare_stmt,omitempty"`
	DisableNestedTransaction                 bool `mapstructure:"disable_nested_transaction" json:"disable_nested_transaction,omitempty"`
	AllowGlobalUpdate                        bool `mapstructure:"allow_global_update" json:"allow_global_update,omitempty"`
	DisableAutomaticPing                     bool `mapstructure:"disable_automatic_ping" json:"disable_automatic_ping,omitempty"`
	DisableForeignKeyConstraintWhenMigrating bool `mapstructure:"disable_foreign_key_constraint_when_migrating" json:"disable_foreign_key_constraint_when_migrating,omitempty"`
	IgnoreRelationshipsWhenMigrating         bool `mapstructure:"ignore_relationships_when_migrating" json:"ignore_relationships_when_migrating,omitempty"`
	TranslateError                           bool `mapstructure:"translate_error" json:"translate_error,omitempty"`
}

type ConfigMySQL struct {
	DriverName                    string `mapstructure:"driver_name" json:"driver_name,omitempty"`
	ServerVersion                 string `mapstructure:"server_version" json:"server_version,omitempty"`
	DSN                           string `mapstructure:"dsn" json:"dsn,omitempty"`
	SkipInitializeWithVersion     bool   `mapstructure:"skip_initialize_with_version" json:"skip_initialize_with_version,omitempty"`
	DefaultStringSize             uint   `mapstructure:"default_string_size" json:"default_string_size,omitempty"`
	DefaultDatetimePrecision      *int   `mapstructure:"default_datetime_precision" json:"default_datetime_precision,omitempty"`
	DisableWithReturning          bool   `mapstructure:"disable_with_returning" json:"disable_with_returning,omitempty"`
	DisableDatetimePrecision      bool   `mapstructure:"disable_datetime_precision" json:"disable_datetime_precision,omitempty"`
	DontSupportRenameIndex        bool   `mapstructure:"dont_support_rename_index" json:"dont_support_rename_index,omitempty"`
	DontSupportRenameColumn       bool   `mapstructure:"dont_support_rename_column" json:"dont_support_rename_column,omitempty"`
	DontSupportRenameColumnUnique bool   `mapstructure:"dont_support_rename_column_unique" json:"dont_support_rename_column_unique,omitempty"`
	DontSupportForShareClause     bool   `mapstructure:"dont_support_for_share_clause" json:"dont_support_for_share_clause,omitempty"`
	DontSupportNullAsDefaultValue bool   `mapstructure:"dont_support_null_as_default_value" json:"dont_support_null_as_default_value,omitempty"`
}

func (c ConfigMySQL) Dialector(conn gorm.ConnPool) gorm.Dialector {
	defaultStringSize := c.DefaultStringSize
	if defaultStringSize == 0 {
		defaultStringSize = 256
	}

	return mysql.New(mysql.Config{
		DriverName:                    c.DriverName,
		ServerVersion:                 c.ServerVersion,
		DSN:                           c.DSN,
		Conn:                          conn,
		SkipInitializeWithVersion:     c.SkipInitializeWithVersion,
		DefaultStringSize:             defaultStringSize,
		DefaultDatetimePrecision:      c.DefaultDatetimePrecision,
		DisableWithReturning:          c.DisableWithReturning,
		DisableDatetimePrecision:      c.DisableDatetimePrecision,
		DontSupportRenameIndex:        c.DontSupportRenameIndex,
		DontSupportRenameColumn:       c.DontSupportRenameColumn,
		DontSupportRenameColumnUnique: c.DontSupportRenameColumnUnique,
		DontSupportForShareClause:     c.DontSupportForShareClause,
		DontSupportNullAsDefaultValue: c.DontSupportNullAsDefaultValue,
	})
}

type ConfigPostgreSQL struct {
	DriverName           string `mapstructure:"driver_name" json:"driver_name,omitempty"`
	DSN                  string `mapstructure:"dsn" json:"dsn,omitempty"`
	PreferSimpleProtocol bool   `mapstructure:"prefer_simple_protocol" json:"prefer_simple_protocol,omitempty"`
	WithoutReturning     bool   `mapstructure:"without_returning" json:"without_returning,omitempty"`
}

func (c ConfigPostgreSQL) Dialector(conn gorm.ConnPool) gorm.Dialector {
	return postgres.New(postgres.Config{
		DriverName:           c.DriverName,
		DSN:                  c.DSN,
		PreferSimpleProtocol: c.PreferSimpleProtocol,
		WithoutReturning:     c.WithoutReturning,
		Conn:                 conn,
	})
}

type ConfigSQLServer struct {
	DriverName        string `mapstructure:"driver_name" json:"driver_name,omitempty"`
	DSN               string `mapstructure:"dsn" json:"dsn,omitempty"`
	DefaultStringSize int    `mapstructure:"default_string_size" json:"default_string_size,omitempty"`
}

func (c ConfigSQLServer) Dialector(conn gorm.ConnPool) gorm.Dialector {
	return sqlserver.New(sqlserver.Config{
		DriverName:        c.DriverName,
		DSN:               c.DSN,
		DefaultStringSize: c.DefaultStringSize,
		Conn:              conn,
	})
}

type ConfigSQLite struct {
	DriverName string `mapstructure:"driver_name" json:"driver_name,omitempty"`
	DSN        string `mapstructure:"dsn" json:"dsn,omitempty"`
}
