package app

import (
	"github.com/gmlalfjr/go-service-kit/logger"
	"github.com/gmlalfjr/go-service-kit/service"
)

// Application struct to manage app lifecycle
type Application struct {
	logger     logger.LoggerConfig
	appService *service.AppService
}

// NewApplication creates a new Application instance
func NewApplication(logger logger.Logger, services ...service.Service) *Application {
	appService := service.NewAppService(services...)
	return &Application{
		logger:     logger,
		appService: appService,
	}
}

// Start starts the application
func (a *Application) Start() error {
	a.logger.Print("Starting application...")

	// Start all registered services
	if err := a.appService.Start(); err != nil {
		a.logger.Error("Failed to start application:", err)
		return err
	}

	return nil
}

// Stop stops the application
func (a *Application) Stop() error {
	a.logger.Info("Stopping application...")

	// Stop all registered services
	if err := a.appService.Stop(); err != nil {
		a.logger.Error("Failed to stop application:", err)
		return err
	}

	return nil
}
