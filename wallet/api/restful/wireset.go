package restful

import "github.com/google/wire"

var RestfulWireSet = wire.NewSet(
	NewGinServer,
)
