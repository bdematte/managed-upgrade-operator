package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HypershiftUpgradeType provides a type to declare upgrade types with
type HypershiftUpgradeType string

// UpgradeConfigSpec defines the desired state of UpgradeConfig and upgrade window and freeze window
type HypershiftUpgradeConfigSpec struct {
	// Specify the desired OpenShift release
	Desired Update `json:"desired"`

	// Specify the upgrade start time
	UpgradeAt string `json:"upgradeAt"`

	// +kubebuilder:validation:Minimum:=0
	// The maximum grace period granted to a node whose drain is blocked by a Pod Disruption Budget, before that drain is forced. Measured in minutes. The minimum accepted value is 0 and in this case it will trigger force drain after the expectedNodeDrainTime lapsed.
	PDBForceDrainTimeout int32 `json:"PDBForceDrainTimeout"`

	// +kubebuilder:validation:Enum={"OSD","ARO"}
	// Type indicates the HypershiftClusterUpgrader implementation to use to perform an upgrade of the cluster
	Type UpgradeType `json:"type"`

	// Specify if scaling up an extra node for capacity reservation before upgrade starts is needed
	CapacityReservation bool `json:"capacityReservation,omitempty"`
}

// +kubebuilder:object:root=true

// UpgradeConfig is the Schema for the upgradeconfigs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=hypershiftupgradeconfigs,scope=Namespaced,shortName=hyperupgrade
// +kubebuilder:printcolumn:name="desired_version",type="string",JSONPath=".spec.desired.version"
// +kubebuilder:printcolumn:name="phase",type="string",JSONPath=".status.history[0].phase"
// +kubebuilder:printcolumn:name="stage",type="string",JSONPath=".status.history[0].conditions[0].type"
// +kubebuilder:printcolumn:name="status",type="string",JSONPath=".status.history[0].conditions[0].status"
// +kubebuilder:printcolumn:name="reason",type="string",JSONPath=".status.history[0].conditions[0].reason"
// +kubebuilder:printcolumn:name="message",type="string",JSONPath=".status.history[0].conditions[0].message"
type HypershiftUpgradeConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HypershiftUpgradeConfigSpec   `json:"spec,omitempty"`
	Status UpgradeConfigStatus `json:"status,omitempty"`
}
