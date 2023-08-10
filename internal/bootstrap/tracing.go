package bootstrap

import (
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/tracing"
	"github.com/pauloo27/logger"
)

func initTracing() {
	if !config.Config.Tracing.Enabled {
		logger.Info("Tracing is disabled")
		tracing.DisableTracer()
		return
	}

	err := tracing.InitTracer(
		config.Config.Tracing.Endpoint,
		config.Config.Tracing.ServiceName,
		config.Config.Env,
	)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Success("Tracing enabled")
}