
package workers

import (
	"github.com/rickferrdev/mail-burrow/internal/outbound/queues/workers/email"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"workers",
	email.Invoke,
	email.Provide,
)
