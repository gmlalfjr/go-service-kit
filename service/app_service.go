package service

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Service adalah interface yang harus dipenuhi oleh semua layanan
type Service interface {
	Start() error
	Stop() error
}

// AppService is a service for managing multiple services
type AppService struct {
	services []Service // Slice untuk menyimpan berbagai layanan
}

// NewAppService creates a new AppService instance with variadic services
func NewAppService(services ...Service) *AppService {
	return &AppService{
		services: services,
	}
}

// Start starts all registered services in AppService
func (s *AppService) Start() error {
	errChan := make(chan error, len(s.services))
	for _, service := range s.services {
		go func(srv Service) {
			if err := srv.Start(); err != nil {
				errChan <- err
			}
		}(service)
	}

	// Handle OS signals for graceful shutdown
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt)
	signal.Notify(quitSignal, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case <-quitSignal:
		log.Println("Received shutdown signal, stopping services...")
		return s.Stop()
	}
}

// Stop stops all registered services in AppService
func (s *AppService) Stop() error {
	for _, service := range s.services {
		if err := service.Stop(); err != nil {
			return err
		}
	}
	return nil
}
