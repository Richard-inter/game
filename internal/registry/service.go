package registry

import (
	"fmt"
	"os"

	"github.com/Richard-inter/game/internal/discovery"
	"github.com/sirupsen/logrus"
)

type ServiceRegistry struct {
	discovery discovery.ServiceDiscovery
	logger    *logrus.Logger
}

func NewServiceRegistry(discovery discovery.ServiceDiscovery, logger *logrus.Logger) *ServiceRegistry {
	return &ServiceRegistry{
		discovery: discovery,
		logger:    logger,
	}
}

func (r *ServiceRegistry) RegisterService(serviceName, host string, port int) error {
	address := fmt.Sprintf("%s:%d", host, port)

	r.logger.WithFields(logrus.Fields{
		"service": serviceName,
		"address": address,
	}).Info("Registering service with etcd")

	err := r.discovery.RegisterService(serviceName, address)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"service": serviceName,
			"address": address,
			"error":   err,
		}).Error("Failed to register service")
		return fmt.Errorf("failed to register service %s: %w", serviceName, err)
	}

	r.logger.WithFields(logrus.Fields{
		"service": serviceName,
		"address": address,
	}).Info("Successfully registered service")

	return nil
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
