package api

import "github.com/google/wire"

var ServerWireSet = wire.NewSet(
	NewGinServer,
)
