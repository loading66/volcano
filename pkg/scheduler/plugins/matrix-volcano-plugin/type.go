/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

// Package main is used for matrix-volcano-plugin
package main

import (
	"volcano.sh/volcano/pkg/scheduler/api"
	"volcano.sh/volcano/pkg/scheduler/plugins/matrix-volcano-plugin/util"
)

// MatrixPluginName indicates name of volcano scheduler plugin.
const MatrixPluginName = "matrix-volcano-plugin"

// MatrixVolcanoPlugin is the struct of matrix-volcano-plugin.
type MatrixVolcanoPlugin struct {
	// Arguments given for the plugin
	Policy string
	Weight util.PriorityWeight
	// Info record by plugin
	CurrJobId           api.JobID
	RackSelectionStatus bool
	ChosenRackId        string
	ChosenShareDomain   string
}
