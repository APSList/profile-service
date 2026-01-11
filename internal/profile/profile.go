package profile

import "go.uber.org/fx"

// ======== EXPORTS ========

// Module exports services present
var Context = fx.Options(
	fx.Provide(GetProfileController),
	fx.Provide(fx.Annotate(
		GetProfileService,
		fx.As(new(Service)),
	)),
	fx.Provide(GetProfileRepository),
	fx.Provide(SetProfileRoutes),
)
