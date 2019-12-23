package v1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 添加`+genclient:noStatus`标记可以不添加自定义资源的Status成员.
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodGroup describes a PodGroup resource
type PodGroup struct {
	// TypeMeta为各资源通用元信息, 包括kind和apiVersion.
    metav1.TypeMeta `json:",inline"`
	// ObjectMeta为特定类型的元信息, 包括name, namespace, selfLink, labels等.
    metav1.ObjectMeta `json:"metadata,omitempty"`
    // spec字段
	Spec PodGroupSpec `json:"spec"`
	// status字段
	Status PodGroupStatus `json:"status"`
}

// PodGroupSpec is the spec for a MyResource resource
type PodGroupSpec struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodGroupList is a list of PodGroup resources
type PodGroupList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata"`

    Items []PodGroup `json:"items"`
}

// PodGroupStatus is the status for a PodStatus resource
type PodGroupStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}
