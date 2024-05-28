package fiber

import (
	"fmt"
	"github.com/gmlalfjr/go-service-kit/env"
	"github.com/gofiber/fiber/v2"
	"log"
)

// FiberService is a service for managing Fiber app
type FiberService struct {
	app *fiber.App
}

// NewFiberService creates a new FiberService instance
func NewFiberService() *FiberService {
	return &FiberService{
		app: fiber.New(),
	}
}

// Start starts the Fiber service
func (s *FiberService) Start() error {
	port := env.GetString("PORT", "8080")
	port = fmt.Sprintf(":%s", port)
	log.Printf("[Fiber] Starting Fiber service on port %s...", port)
	return s.app.Listen(port)
}

// Stop stops the Fiber service
func (s *FiberService) Stop() error {
	log.Println("[Fiber] Stopping Fiber service...")
	return s.app.Shutdown()
}
