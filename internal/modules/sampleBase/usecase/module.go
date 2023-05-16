package usecase

import (
	"go.uber.org/fx"
	repocontract "micro/internal/modules/sampleBase"
)

var Module = fx.Options(
	fx.Provide(func(repo repocontract.IBaseRepository) repocontract.IBaseService {
		return NewBaseService(repo)
	}),
)
