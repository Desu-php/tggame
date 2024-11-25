package adapter

import "go.uber.org/fx"

var Module = fx.Provide(
	NewUserSessionCacheAdapter,
)