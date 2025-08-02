package contract

import (
	"github.com/hbttundar/scg-database/config"
)

type (
	DBAdapter interface {
		// Connect's only job is to create our rich Connection object from a config struct.
		// Options should be applied *before* this is called.
		Connect(cfg *config.Config) (Connection, error)
		Name() string
	}
)
