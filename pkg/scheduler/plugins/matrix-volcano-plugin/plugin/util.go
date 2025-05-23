/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

// Package plugin is used for matrix-volcano-plugin
package plugin

import "sync"

var (
	// nodeRankMapCache 节点到Rank的映射: nodeName -> rankID
	nodeRankMapCache = make(map[string]string)
	// rankNodesMapCache Rank到节点列表的映射: rankID -> []nodeName
	rankNodesMapCache = make(map[string][]string)
	// nodeShareDomainMapCache 节点到share domain的映射: nodeName -> []shareDomain
	nodeShareDomainMapCache = make(map[string][][]string)
	// cacheLock 数据读写锁
	cacheLock *sync.RWMutex
)
