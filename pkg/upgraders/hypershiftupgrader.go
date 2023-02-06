package upgraders

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	upgradev1alpha1 "github.com/openshift/managed-upgrade-operator/api/v1alpha1"
	cv "github.com/openshift/managed-upgrade-operator/pkg/clusterversion"
	"github.com/openshift/managed-upgrade-operator/pkg/configmanager"
	"github.com/openshift/managed-upgrade-operator/pkg/eventmanager"
	"github.com/openshift/managed-upgrade-operator/pkg/machinery"
	"github.com/openshift/managed-upgrade-operator/pkg/metrics"
	"github.com/openshift/managed-upgrade-operator/pkg/notifier"
	"github.com/openshift/managed-upgrade-operator/pkg/scaler"
	"github.com/openshift/managed-upgrade-operator/pkg/upgradesteps"
)

// osdUpgrader is a cluster upgrader suitable for OpenShift Dedicated clusters.
// It inherits from the base clusterUpgrader.
type hcpUpgrader struct {
	*clusterUpgrader
}

// NewOSDUpgrader creates a new instance of an osdUpgrader
func NewHCPUpgrader(c client.Client, cfm configmanager.ConfigManager, mc metrics.Metrics, notifier eventmanager.EventManager) (*hcpUpgrader, error) {
	cfg := &upgraderConfig{}
	err := cfm.Into(cfg)
	if err != nil {
		return nil, err
	}

	hu := hcpUpgrader{
		clusterUpgrader: &clusterUpgrader{
			client:               c,
			metrics:              mc,
			cvClient:             cv.NewCVClient(c),
			config:               cfg,
			machinery:            machinery.NewMachinery(),
		},
	}

	steps := []upgradesteps.UpgradeStep{
		upgradesteps.Action(string(upgradev1alpha1.IsHcpClusterUpgradable), hu.IsHcpUpgradeable),
		upgradesteps.Action(string(upgradev1alpha1.CommenceUpgrade), hu.CommenceHcpUpgrade),
		upgradesteps.Action(string(upgradev1alpha1.ControlPlaneUpgraded), hu.HcpUpgraded),
	}
	hu.steps = steps

	return &hu, nil
}

// UpgradeCluster performs the upgrade of the cluster and returns an indication of the
// last-executed upgrade phase and any error associated with the phase execution.
//
// The UpgradeCluster enforces OSD policy around expiring upgrades if they do not commence
// within a given time period.
func (u *hcpUpgrader) UpgradeCluster(ctx context.Context, upgradeConfig *upgradev1alpha1.UpgradeConfig, logger logr.Logger) (upgradev1alpha1.UpgradePhase, error) {
	u.upgradeConfig = upgradeConfig

	// OSD upgrader enforces a 'failure' policy if the upgrade does not commence within a time period
	if cancelUpgrade, _ := shouldFailUpgrade(u.cvClient, u.config, u.upgradeConfig); cancelUpgrade {
		return hcpPerformUpgradeFailure(u.client, u.metrics, u.scaler, u.notifier, u.upgradeConfig, logger)
	}

	return u.runSteps(ctx, logger, u.steps)
}

// shouldFailUpgrade checks if the cluster has reached a condition during upgrade
// where it should be treated as failed.
// If the cluster should fail its upgrade a condition of 'true' is returned.
// Any error encountered in making this decision is returned.
func hcpShouldFailUpgrade(cvClient cv.ClusterVersion, cfg *upgraderConfig, upgradeConfig *upgradev1alpha1.UpgradeConfig) (bool, error) {
	commenced, err := cvClient.HasUpgradeCommenced(upgradeConfig)
	if err != nil {
		return false, err
	}
	// If the upgrade has commenced, there's no going back
	if commenced {
		return false, nil
	}

	// Get the managed upgrade start time from upgrade config history
	h := upgradeConfig.Status.History.GetHistory(upgradeConfig.Spec.Desired.Version)
	if h == nil {
		return false, nil
	}
	startTime := h.StartTime.Time

	upgradeWindowDuration := cfg.UpgradeWindow.GetUpgradeWindowTimeOutDuration()
	if !startTime.IsZero() && upgradeWindowDuration > 0 && time.Now().After(startTime.Add(upgradeWindowDuration)) {
		return true, nil
	}
	return false, nil
}

// performUpgradeFailure carries out routines related to moving to an upgrade-failed state
func hcpPerformUpgradeFailure(c client.Client, metricsClient metrics.Metrics, s scaler.Scaler, nc eventmanager.EventManager, upgradeConfig *upgradev1alpha1.UpgradeConfig, logger logr.Logger) (upgradev1alpha1.UpgradePhase, error) {
	// Set up return condition
	h := upgradeConfig.Status.History.GetHistory(upgradeConfig.Spec.Desired.Version)
	condition := &upgradev1alpha1.UpgradeCondition{
		Type:    "FailedUpgrade",
		Status:  corev1.ConditionFalse,
		Reason:  "Upgrade failed",
		Message: "FailedUpgrade notification sent",
	}

	// TearDown the extra machineset
	_, err := s.EnsureScaleDownNodes(c, nil, logger)
	if err != nil {
		logger.Error(err, "Failed to scale down the temporary upgrade machine when upgrade failed")
		h.Conditions.SetCondition(*condition)
		return h.Phase, nil
	}

	// Notify of failure
	err = nc.Notify(notifier.MuoStateFailed)
	if err != nil {
		logger.Error(err, "Failed to notify of upgrade failure")
		h.Conditions.SetCondition(*condition)
		return h.Phase, nil
	}

	// flag window breached metric
	metricsClient.UpdateMetricUpgradeWindowBreached(upgradeConfig.Name)

	// cancel previously triggered metrics
	metricsClient.ResetFailureMetrics()

	// Update condition state to successful
	condition.Status = corev1.ConditionTrue

	h.Conditions.SetCondition(*condition)

	return upgradev1alpha1.UpgradePhaseFailed, nil
}
