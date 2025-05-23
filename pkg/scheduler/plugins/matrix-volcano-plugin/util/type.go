/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

// Package util is used for matrix-volcano-plugin
package util

import "volcano.sh/volcano/pkg/scheduler/api"

// PriorityWeight weight for plugin
type PriorityWeight struct {
	BinPackCPU    int
	BinPackMemory int
}

// TaskReqRes task request resource
type TaskReqRes struct {
	CPU float64
	Mem float64
}

// NodeRes node resource info
type NodeRes struct {
	AllocateCPU float64
	AllocateMem float64
	UsedCPU     float64
	UsedMem     float64
}

// BinPackInfo bin pack info
type BinPackInfo struct {
	MatrixAllocateCPU    float64
	MatrixAllocateMemory float64
	MatrixUsedCPU        float64
	MatrixUsedMemory     float64
	NodesResInfo         map[string]*NodeRes
	IsSchedulable        bool
}

// MatrixSpreadInfo matrix spread info
type MatrixSpreadInfo struct {
	MatrixInfo   *MatrixInfo
	SpreadRecord map[string][]string
}

// MatrixShareDomainInfo matrix share domain info
type MatrixShareDomainInfo struct {
	MatrixInfo *MatrixInfo
}

// MatrixInfo matrix info
type MatrixInfo struct {
	// Arguments given for the plugin
	Policy string
	Weight PriorityWeight
	// Info record by plugin
	CurrJobId           api.JobID
	RackSelectionStatus bool
	ChosenRackId        string
	ChosenShareDomain   string
}
