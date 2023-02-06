package upgraders

import (
	"context"

	"github.com/go-logr/logr"
)

// UpgradeDelayedCheck will raise a 'delayed' event if the cluster has not commenced
// upgrade within a configurable amount of time.
func (h *hcpUpgrader) CommenceHcpUpgrade(ctx context.Context, logger logr.Logger) (bool, error) {
	return true, nil
}

func (h *hcpUpgrader) HcpUpgraded(ctx context.Context, logger logr.Logger) (bool, error) {
	return true, nil
}

func (h *hcpUpgrader) IsHcpUpgradeable(ctx context.Context, logger logr.Logger) (bool, error) {
	return true, nil
}
