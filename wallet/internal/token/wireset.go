package token

import "github.com/google/wire"

var TokenWireSet = wire.NewSet(
	// NewJWTMaker,
	NewPasetoMaker,
)
