package registry

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/1nterdigital/game/internal/discovery"
)

type ServiceRegistry struct {
	discovery discovery.ServiceDiscovery
	logger    *zap.SugaredLogger
}

func NewServiceRegistry(discovery discovery.ServiceDiscovery, logger *zap.SugaredLogger) *ServiceRegistry {
	return &ServiceRegistry{
		discovery: discovery,
		logger:    logger,
	}
}

func (r *ServiceRegistry) RegisterService(serviceName, host string, port int) error {
	address := fmt.Sprintf("%s:%d", host, port)

	r.logger.Infow("Registering service with etcd", "service", serviceName, "address", address)

	err := r.discovery.RegisterService(serviceName, address)
	if err != nil {
		r.logger.Errorw("Failed to register service", "service", serviceName, "address", address, "error", err)
		return fmt.Errorf("failed to register service %s: %w", serviceName, err)
	}

	r.logger.Infow("Successfully registered service", "service", serviceName, "address", address)

	return nil
}

func (r *ServiceRegistry) GetService(serviceName string) (string, error) {
	address, err := r.discovery.GetService(serviceName)
	if err != nil {
		r.logger.Errorw("Failed to get service from etcd", "service", serviceName, "error", err)
		return "", fmt.Errorf("failed to get service %s: %w", serviceName, err)
	}

	r.logger.Debugw("Retrieved service from etcd", "service", serviceName, "address", address)
	return address, nil
}

// Helper function to get service address from environment
func GetServiceAddress(serviceName string) (string, int, error) {
	host := os.Getenv(fmt.Sprintf("%s_HOST", serviceName))
	portStr := os.Getenv(fmt.Sprintf("%s_PORT", serviceName))

	if host == "" {
		host = "localhost"
	}

	if portStr == "" {
		return "", 0, fmt.Errorf("port not found for service %s", serviceName)
	}

	return host, 0, nil
}
