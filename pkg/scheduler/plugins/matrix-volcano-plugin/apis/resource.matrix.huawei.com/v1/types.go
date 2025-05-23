package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:group=resource.matrix.huawei.com,version=v1

type MatrixMetric struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +optional
	// Spec defines a specification of a volume owned by the cluster.
	Spec []Metric `json:"spec"`
}

// Metric defines single item
type Metric struct {
	MetricName string       `json:"metricName"` // 指标名称
	Timestamp  int64        `json:"timestamp"`  // 时间戳
	Metrics    []MetricItem `json:"metrics"`    // 指标项数组
}

// MetricItem is the specification of a volume.
type MetricItem struct {
	NumaId      string     `json:"numaId,omitempty"`      // NUMA 节点 ID
	CpuCount    string     `json:"cpuCount,omitempty"`    // CPU 核心数
	Total       *HugePages `json:"total,omitempty"`       // HUGE 页总内存
	Allocatable *HugePages `json:"allocatable,omitempty"` // HUGE 页可分配内存
	LocalIp     string     `json:"localIp,omitempty"`     // 本地 IP
	LocalName   string     `json:"localName,omitempty"`   // 本地名称
	Target      []Target   `json:"target,omitempty"`      // 目标信息
}

// HugePages defines HUGE 页内存的结构
type HugePages struct {
	Hugepages1Gi  string `json:"hugepages-1Gi,omitempty"`  // 1Gi HUGE 页大小的内存
	Hugepages2Mi  string `json:"hugepages-2Mi,omitempty"`  // 2Mi HUGE 页大小的内存
	Hugepages32Mi string `json:"hugepages-32Mi,omitempty"` // 32Mi HUGE 页大小的内存
	Hugepages64Ki string `json:"hugepages-64Ki,omitempty"` // 64Ki HUGE 页大小的内存
}

func (hp *HugePages) Fields() []string {
	return []string{"Hugepages1Gi", "Hugepages2Mi", "Hugepages32Mi", "Hugepages64Ki"}
}

// Target defines ip target
type Target struct {
	TargetIps []string `json:"TargetIps,omitempty"` // 目标 IP 列表
	TtlNums   uint32   `json:"TtlNums,omitempty"`   // TTL 数量
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MatrixMetricList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []MatrixMetric `json:"items"`
}
