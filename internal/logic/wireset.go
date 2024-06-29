package logic

import "github.com/google/wire"

var LogicWireSet = wire.NewSet(
	NewAccountLogic,
	NewUserLogic,
	NewtranserLogic,
)
