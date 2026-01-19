package grpc

import (
	"fmt"
)

type ClientManager struct {
	discovery          *DiscoveryManager
	player             *PlayerClient
	clawmachine        *ClawMachineClient
	clawmachineRuntime *ClawMachineRuntimeClient
	gachaMachine       *GachaMachineClient
	// Direct connection addresses for when discovery is disabled
	playerAddr       string
	clawmachineAddr  string
	runtimeAddr      string
	gachaMachineAddr string
}

type ClientManagerConfig struct {
	EtcdEndpoints    []string
	PlayerAddr       string // Direct address when discovery disabled
	ClawmachineAddr  string // Direct address when discovery disabled
	RuntimeAddr      string // Direct address when discovery disabled
	GachaMachineAddr string // Direct address when discovery disabled
}

func NewClientManager(cfg *ClientManagerConfig) (*ClientManager, error) {
	var discovery *DiscoveryManager
	var err error

	// Only create discovery manager if etcd endpoints are provided
	if len(cfg.EtcdEndpoints) > 0 {
		discovery, err = NewDiscoveryManager(cfg.EtcdEndpoints)
		if err != nil {
			return nil, fmt.Errorf("failed to create service discovery: %w", err)
		}
	}

	return &ClientManager{
		discovery:        discovery,
		playerAddr:       cfg.PlayerAddr,
		clawmachineAddr:  cfg.ClawmachineAddr,
		runtimeAddr:      cfg.RuntimeAddr,
		gachaMachineAddr: cfg.GachaMachineAddr,
	}, nil
}

func (cm *ClientManager) GetPlayerClient() (*PlayerClient, error) {
	if cm.player == nil {
		var playerAddr string
		var err error

		// Use discovery if available, otherwise use direct address
		if cm.discovery != nil {
			playerAddr, err = cm.discovery.GetService("player-service")
			if err != nil {
				return nil, fmt.Errorf("failed to get player service address: %w", err)
			}
		} else {
			if cm.playerAddr == "" {
				return nil, fmt.Errorf("player service address not configured and discovery is disabled")
			}
			playerAddr = cm.playerAddr
		}

		// Create client
		player, err := NewPlayerClient(playerAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to create player client: %w", err)
		}

		cm.player = player
	}

	return cm.player, nil
}

func (cm *ClientManager) GetClawMachineClient() (*ClawMachineClient, error) {
	if cm.clawmachine == nil {
		var clawmachineAddr string
		var err error

		// Use discovery if available, otherwise use direct address
		if cm.discovery != nil {
			clawmachineAddr, err = cm.discovery.GetService("clawmachine-service")
			if err != nil {
				return nil, fmt.Errorf("failed to get clawmachine service address: %w", err)
			}
		} else {
			if cm.clawmachineAddr == "" {
				return nil, fmt.Errorf("clawmachine service address not configured and discovery is disabled")
			}
			clawmachineAddr = cm.clawmachineAddr
		}

		// Create client
		clawmachine, err := NewClawMachineClient(clawmachineAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to create clawmachine client: %w", err)
		}

		cm.clawmachine = clawmachine
	}

	return cm.clawmachine, nil
}

func (cm *ClientManager) GetClawMachineRuntimeClient() (*ClawMachineRuntimeClient, error) {
	if cm == nil {
		return nil, fmt.Errorf("client manager is nil")
	}

	if cm.clawmachineRuntime == nil {
		var runtimeAddr string
		var err error

		// Use discovery if available, otherwise use direct address
		if cm.discovery != nil {
			runtimeAddr, err = cm.discovery.GetService("clawmachine-runtime-service")
			if err != nil {
				return nil, fmt.Errorf("failed to get clawmachine runtime service address: %w", err)
			}
		} else {
			if cm.runtimeAddr == "" {
				return nil, fmt.Errorf("clawmachine runtime service address not configured and discovery is disabled")
			}
			runtimeAddr = cm.runtimeAddr
		}

		// Create client
		clawmachineRuntime, err := NewClawMachineRuntimeClient(runtimeAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to create clawmachine runtime client: %w", err)
		}
		cm.clawmachineRuntime = clawmachineRuntime
	}

	return cm.clawmachineRuntime, nil
}

func (cm *ClientManager) GetGachaMachineClient() (*GachaMachineClient, error) {
	if cm == nil {
		return nil, fmt.Errorf("client manager is nil")
	}

	if cm.gachaMachine == nil {
		var gachaMachineAddr string
		var err error

		// Use discovery if available, otherwise use direct address
		if cm.discovery != nil {
			gachaMachineAddr, err = cm.discovery.GetService("gachamachine-service")
			if err != nil {
				return nil, fmt.Errorf("failed to get gachamachine service address: %w", err)
			}
		} else {
			if cm.gachaMachineAddr == "" {
				return nil, fmt.Errorf("gachamachine service address not configured and discovery is disabled")
			}
			gachaMachineAddr = cm.gachaMachineAddr
		}

		// Create client
		gachaMachine, err := NewGachaMachineClient(gachaMachineAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to create gachamachine client: %w", err)
		}
		cm.gachaMachine = gachaMachine
	}

	return cm.gachaMachine, nil
}

func (cm *ClientManager) Close() error {
	var errors []string

	if cm.player != nil {
		if err := cm.player.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("player client: %v", err))
		}
	}

	if cm.clawmachine != nil {
		if err := cm.clawmachine.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("clawmachine client: %v", err))
		}
	}
	if cm.clawmachineRuntime != nil {
		// Note: Runtime client doesn't have Close method in the interface
		// Connection will be closed when the process exits
	}

	if cm.discovery != nil {
		if err := cm.discovery.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("discovery: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple errors during close: %v", errors)
	}

	return nil
}
