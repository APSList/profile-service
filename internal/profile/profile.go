package profile

import "go.uber.org/fx"

// ======== EXPORTS ========

// Module exports services present
var Context = fx.Options(
	fx.Provide(GetProfileController),
	fx.Provide(GetProfileService),
	fx.Provide(GetProfileRepository),
	fx.Provide(SetProfileRoutes),
)
