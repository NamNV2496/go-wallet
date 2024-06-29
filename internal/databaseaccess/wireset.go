package databaseaccess

import (
	"github.com/google/wire"
)

var DatabaseWireSet = wire.NewSet(
	NewDatabase,
)
