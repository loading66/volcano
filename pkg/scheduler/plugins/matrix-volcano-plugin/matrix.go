/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

// Package main is used for matrix-volcano-plugin
package main

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"volcano.sh/volcano/pkg/scheduler/api"
	"volcano.sh/volcano/pkg/scheduler/framework"
	"volcano.sh/volcano/pkg/scheduler/plugins/matrix-volcano-plugin/plugin"
	"volcano.sh/volcano/pkg/scheduler/plugins/matrix-volcano-plugin/util"
)

// New return matrix plugin.
func New(arguments framework.Arguments) framework.Plugin {
	return &MatrixVolcanoPlugin{
		Policy: "",
		Weight: util.PriorityWeight{
			BinPackCPU:    0,
			BinPackMemory: 0,
		},
		CurrJobId:           "",
		RackSelectionStatus: false,
		ChosenRackId:        "",
		ChosenShareDomain:   "",
	}
}

// Name This need by volcano frame init plugin.
func (m *MatrixVolcanoPlugin) Name() string {
	return MatrixPluginName
}

// OnSessionOpen open session for frame.
func (m *MatrixVolcanoPlugin) OnSessionOpen(ssn *framework.Session) {
	klog.V(util.LogDebugLev).Infof("Enter matrix-volcano-plugin ...")
	defer func() {
		klog.V(util.LogDebugLev).Infof("Leave matrix-volcano-plugin.")
	}()
	batchNodeOrderFn := func(task *api.TaskInfo, nodeInfo []*api.NodeInfo) (map[string]float64, error) {
		nodeScores := make(map[string]float64, len(nodeInfo))
		var err error
		if err = m.GetMatrixSchedulePolicyArg(ssn, task); err != nil {
			klog.V(util.LogErrorLev).Infof("task:%s GetMatrixSchedulePolicyArg err: %v.", task.Name, err)
			return nil, err
		}
		// 根据配置信息走不同分支
		switch m.Policy {
		case util.ShareDomain:
			// 拓扑感知调度Topology awareness
			break
		case util.BinPack:
			// 装箱调度Bin Pack
			if nodeScores, err = plugin.MatrixBinPack(m.GetMatrixInfo(), ssn, task, nodeInfo); err != nil {
				klog.V(util.LogErrorLev).Infof("task: %v matrix bin pack err: %v", task.Name, err)
				return nil, err
			}
			break
		case util.Spread:
			// Spread调度
			break
		default:
			// 配置信息异常
			klog.V(util.LogErrorLev).Infof("task: %v invalid matrix schedule policy %s", task.Name, m.Policy)
			return nil, errors.Errorf(fmt.Sprintf("invalid matrix schedule policy %s", m.Policy))
		}

		return nodeScores, nil
	}
	ssn.AddBatchNodeOrderFn(m.Name(), batchNodeOrderFn)
}

// GetMatrixInfo get info for matrix volcano plugin
func (m *MatrixVolcanoPlugin) GetMatrixInfo() *util.MatrixInfo {
	return &util.MatrixInfo{
		Policy: m.Policy,
		Weight: util.PriorityWeight{
			BinPackCPU:    m.Weight.BinPackCPU,
			BinPackMemory: m.Weight.BinPackMemory,
		},
		CurrJobId:           m.CurrJobId,
		RackSelectionStatus: m.RackSelectionStatus,
		ChosenRackId:        m.ChosenRackId,
		ChosenShareDomain:   m.ChosenShareDomain,
	}
}

// GetMatrixWeight get weight for matrix volcano plugin
func (m *MatrixVolcanoPlugin) GetMatrixWeight(job *api.JobInfo) error {
	var binPackCPU, binPackMem = 1, 1
	var err error
	if weightCPU, found := job.PodGroup.Labels[util.BinPackWeightCPUKey]; found {
		binPackCPU, err = strconv.Atoi(weightCPU)
		if err != nil || binPackCPU <= 0 {
			return errors.Errorf(fmt.Sprintf("invalid matrix plugin cpu weight %s", weightCPU))
		}
	}
	if weightMem, found := job.PodGroup.Labels[util.BinPackWeightMemoryKey]; found {
		binPackMem, err = strconv.Atoi(weightMem)
		if err != nil || binPackMem <= 0 {
			return errors.Errorf(fmt.Sprintf("invalid matrix plugin memory weight %s", weightMem))
		}
	}
	m.Weight.BinPackMemory = binPackMem
	m.Weight.BinPackCPU = binPackCPU
	return nil
}

// GetMatrixSchedulePolicyArg get args for matrix volcano plugin
func (m *MatrixVolcanoPlugin) GetMatrixSchedulePolicyArg(ssn *framework.Session, task *api.TaskInfo) error {
	job, found := ssn.Jobs[task.Job]
	if !found {
		return errors.Errorf("not found job in ssn")
	}
	policy, found := job.PodGroup.Labels[util.MatrixSchedulePolicyKey]
	if !found {
		return errors.Errorf("not found matrixSchedulePolicy in jobinfo")
	}
	if err := m.GetMatrixWeight(job); err != nil {
		return err
	}
	m.Policy = policy
	return nil
}

// OnSessionClose Close session by volcano frame.
func (m *MatrixVolcanoPlugin) OnSessionClose(ssn *framework.Session) {
}
