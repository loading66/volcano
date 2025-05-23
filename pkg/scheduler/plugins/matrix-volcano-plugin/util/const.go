/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

// Package util is used for matrix-volcano-plugin
package util

const (
	// ShareDomain matrix schedule policy value
	ShareDomain = "shareDomain"
	// BinPack matrix schedule policy value
	BinPack = "binPack"
	// Spread matrix schedule policy value
	Spread = "spread"
	// MatrixSchedulePolicyKey matrix schedule policy key in yaml
	MatrixSchedulePolicyKey = "matrixSchedulePolicy"
	// BinPackWeightCPUKey bin pack weight CPU key in yaml
	BinPackWeightCPUKey = "matrixBinpackCPUWeight"
	// BinPackWeightMemoryKey bin pack weight memory key in yaml
	BinPackWeightMemoryKey = "matrixBinpackMemWeight"
	// MatrixScheduleSuccess The RackSelectionStatus recorded when scheduling fails is false.
	MatrixScheduleSuccess = true
	// MatrixScheduleFail The RackSelectionStatus recorded when scheduling fails is false.
	MatrixScheduleFail = false
	// MatrixMaxScore matrix get this score when be chosen
	MatrixMaxScore = 10000
)

const (
	// NumOne the number 1
	NumOne = 1
	// NumFive the number 5
	NumFive = 5
	// NumTen the number 10
	NumTen = 10
	// NumThirty the number 30
	NumThirty = 30
	// NumFifty the number 50
	NumFifty = 50
	// NumHundred the number 100
	NumHundred = 100
)

const (
	// LogErrorLev for error information.
	LogErrorLev = 1
	// LogWarningLev for warning information.
	LogWarningLev = 2
	// LogInfoLev for Info information.
	LogInfoLev = 3
	// LogDebugLev for debug information.
	LogDebugLev = 4
)
