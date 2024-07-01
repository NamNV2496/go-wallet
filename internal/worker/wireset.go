package worker

import "github.com/google/wire"

var WorkerWireSet = wire.NewSet(
	NewTaskProcessor,
)
