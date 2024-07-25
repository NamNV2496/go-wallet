package api

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-wallet/api/restful"
)

var ServerWireSet = wire.NewSet(
	restful.RestfulWireSet,
)
