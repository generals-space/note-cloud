package v1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 添加`+genclient:noStatus`标记可以不添加自定义资源的Status成员.
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MyCrd describes a MyCrd resource
type MyCrd struct {
	// TypeMeta为各资源通用元信息, 包括kind和apiVersion.
    metav1.TypeMeta `json:",inline"`
	// ObjectMeta为特定类型的元信息, 包括name, namespace, selfLink, labels等.
    metav1.ObjectMeta `json:"metadata,omitempty"`
    // spec字段
	Spec MyCrdSpec `json:"spec"`
	// status字段
	Status MyCrdStatus `json:"status"`
}

// MyCrdSpec is the spec for a MyResource resource
type MyCrdSpec struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MyCrdList is a list of MyCrd resources
type MyCrdList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata"`

    Items []MyCrd `json:"items"`
}

// MyCrdStatus is the status for a PodStatus resource
type MyCrdStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}
