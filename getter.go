package gorm

import (
	"errors"
	"log/slog"
	"sync"

	"gorm.io/gorm"
)

var ErrConfigNotFound = errors.New("gorm config not found")

type DBGetter interface {
	GetDB(name string) (*gorm.DB, error)
}

type Dialector interface {
	Dialector(conn gorm.ConnPool) gorm.Dialector
}

type Channel struct {
	Config    Config
	Dialector Dialector
	Log       *slog.Logger

	once sync.Once
	db   *gorm.DB
}

type Getter struct {
	sync.Mutex
	channels map[string]*Channel
}

func (g *Getter) GetDB(name string) (*gorm.DB, error) {
	g.Lock()
	defer g.Unlock()

	if channel, ok := g.channels[name]; ok {
		return channel.GetDB(nil)
	}

	return nil, ErrConfigNotFound
}

func (c *Channel) GetDB(conn gorm.ConnPool) (*gorm.DB, error) {
	var err error

	c.once.Do(func() {
		c.db, err = gorm.Open(c.Dialector.Dialector(conn), &gorm.Config{
			SkipDefaultTransaction:                   c.Config.SkipDefaultTransaction,
			PrepareStmt:                              c.Config.PrepareStmt,
			DisableNestedTransaction:                 c.Config.DisableNestedTransaction,
			AllowGlobalUpdate:                        c.Config.AllowGlobalUpdate,
			DisableAutomaticPing:                     c.Config.DisableAutomaticPing,
			DisableForeignKeyConstraintWhenMigrating: c.Config.DisableForeignKeyConstraintWhenMigrating,
			IgnoreRelationshipsWhenMigrating:         c.Config.IgnoreRelationshipsWhenMigrating,
			TranslateError:                           c.Config.TranslateError,
			Logger:                                   &gormLogger{log: c.Log},
		})
	})

	return c.db, err
}
