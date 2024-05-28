package fiber

import (
	"fmt"
	"github.com/gmlalfjr/go-service-kit/env"
	"github.com/gofiber/fiber/v2"
	"log"
)

// FiberService is a service for managing Fiber app
type FiberService struct {
	App *fiber.App
}

// NewFiberService creates a new FiberService instance
func NewFiberService() *FiberService {
	return &FiberService{
		App: fiber.New(),
	}
}

// Start starts the Fiber service
func (s *FiberService) Start() error {
	port := env.GetString("PORT", "8080")
	port = fmt.Sprintf(":%s", port)
	log.Printf("[Fiber] Starting Fiber service on port %s...", port)
	return s.App.Listen(port)
}

// Stop stops the Fiber service
func (s *FiberService) Stop() error {
	log.Println("[Fiber] Stopping Fiber service...")
	return s.App.Shutdown()
}
