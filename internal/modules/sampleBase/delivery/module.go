package delivery

import (
	"micro/internal/modules/sampleBase/delivery/grpc"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(controller.NewBaseController),
)
