package config

import "github.com/google/wire"

var ConfigWireSet = wire.NewSet(
	LoadAllConfig,
)
