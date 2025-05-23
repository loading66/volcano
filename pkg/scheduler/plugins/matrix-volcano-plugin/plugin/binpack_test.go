/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

// Package plugin is used for matrix-volcano-plugin
package plugin

import (
	"sync"
	"testing"

	"volcano.sh/volcano/pkg/scheduler/api"
	"volcano.sh/volcano/pkg/scheduler/framework"
	"volcano.sh/volcano/pkg/scheduler/plugins/matrix-volcano-plugin/util"
)

type MatrixBinPackArgs struct {
	ssn              *framework.Session
	task             *api.TaskInfo
	nodeInfo         []*api.NodeInfo
	nodeRankMapCache map[string]string
}

type MatrixBinPackWant struct {
	score   map[string]float64
	rackId  string
	wantErr bool
}

type MatrixBinPackTest struct {
	name string
	tp   *MatrixBinPackInfo
	args MatrixBinPackArgs
	want MatrixBinPackWant
}

func buildMatrixBinPackInfo(weightCPU, weightMemory int, jobId api.JobID, status bool,
	rackId string) *MatrixBinPackInfo {
	return &MatrixBinPackInfo{
		MatrixInfo: &util.MatrixInfo{
			Policy: util.BinPack,
			Weight: util.PriorityWeight{
				BinPackCPU:    weightCPU,
				BinPackMemory: weightMemory,
			},
			CurrJobId:           jobId,
			RackSelectionStatus: status,
			ChosenRackId:        rackId,
		},
	}
}

func buildResource(cpu, memory float64) *api.Resource {
	return &api.Resource{
		MilliCPU: cpu,
		Memory:   memory,
	}
}

func matrixBinPackTestCase1() MatrixBinPackTest {
	return MatrixBinPackTest{
		name: "MatrixBinPack success chosen rack 1",
		tp:   buildMatrixBinPackInfo(util.NumFive, util.NumOne, "", false, ""),
		args: MatrixBinPackArgs{
			ssn: &framework.Session{
				Jobs: map[api.JobID]*api.JobInfo{
					"job1": {
						Tasks: map[api.TaskID]*api.TaskInfo{
							"task1": {Resreq: buildResource(util.NumTen, util.NumTen)},
							"task2": {Resreq: buildResource(util.NumTen, util.NumTen)},
							"task3": {Resreq: buildResource(util.NumTen, util.NumTen)},
							"task4": {Resreq: buildResource(util.NumTen, util.NumTen)},
						},
					},
				},
			},
			task: &api.TaskInfo{Job: "job1", Name: "task1"},
			nodeInfo: []*api.NodeInfo{
				{
					Name:        "node1",
					Used:        buildResource(util.NumFifty, util.NumThirty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
				{
					Name:        "node2",
					Used:        buildResource(util.NumFifty, util.NumThirty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
				{
					Name:        "node3",
					Used:        buildResource(util.NumThirty, util.NumFifty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
				{
					Name:        "node4",
					Used:        buildResource(util.NumThirty, util.NumFifty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
			},
			nodeRankMapCache: map[string]string{"node1": "1", "node2": "1", "node3": "2", "node4": "2"},
		},
		want: MatrixBinPackWant{
			score:   map[string]float64{"node1": util.MatrixMaxScore, "node2": util.MatrixMaxScore, "node3": 0, "node4": 0},
			rackId:  "1",
			wantErr: false,
		},
	}
}

func matrixBinPackTestCase2() MatrixBinPackTest {
	return MatrixBinPackTest{
		name: "MatrixBinPack success chosen rack 2",
		tp:   buildMatrixBinPackInfo(util.NumOne, util.NumFive, "", false, ""),
		args: MatrixBinPackArgs{
			ssn: &framework.Session{
				Jobs: map[api.JobID]*api.JobInfo{
					"job1": {
						Tasks: map[api.TaskID]*api.TaskInfo{
							"task1": {Resreq: buildResource(util.NumTen, util.NumTen)},
							"task2": {Resreq: buildResource(util.NumTen, util.NumTen)},
							"task3": {Resreq: buildResource(util.NumTen, util.NumTen)},
							"task4": {Resreq: buildResource(util.NumTen, util.NumTen)},
						},
					},
				},
			},
			task: &api.TaskInfo{Job: "job1", Name: "task1"},
			nodeInfo: []*api.NodeInfo{
				{
					Name:        "node1",
					Used:        buildResource(util.NumFifty, util.NumThirty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
				{
					Name:        "node2",
					Used:        buildResource(util.NumFifty, util.NumThirty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
				{
					Name:        "node3",
					Used:        buildResource(util.NumThirty, util.NumFifty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
				{
					Name:        "node4",
					Used:        buildResource(util.NumThirty, util.NumFifty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
			},
			nodeRankMapCache: map[string]string{"node1": "1", "node2": "1", "node3": "2", "node4": "2"},
		},
		want: MatrixBinPackWant{
			score:   map[string]float64{"node1": 0, "node2": 0, "node3": util.MatrixMaxScore, "node4": util.MatrixMaxScore},
			rackId:  "2",
			wantErr: false,
		},
	}
}

func matrixBinPackTestCase3() MatrixBinPackTest {
	return MatrixBinPackTest{
		name: "MatrixBinPack success chosen rack 1 when RackSelectionStatus is true",
		tp:   buildMatrixBinPackInfo(util.NumOne, util.NumFive, "job1", true, "1"),
		args: MatrixBinPackArgs{
			ssn: &framework.Session{
				Jobs: map[api.JobID]*api.JobInfo{
					"job1": {
						Tasks: map[api.TaskID]*api.TaskInfo{
							"task1": {Resreq: buildResource(util.NumTen, util.NumTen)},
							"task2": {Resreq: buildResource(util.NumTen, util.NumTen)},
							"task3": {Resreq: buildResource(util.NumTen, util.NumTen)},
							"task4": {Resreq: buildResource(util.NumTen, util.NumTen)},
						},
					},
				},
			},
			task: &api.TaskInfo{Job: "job1", Name: "task1"},
			nodeInfo: []*api.NodeInfo{
				{
					Name:        "node1",
					Used:        buildResource(util.NumFifty, util.NumThirty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
				{
					Name:        "node2",
					Used:        buildResource(util.NumFifty, util.NumThirty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
				{
					Name:        "node3",
					Used:        buildResource(util.NumThirty, util.NumFifty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
				{
					Name:        "node4",
					Used:        buildResource(util.NumThirty, util.NumFifty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
			},
			nodeRankMapCache: map[string]string{"node1": "1", "node2": "1", "node3": "2", "node4": "2"},
		},
		want: MatrixBinPackWant{
			score:   map[string]float64{"node1": util.MatrixMaxScore, "node2": util.MatrixMaxScore, "node3": 0, "node4": 0},
			rackId:  "1",
			wantErr: false,
		},
	}
}

func buildMatrixBinPackTestCases() []MatrixBinPackTest {
	MatrixBinPackTests := []MatrixBinPackTest{
		matrixBinPackTestCase1(),
		matrixBinPackTestCase2(),
		matrixBinPackTestCase3(),
	}
	return MatrixBinPackTests
}

func TestMatrixBinPack(t *testing.T) {
	buildTests := buildMatrixBinPackTestCases()
	cacheLock = new(sync.RWMutex)
	for _, tt := range buildTests {
		t.Run(tt.name, func(t *testing.T) {
			nodeRankMapCache = tt.args.nodeRankMapCache
			res, err := MatrixBinPack(tt.tp.MatrixInfo, tt.args.ssn, tt.args.task, tt.args.nodeInfo)
			if (err != nil) != tt.want.wantErr {
				t.Errorf("MatrixBinPack() err = %v, wantErr %v", err, tt.want.wantErr)
			}
			if res["node1"] != tt.want.score["node1"] || res["node2"] != tt.want.score["node2"] ||
				res["node3"] != tt.want.score["node3"] || res["node4"] != tt.want.score["node4"] {
				t.Errorf("MatrixBinPack() res = %v, want %v", res, tt.want.score)
			}
			if tt.tp.MatrixInfo.ChosenRackId != tt.want.rackId {
				t.Errorf("MatrixBinPack() ChosenRackID = %v, want %v",
					tt.tp.MatrixInfo.ChosenRackId, tt.want.rackId)
			}
		})
	}
}

func matrixBinPackFailTestCase1() MatrixBinPackTest {
	return MatrixBinPackTest{
		name: "MatrixBinPack fail when rackID not found",
		tp:   buildMatrixBinPackInfo(util.NumFive, util.NumOne, "", false, ""),
		args: MatrixBinPackArgs{
			nodeInfo: []*api.NodeInfo{
				{Name: "node1"},
			},
			nodeRankMapCache: map[string]string{},
		},
		want: MatrixBinPackWant{wantErr: true},
	}
}

func matrixBinPackFailTestCase2() MatrixBinPackTest {
	return MatrixBinPackTest{
		name: "MatrixBinPack fail when RackSelectionStatus is false",
		tp:   buildMatrixBinPackInfo(util.NumOne, util.NumFive, "job1", false, ""),
		args: MatrixBinPackArgs{
			task:             &api.TaskInfo{Job: "job1", Name: "task1"},
			nodeRankMapCache: map[string]string{"node1": "1"},
		},
		want: MatrixBinPackWant{wantErr: true},
	}
}

func matrixBinPackFailTestCase3() MatrixBinPackTest {
	return MatrixBinPackTest{
		name: "MatrixBinPack fail when resource not enough",
		tp:   buildMatrixBinPackInfo(util.NumOne, util.NumFive, "", false, ""),
		args: MatrixBinPackArgs{
			ssn: &framework.Session{
				Jobs: map[api.JobID]*api.JobInfo{
					"job1": {
						Tasks: map[api.TaskID]*api.TaskInfo{
							"task1": {Resreq: buildResource(util.NumTen, util.NumTen)},
						},
					},
				},
			},
			task: &api.TaskInfo{
				Job:  "job1",
				Name: "task1",
			},
			nodeInfo: []*api.NodeInfo{
				{
					Name:        "node1",
					Used:        buildResource(util.NumHundred, util.NumHundred),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
			},
			nodeRankMapCache: map[string]string{"node1": "1"},
		},
		want: MatrixBinPackWant{wantErr: true},
	}
}

func matrixBinPackFailTestCase4() MatrixBinPackTest {
	return MatrixBinPackTest{
		name: "MatrixBinPack fail when MatrixInfo is nil",
		tp: &MatrixBinPackInfo{
			MatrixInfo: nil,
		},
		want: MatrixBinPackWant{wantErr: true},
	}
}

func matrixBinPackFailTestCase5() MatrixBinPackTest {
	return MatrixBinPackTest{
		name: "MatrixBinPack fail when binPackRecord not found rackId",
		tp:   buildMatrixBinPackInfo(util.NumOne, util.NumFive, "job1", true, "1"),
		args: MatrixBinPackArgs{
			task: &api.TaskInfo{
				Job:  "job1",
				Name: "task1",
			},
			nodeInfo: []*api.NodeInfo{
				{
					Name:        "node1",
					Used:        buildResource(util.NumFifty, util.NumThirty),
					Allocatable: buildResource(util.NumHundred, util.NumHundred),
				},
			},
			nodeRankMapCache: map[string]string{"node1": "2"},
		},
		want: MatrixBinPackWant{wantErr: true},
	}
}

func buildMatrixBinPackFailTestCases() []MatrixBinPackTest {
	MatrixBinPackFailTests := []MatrixBinPackTest{
		matrixBinPackFailTestCase1(),
		matrixBinPackFailTestCase2(),
		matrixBinPackFailTestCase3(),
		matrixBinPackFailTestCase4(),
		matrixBinPackFailTestCase5(),
	}
	return MatrixBinPackFailTests
}

func TestMatrixBinPackFail(t *testing.T) {
	buildTests := buildMatrixBinPackFailTestCases()
	cacheLock = new(sync.RWMutex)
	for _, tt := range buildTests {
		t.Run(tt.name, func(t *testing.T) {
			nodeRankMapCache = tt.args.nodeRankMapCache
			_, err := MatrixBinPack(tt.tp.MatrixInfo, tt.args.ssn, tt.args.task, tt.args.nodeInfo)
			if (err != nil) != tt.want.wantErr {
				t.Errorf("MatrixBinPack() err = %v, wantErr %v", err, tt.want.wantErr)
			}
		})
	}
}

type SelectSuitableRackIdTest struct {
	name    string
	tp      *MatrixBinPackInfo
	wantErr bool
}

func buildSelectSuitableRackIdTestCases() []SelectSuitableRackIdTest {
	GetMatrixInfoTests := []SelectSuitableRackIdTest{
		{
			name: "SelectSuitableRackId fail when MatrixAllocateCPU is 0",
			tp: &MatrixBinPackInfo{
				BinPackRecord: map[string]*util.BinPackInfo{
					"rack1": {
						IsSchedulable:        true,
						MatrixAllocateCPU:    0,
						MatrixAllocateMemory: util.NumOne,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "SelectSuitableRackId fail when MatrixAllocateMemory is 0",
			tp: &MatrixBinPackInfo{
				BinPackRecord: map[string]*util.BinPackInfo{
					"rack1": {
						IsSchedulable:        true,
						MatrixAllocateCPU:    util.NumOne,
						MatrixAllocateMemory: 0,
					},
				},
			},
			wantErr: true,
		},
	}
	return GetMatrixInfoTests
}

func TestSelectSuitableRackId(t *testing.T) {
	buildTests := buildSelectSuitableRackIdTestCases()
	for _, tt := range buildTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.tp.SelectSuitableRackId(); (res != nil) != tt.wantErr {
				t.Errorf("SelectSuitableRackId() err = %v, wantErr %v", res, tt.wantErr)
			}
		})
	}
}

type GroupNodesByRackIDTest struct {
	name    string
	tp      *MatrixBinPackInfo
	args    []*api.NodeInfo
	wantErr bool
}

func buildGroupNodesByRackIDTestCases() []GroupNodesByRackIDTest {
	GroupNodesByRackIDTests := []GroupNodesByRackIDTest{
		{
			name: "GroupNodesByRackID fail when node allocate CPU is 0",
			tp:   &MatrixBinPackInfo{},
			args: []*api.NodeInfo{
				{
					Name:        "node1",
					Allocatable: buildResource(0, util.NumOne),
				},
			},
			wantErr: true,
		},
		{
			name: "GroupNodesByRackID fail when node allocate memory is 0",
			tp:   &MatrixBinPackInfo{},
			args: []*api.NodeInfo{
				{
					Name:        "node1",
					Allocatable: buildResource(util.NumOne, 0),
				},
			},
			wantErr: true,
		},
	}
	return GroupNodesByRackIDTests
}

func TestGroupNodesByRackID(t *testing.T) {
	buildTests := buildGroupNodesByRackIDTestCases()
	cacheLock = new(sync.RWMutex)
	for _, tt := range buildTests {
		t.Run(tt.name, func(t *testing.T) {
			nodeRankMapCache = map[string]string{"node1": "rack1"}
			if res := tt.tp.GroupNodesByRackID(tt.args); (res != nil) != tt.wantErr {
				t.Errorf("GroupNodesByRackID() err = %v, wantErr %v", res, tt.wantErr)
			}
		})
	}
}
