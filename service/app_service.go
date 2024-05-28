package service

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Service interface {
	Start() error
	Stop() error
}

type AppService struct {
	services []Service
}

func NewAppService(services ...Service) *AppService {
	return &AppService{
		services: services,
	}
}

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
