/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

// Package plugin is used for matrix-volcano-plugin
package plugin

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"volcano.sh/volcano/pkg/scheduler/api"
	"volcano.sh/volcano/pkg/scheduler/framework"
	"volcano.sh/volcano/pkg/scheduler/plugins/matrix-volcano-plugin/util"
)

// MatrixBinPackInfo matrix bin pack info
type MatrixBinPackInfo struct {
	MatrixInfo    *util.MatrixInfo
	BinPackRecord map[string]*util.BinPackInfo
}

// MatrixBinPack start matrix bin pack schedule
func MatrixBinPack(matrixInfo *util.MatrixInfo, ssn *framework.Session, task *api.TaskInfo,
	nodeInfo []*api.NodeInfo) (map[string]float64, error) {
	klog.V(util.LogDebugLev).Infof("Enter matrix %s.", util.BinPack)
	defer klog.V(util.LogDebugLev).Infof("Leave matrix %s.", util.BinPack)
	m, err := GetMatrixBinPackInfo(matrixInfo)
	if err != nil {
		klog.V(util.LogErrorLev).Infof("task:%s GetMatrixBinPackInfo err: %v.", task.Name, err)
		return nil, err
	}
	// node按照Rack分组
	if err = m.GroupNodesByRackID(nodeInfo); err != nil {
		klog.V(util.LogErrorLev).Infof("task:%s GroupNodesByRackID err: %v.", task.Name, err)
		return nil, err
	}
	// 当前task已经记录并且rack已经选出，直接按照已选rackID打分并返回
	if m.MatrixInfo.CurrJobId == task.Job && m.MatrixInfo.RackSelectionStatus {
		nodeScore, err := m.ScoreMatrixNode(nodeInfo)
		if err != nil {
			klog.V(util.LogErrorLev).Infof("task:%s ScoreMatrixNode err: %v.", task.Name, err)
			return nil, err
		}
		return nodeScore, nil
	}
	// 当前task已经记录，但是rack未选出，调度失败
	if m.MatrixInfo.CurrJobId == task.Job && !m.MatrixInfo.RackSelectionStatus {
		return nil, errors.Errorf("failed to select a suitable matrix node")
	}
	// 当前task未记录，初始化记录并开始选择rack
	m.InitMatrixBinPackInfo(task)
	// 获取各个task资源req，从大到小排序
	// 遍历所有rack，依次找到能满足所有task的rack，
	if err = m.SelectMatrixForTasks(m.GetTaskList(ssn, task)); err != nil {
		klog.V(util.LogErrorLev).Infof("task:%s SelectMatrixForTasks err: %v.", task.Name, err)
		return nil, err
	}
	m.UpdateMatrixInfo(matrixInfo)
	if nodeScores, err := m.ScoreMatrixNode(nodeInfo); err != nil {
		klog.V(util.LogErrorLev).Infof("task:%s ScoreMatrixNode err: %v.", task.Name, err)
		return nil, err
	} else {
		klog.V(util.LogInfoLev).Infof("matrix %s score: %v.", util.BinPack, nodeScores)
		return nodeScores, nil
	}
}

// GetMatrixBinPackInfo get info for matrix bin pack
func GetMatrixBinPackInfo(matrixInfo *util.MatrixInfo) (*MatrixBinPackInfo, error) {
	if matrixInfo == nil {
		return nil, errors.Errorf(fmt.Sprintf("matrixInfo is nil"))
	}
	return &MatrixBinPackInfo{
		MatrixInfo:    matrixInfo,
		BinPackRecord: make(map[string]*util.BinPackInfo),
	}, nil
}

// InitMatrixBinPackInfo init info when job's first task start schedule
func (m *MatrixBinPackInfo) InitMatrixBinPackInfo(task *api.TaskInfo) {
	m.MatrixInfo.CurrJobId = task.Job
	m.MatrixInfo.RackSelectionStatus = util.MatrixScheduleFail
	m.MatrixInfo.ChosenRackId = ""
}

// UpdateMatrixInfo update info when finish
func (m *MatrixBinPackInfo) UpdateMatrixInfo(matrixInfo *util.MatrixInfo) {
	matrixInfo.ChosenRackId = m.MatrixInfo.ChosenRackId
	matrixInfo.CurrJobId = m.MatrixInfo.CurrJobId
	matrixInfo.RackSelectionStatus = m.MatrixInfo.RackSelectionStatus
	return
}

// GroupNodesByRackID use rack ID group all nodes
func (m *MatrixBinPackInfo) GroupNodesByRackID(nodeInfo []*api.NodeInfo) error {
	nodeGroup := make(map[string]*util.BinPackInfo)
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	for _, node := range nodeInfo {
		rackID, found := nodeRankMapCache[node.Name]
		if !found {
			return errors.Errorf(fmt.Sprintf("not found rackID in node %s", node.Name))
		}
		if node.Allocatable.MilliCPU == 0 || node.Allocatable.Memory == 0 {
			return errors.Errorf(fmt.Sprintf("node %s resource info is invaild", node.Name))
		}
		if _, found = nodeGroup[rackID]; !found {
			nodeGroup[rackID] = &util.BinPackInfo{
				MatrixAllocateCPU:    0,
				MatrixAllocateMemory: 0,
				MatrixUsedCPU:        0,
				MatrixUsedMemory:     0,
				NodesResInfo:         make(map[string]*util.NodeRes),
				IsSchedulable:        false,
			}
		}
		nodeGroup[rackID].MatrixAllocateCPU += node.Allocatable.MilliCPU
		nodeGroup[rackID].MatrixAllocateMemory += node.Allocatable.Memory
		nodeGroup[rackID].MatrixUsedCPU += node.Used.MilliCPU
		nodeGroup[rackID].MatrixUsedMemory += node.Used.Memory
		nodeGroup[rackID].NodesResInfo[node.Name] = &util.NodeRes{
			AllocateCPU: node.Allocatable.MilliCPU,
			AllocateMem: node.Allocatable.Memory,
			UsedCPU:     node.Used.MilliCPU,
			UsedMem:     node.Used.Memory,
		}
	}
	m.BinPackRecord = nodeGroup
	return nil
}

// GetTaskList get all tasks and sort
func (m *MatrixBinPackInfo) GetTaskList(ssn *framework.Session, task *api.TaskInfo) []util.TaskReqRes {
	taskList := make([]util.TaskReqRes, len(ssn.Jobs[task.Job].Tasks))
	taskNum := 0
	for _, currTask := range ssn.Jobs[task.Job].Tasks {
		taskList[taskNum].CPU = currTask.Resreq.MilliCPU
		taskList[taskNum].Mem = currTask.Resreq.Memory
		taskNum += 1
	}
	sort.Slice(taskList, func(i, j int) bool {
		if m.MatrixInfo.Weight.BinPackCPU >= m.MatrixInfo.Weight.BinPackMemory {
			return taskList[i].CPU > taskList[j].CPU
		}
		return taskList[i].Mem > taskList[j].Mem
	})
	return taskList
}

// SelectNodeForTask select all suitable node for task
func (m *MatrixBinPackInfo) SelectNodeForTask(rackID string, rackInfo *util.BinPackInfo, task util.TaskReqRes) {
	var usedFinallyCPU, usedFinallyMemory float64
	var minRemainNode = ""
	for nodeName, nodeResInfo := range rackInfo.NodesResInfo {
		usedFinallyCPU = nodeResInfo.UsedCPU + task.CPU
		usedFinallyMemory = nodeResInfo.UsedMem + task.Mem
		if usedFinallyCPU > nodeResInfo.AllocateCPU || usedFinallyMemory > nodeResInfo.AllocateMem {
			continue
		}
		minRemainNode = nodeName
		break
	}
	if minRemainNode != "" {
		rackInfo.NodesResInfo[minRemainNode].UsedMem += task.Mem
		rackInfo.NodesResInfo[minRemainNode].UsedCPU += task.CPU
		rackInfo.MatrixUsedCPU += task.CPU
		rackInfo.MatrixUsedMemory += task.Mem
		m.BinPackRecord[rackID].IsSchedulable = true
		return
	}
	m.BinPackRecord[rackID].IsSchedulable = false
}

// SelectSuitableRackId select rack ID for job
func (m *MatrixBinPackInfo) SelectSuitableRackId() error {
	var maxScore = 0.0
	var ChosenRackId = ""
	for rackID, rackInfo := range m.BinPackRecord {
		if !rackInfo.IsSchedulable {
			continue
		}
		if rackInfo.MatrixAllocateCPU == 0 || rackInfo.MatrixAllocateMemory == 0 {
			return errors.Errorf("matrix node %v resource info is invaild", rackID)
		}
		currMatrixScore := float64(m.MatrixInfo.Weight.BinPackCPU)*rackInfo.MatrixUsedCPU/rackInfo.MatrixAllocateCPU +
			float64(m.MatrixInfo.Weight.BinPackMemory)*rackInfo.MatrixUsedMemory/rackInfo.MatrixAllocateMemory
		if currMatrixScore > maxScore {
			maxScore = currMatrixScore
			ChosenRackId = rackID
		}
	}
	if ChosenRackId != "" {
		m.MatrixInfo.ChosenRackId = ChosenRackId
		m.MatrixInfo.RackSelectionStatus = util.MatrixScheduleSuccess
		return nil
	}
	m.MatrixInfo.RackSelectionStatus = util.MatrixScheduleFail
	return errors.Errorf("no matrix node could been choosen")
}

// SelectMatrixForTasks select matrix for tasks
func (m *MatrixBinPackInfo) SelectMatrixForTasks(taskList []util.TaskReqRes) error {
	for rackID, rackInfo := range m.BinPackRecord {
		for _, currTask := range taskList {
			m.SelectNodeForTask(rackID, rackInfo, currTask)
			if !m.BinPackRecord[rackID].IsSchedulable {
				break
			}
		}
	}
	if err := m.SelectSuitableRackId(); err != nil {
		return err
	}
	return nil
}

// ScoreMatrixNode score node for matrix bin pack
func (m *MatrixBinPackInfo) ScoreMatrixNode(nodeInfo []*api.NodeInfo) (map[string]float64, error) {
	nodeScores := make(map[string]float64, len(nodeInfo))
	if _, found := m.BinPackRecord[m.MatrixInfo.ChosenRackId]; !found {
		return nil, errors.Errorf(fmt.Sprintf("not found rackID %v in binPackRecord", m.MatrixInfo.ChosenRackId))
	}
	for nodeName := range m.BinPackRecord[m.MatrixInfo.ChosenRackId].NodesResInfo {
		nodeScores[nodeName] = util.MatrixMaxScore
	}
	return nodeScores, nil
}
