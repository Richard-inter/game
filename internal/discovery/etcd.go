package discovery

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdDiscovery struct {
	client   *clientv3.Client
	services map[string]string
}

type ServiceDiscovery interface {
	GetService(serviceName string) (string, error)
	RegisterService(serviceName, address string) error
	Close() error
}

func NewEtcdDiscovery(etcdEndpoints []string) (*EtcdDiscovery, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}

	discovery := &EtcdDiscovery{
		client:   client,
		services: make(map[string]string),
	}

	// Load existing services
	err = discovery.loadServices()
	if err != nil {
		return nil, fmt.Errorf("failed to load services: %w", err)
	}

	return discovery, nil
}

func (d *EtcdDiscovery) GetService(serviceName string) (string, error) {
	if address, exists := d.services[serviceName]; exists {
		return address, nil
	}

	// Fallback: fetch from etcd
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := d.client.Get(ctx, serviceName)
	if err != nil {
		return "", fmt.Errorf("failed to get service %s from etcd: %w", serviceName, err)
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("service %s not found in etcd", serviceName)
	}

	address := string(resp.Kvs[0].Value)
	d.services[serviceName] = address // Cache locally
	return address, nil
}

func (d *EtcdDiscovery) RegisterService(serviceName, address string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Register with TTL (30 seconds)
	lease, err := d.client.Grant(ctx, 30)
	if err != nil {
		return fmt.Errorf("failed to create lease: %w", err)
	}

	// Keep lease alive
	_, err = d.client.KeepAlive(ctx, lease.ID)
	if err != nil {
		return fmt.Errorf("failed to keep lease alive: %w", err)
	}

	// Register service
	_, err = d.client.Put(ctx, serviceName, address, clientv3.WithLease(lease.ID))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	d.services[serviceName] = address
	return nil
}

func (d *EtcdDiscovery) Close() error {
	return d.client.Close()
}

func (d *EtcdDiscovery) loadServices() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := d.client.Get(ctx, "", clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("failed to load services: %w", err)
	}

	for _, kv := range resp.Kvs {
		d.services[string(kv.Key)] = string(kv.Value)
	}

	return nil
}
