//go:build !cgo

package gorm

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func (c ConfigSQLite) Dialector(conn gorm.ConnPool) gorm.Dialector {
	return &sqlite.Dialector{
		DriverName: c.DriverName,
		DSN:        c.DSN,
		Conn:       conn,
	}
}
