package queues

import (
	"github.com/rickferrdev/mail-burrow/internal/outbound/queues/publisher"
	"github.com/rickferrdev/mail-burrow/internal/outbound/queues/topology"
	"github.com/rickferrdev/mail-burrow/internal/outbound/queues/workers"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"queues",
	topology.Invoke,
	topology.Provide,
	publisher.Provide,
	workers.Module,
)
