package gorm

import (
	"fmt"

	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
)

const PluginName = "gorm"

var drivers = []Driver{MySQL, PostgreSQL, SQLServer, SQLite}

type Plugin struct {
	getter *Getter
}

func (p *Plugin) Init(cfg Configurer, logger Logger) error {
	const op = errors.Op("gorm_plugin_init")

	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	var channelCfg ChannelConfig
	if err := cfg.UnmarshalKey(PluginName, &channelCfg); err != nil {
		return errors.E(op, err)
	}

	log := logger.NamedLogger(PluginName)
	channels := map[string]*Channel{}

	for name, config := range channelCfg {
		for _, driver := range drivers {
			if dialector, ok := getDialector(cfg, name, driver); ok {
				channels[name] = &Channel{
					Config:    config,
					Dialector: dialector,
					Log:       log.WithGroup(name).WithGroup(string(driver)),
				}
				break
			}
		}
	}

	if len(channels) == 0 {
		return errors.E(op, errors.Disabled)
	}

	p.getter = &Getter{channels: channels}

	return nil
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*DBGetter)(nil), p.Getter),
	}
}

func (p *Plugin) Getter() *Getter {
	return p.getter
}

func getDialector(cfg Configurer, name string, driver Driver) (Dialector, bool) {
	key := fmt.Sprintf("%s.%s.%s", PluginName, name, driver)

	if cfg.Has(key) {
		switch driver {
		case MySQL:
			var dialector ConfigMySQL
			err := cfg.UnmarshalKey(key, &dialector)
			return dialector, err == nil
		case PostgreSQL:
			var dialector ConfigPostgreSQL
			err := cfg.UnmarshalKey(key, &dialector)
			return dialector, err == nil
		case SQLServer:
			var dialector ConfigSQLServer
			err := cfg.UnmarshalKey(key, &dialector)
			return dialector, err == nil
		case SQLite:
			var dialector ConfigSQLite
			err := cfg.UnmarshalKey(key, &dialector)
			return dialector, err == nil
		}
	}

	return nil, false
}
