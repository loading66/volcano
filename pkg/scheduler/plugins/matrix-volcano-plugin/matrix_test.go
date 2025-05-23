/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

// Package main is used for matrix-volcano-plugin
package main

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"volcano.sh/apis/pkg/apis/scheduling"
	"volcano.sh/volcano/pkg/scheduler/api"
	"volcano.sh/volcano/pkg/scheduler/framework"
	"volcano.sh/volcano/pkg/scheduler/plugins/matrix-volcano-plugin/util"
)

type nameTest struct {
	name string
	tp   *MatrixVolcanoPlugin
	want string
}

func buildNameTestCases() []nameTest {
	nameTests := []nameTest{
		{
			name: "Name success test",
			tp:   &MatrixVolcanoPlugin{},
			want: MatrixPluginName,
		},
	}
	return nameTests
}

func TestName(t *testing.T) {
	buildTests := buildNameTestCases()
	for _, tt := range buildTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.tp.Name(); res != tt.want {
				t.Errorf("Name() res = %v, want %v", res, tt.want)
			}
		})
	}
}

type newTest struct {
	name string
	args framework.Arguments
	want *MatrixVolcanoPlugin
}

func buildNewTestCases() []newTest {
	newTests := []newTest{
		{
			name: "New success test",
			args: make(framework.Arguments),
			want: &MatrixVolcanoPlugin{},
		},
	}
	return newTests
}

func TestNew(t *testing.T) {
	buildTests := buildNewTestCases()
	for _, tt := range buildTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := New(tt.args); res == nil {
				t.Errorf("Name() res = %v, want %v", res, tt.want)
			}
		})
	}
}

type OnSessionOpenTest struct {
	name string
	tp   *MatrixVolcanoPlugin
	args *framework.Session
}

func buildOnSessionOpenTestCases() []OnSessionOpenTest {
	OnSessionOpenTests := []OnSessionOpenTest{
		{
			name: "OnSessionOpen success test",
			tp:   &MatrixVolcanoPlugin{},
			args: &framework.Session{
				BatchNodeOrderFns: map[string]api.BatchNodeOrderFn{
					MatrixPluginName: nil,
				},
			},
		},
	}
	return OnSessionOpenTests
}

func TestOnSessionOpen(t *testing.T) {
	buildTests := buildOnSessionOpenTestCases()
	for _, tt := range buildTests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tp.OnSessionOpen(tt.args)
		})
	}
}

type GetMatrixSchedulePolicyArgArgs struct {
	ssn  *framework.Session
	task *api.TaskInfo
}

type GetMatrixSchedulePolicyArgWant struct {
	policy  string
	weight  *util.PriorityWeight
	wantErr bool
}

type GetMatrixSchedulePolicyArgTest struct {
	name string
	tp   *MatrixVolcanoPlugin
	args GetMatrixSchedulePolicyArgArgs
	want GetMatrixSchedulePolicyArgWant
}

func getMatrixSchedulePolicyArgTestCase1() GetMatrixSchedulePolicyArgTest {
	return GetMatrixSchedulePolicyArgTest{
		name: "GetMatrixSchedulePolicyArg success test",
		tp:   &MatrixVolcanoPlugin{},
		args: GetMatrixSchedulePolicyArgArgs{
			ssn: &framework.Session{
				Jobs: map[api.JobID]*api.JobInfo{
					"job1": {
						PodGroup: &api.PodGroup{
							PodGroup: scheduling.PodGroup{
								ObjectMeta: v1.ObjectMeta{
									Labels: map[string]string{
										util.MatrixSchedulePolicyKey: util.BinPack,
										util.BinPackWeightCPUKey:     "1",
										util.BinPackWeightMemoryKey:  "1",
									},
								},
							},
						},
					},
				},
			},
			task: &api.TaskInfo{Job: "job1"},
		},
		want: GetMatrixSchedulePolicyArgWant{
			policy: util.BinPack,
			weight: &util.PriorityWeight{
				BinPackCPU:    1,
				BinPackMemory: 1,
			},
			wantErr: false,
		},
	}
}

func getMatrixSchedulePolicyArgTestCase2() GetMatrixSchedulePolicyArgTest {
	return GetMatrixSchedulePolicyArgTest{
		name: "GetMatrixSchedulePolicyArg fail test when job not found",
		tp:   &MatrixVolcanoPlugin{},
		args: GetMatrixSchedulePolicyArgArgs{
			ssn: &framework.Session{
				Jobs: map[api.JobID]*api.JobInfo{
					"job1": {
						PodGroup: &api.PodGroup{
							PodGroup: scheduling.PodGroup{
								ObjectMeta: v1.ObjectMeta{
									Labels: map[string]string{
										util.MatrixSchedulePolicyKey: util.BinPack,
										util.BinPackWeightCPUKey:     "1",
										util.BinPackWeightMemoryKey:  "1",
									},
								},
							},
						},
					},
				},
			},
			task: &api.TaskInfo{Job: "job2"},
		},
		want: GetMatrixSchedulePolicyArgWant{
			policy:  "",
			weight:  &util.PriorityWeight{},
			wantErr: true,
		},
	}
}

func getMatrixSchedulePolicyArgTestCase3() GetMatrixSchedulePolicyArgTest {
	return GetMatrixSchedulePolicyArgTest{
		name: "GetMatrixSchedulePolicyArg fail test when MatrixSchedulePolicyKey not found",
		tp:   &MatrixVolcanoPlugin{},
		args: GetMatrixSchedulePolicyArgArgs{
			ssn: &framework.Session{
				Jobs: map[api.JobID]*api.JobInfo{
					"job1": {
						PodGroup: &api.PodGroup{
							PodGroup: scheduling.PodGroup{
								ObjectMeta: v1.ObjectMeta{
									Labels: map[string]string{
										"policy":                    util.BinPack,
										util.BinPackWeightCPUKey:    "1",
										util.BinPackWeightMemoryKey: "1",
									},
								},
							},
						},
					},
				},
			},
			task: &api.TaskInfo{Job: "job1"},
		},
		want: GetMatrixSchedulePolicyArgWant{
			policy:  "",
			weight:  &util.PriorityWeight{},
			wantErr: true,
		},
	}
}

func getMatrixSchedulePolicyArgTestCase4() GetMatrixSchedulePolicyArgTest {
	return GetMatrixSchedulePolicyArgTest{
		name: "GetMatrixSchedulePolicyArg fail test when BinPackWeightCPU is 0",
		tp:   &MatrixVolcanoPlugin{},
		args: GetMatrixSchedulePolicyArgArgs{
			ssn: &framework.Session{
				Jobs: map[api.JobID]*api.JobInfo{
					"job1": {
						PodGroup: &api.PodGroup{
							PodGroup: scheduling.PodGroup{
								ObjectMeta: v1.ObjectMeta{
									Labels: map[string]string{
										util.MatrixSchedulePolicyKey: util.BinPack,
										util.BinPackWeightCPUKey:     "0",
										util.BinPackWeightMemoryKey:  "1",
									},
								},
							},
						},
					},
				},
			},
			task: &api.TaskInfo{Job: "job1"},
		},
		want: GetMatrixSchedulePolicyArgWant{
			policy:  "",
			weight:  &util.PriorityWeight{},
			wantErr: true,
		},
	}
}

func getMatrixSchedulePolicyArgTestCase5() GetMatrixSchedulePolicyArgTest {
	return GetMatrixSchedulePolicyArgTest{
		name: "GetMatrixSchedulePolicyArg fail test when BinPackWeightMemory is 0",
		tp:   &MatrixVolcanoPlugin{},
		args: GetMatrixSchedulePolicyArgArgs{
			ssn: &framework.Session{
				Jobs: map[api.JobID]*api.JobInfo{
					"job1": {
						PodGroup: &api.PodGroup{
							PodGroup: scheduling.PodGroup{
								ObjectMeta: v1.ObjectMeta{
									Labels: map[string]string{
										util.MatrixSchedulePolicyKey: util.BinPack,
										util.BinPackWeightCPUKey:     "1",
										util.BinPackWeightMemoryKey:  "0",
									},
								},
							},
						},
					},
				},
			},
			task: &api.TaskInfo{Job: "job1"},
		},
		want: GetMatrixSchedulePolicyArgWant{
			policy:  "",
			weight:  &util.PriorityWeight{},
			wantErr: true,
		},
	}
}

func buildGetMatrixSchedulePolicyArgTestCases() []GetMatrixSchedulePolicyArgTest {
	GetMatrixSchedulePolicyArgTests := []GetMatrixSchedulePolicyArgTest{
		getMatrixSchedulePolicyArgTestCase1(),
		getMatrixSchedulePolicyArgTestCase2(),
		getMatrixSchedulePolicyArgTestCase3(),
		getMatrixSchedulePolicyArgTestCase4(),
		getMatrixSchedulePolicyArgTestCase5(),
	}
	return GetMatrixSchedulePolicyArgTests
}

func TestGetMatrixSchedulePolicyArg(t *testing.T) {
	buildTests := buildGetMatrixSchedulePolicyArgTestCases()
	for _, tt := range buildTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.tp.GetMatrixSchedulePolicyArg(tt.args.ssn, tt.args.task); (res != nil) != tt.want.wantErr {
				t.Errorf("GetMatrixSchedulePolicyArg() res = %v, wantErr %v", res, tt.want.wantErr)
			}
			if tt.tp.Policy != tt.want.policy {
				t.Errorf("GetMatrixSchedulePolicyArg() Policy = %v, want = %v", tt.tp.Policy, tt.want.policy)
			}
			if tt.tp.Weight.BinPackCPU != tt.want.weight.BinPackCPU ||
				tt.tp.Weight.BinPackMemory != tt.want.weight.BinPackMemory {
				t.Errorf("GetMatrixSchedulePolicyArg() Weight CPU = %v Mem = %v, want CPU = %v Mem = %v",
					tt.tp.Weight.BinPackCPU, tt.tp.Weight.BinPackMemory,
					tt.want.weight.BinPackCPU, tt.want.weight.BinPackMemory)
			}
		})
	}
}

type GetMatrixInfoTest struct {
	name string
	tp   *MatrixVolcanoPlugin
	want MatrixVolcanoPlugin
}

func buildGetMatrixInfoTestCases() []GetMatrixInfoTest {
	GetMatrixInfoTests := []GetMatrixInfoTest{
		{
			name: "GetMatrixInfo success test",
			tp: &MatrixVolcanoPlugin{
				Policy: util.BinPack,
				Weight: util.PriorityWeight{
					BinPackCPU:    1,
					BinPackMemory: 1,
				},
				CurrJobId:           "job1",
				RackSelectionStatus: true,
				ChosenRackId:        "rack1",
				ChosenShareDomain:   "share1",
			},
			want: MatrixVolcanoPlugin{
				Policy: util.BinPack,
				Weight: util.PriorityWeight{
					BinPackCPU:    1,
					BinPackMemory: 1,
				},
				CurrJobId:           "job1",
				RackSelectionStatus: true,
				ChosenRackId:        "rack1",
				ChosenShareDomain:   "share1",
			},
		},
	}
	return GetMatrixInfoTests
}

func TestGetMatrixInfo(t *testing.T) {
	buildTests := buildGetMatrixInfoTestCases()
	for _, tt := range buildTests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.tp.GetMatrixInfo()
			if res.Policy != tt.want.Policy {
				t.Errorf("GetMatrixInfo() policy = %v want policy = %v", res.Policy, tt.want.Policy)
			}
			if res.Weight.BinPackCPU != tt.want.Weight.BinPackCPU {
				t.Errorf("GetMatrixInfo() cpu weight = %v want cpu weight = %v",
					res.Weight.BinPackCPU, tt.want.Weight.BinPackCPU)
			}
			if res.Weight.BinPackMemory != tt.want.Weight.BinPackMemory {
				t.Errorf("GetMatrixInfo() mem weight = %v want mem weight = %v",
					res.Weight.BinPackMemory, tt.want.Weight.BinPackMemory)
			}
			if res.CurrJobId != tt.want.CurrJobId {
				t.Errorf("GetMatrixInfo() CurrJobId = %v want CurrJobId = %v",
					res.CurrJobId, tt.want.CurrJobId)
			}
			if res.RackSelectionStatus != tt.want.RackSelectionStatus {
				t.Errorf("GetMatrixInfo() RackSelectionStatus = %v want RackSelectionStatus = %v",
					res.RackSelectionStatus, tt.want.RackSelectionStatus)
			}
			if res.ChosenRackId != tt.want.ChosenRackId {
				t.Errorf("GetMatrixInfo() ChosenRackId = %v want ChosenRackId = %v",
					res.ChosenRackId, tt.want.ChosenRackId)
			}
			if res.ChosenShareDomain != tt.want.ChosenShareDomain {
				t.Errorf("GetMatrixInfo() ChosenShareDomain = %v want ChosenShareDomain = %v",
					res.ChosenShareDomain, tt.want.ChosenShareDomain)
			}
		})
	}
}
