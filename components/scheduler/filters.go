package main

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"math"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/utils"
)

func GetNodeResourceInfo(etcdClient *clientv3.Client, nodeID int) *def.NodeResourceSchedulerCache {
	nodeResourceSchedulerCache := def.NodeResourceSchedulerCache{
		Validate: false,
	}
	nodeResource := utils.GetResourceUsageSequenceByNodeID(etcdClient, nodeID)
	//length := len(nodeResource.Sequence)
	//if length != 0 {
	//	latestRecode := nodeResource.Sequence[length-1]
	if nodeResource.Valid {
		latestRecode := nodeResource
		nodeResourceSchedulerCache = def.NodeResourceSchedulerCache{
			CPULoad:     float64(latestRecode.CPULoad) / 1000,
			CPUNum:      latestRecode.CPUNum,
			MemoryUsage: latestRecode.MemoryUsage,
			MemoryTotal: latestRecode.MemoryUsage,
			Validate:    true,
		}
	}
	return &nodeResourceSchedulerCache
}

func GetInfoOfANode(etcdClient *clientv3.Client, nodeID int) (*def.NodeInfoSchedulerCache, []string) {
	nodeInfo := def.NodeInfoSchedulerCache{
		NodeID:          nodeID,
		PodInstanceList: []def.PodInstanceSchedulerCache{},
	}
	replicaNameList := utils.GetAllPodInstancesOfANode(nodeID, etcdClient)
	// for every replica, generate a PodInstanceSchedulerCache
	for _, replicaName := range replicaNameList {
		podInstance := utils.GetPodInstanceByName(etcdClient, replicaName)
		if podInstance.Status != def.FAILED {
			replicaInfo := def.PodInstanceSchedulerCache{
				InstanceName: replicaName,
				PodName:      podInstance.Pod.Metadata.Name,
			}
			nodeInfo.PodInstanceList = append(nodeInfo.PodInstanceList, replicaInfo)
		}
	}
	return &nodeInfo, replicaNameList
}

func NotWithFilter(nodes []int, notWith string, allNodesInfo []*def.NodeInfoSchedulerCache) []int {
	fmt.Println("[scheduler notWithFilter]", notWith)
	var afterFilterNodes []int
	if notWith == "" {
		return nodes
	}
	for _, node := range nodes {
		fmt.Println("[node info]", node, *allNodesInfo[node])
		in := false
		for _, instance := range allNodesInfo[node].PodInstanceList {
			fmt.Println("[podInstance on node]", node, instance.PodName)
			if instance.PodName == notWith {
				in = true
				break
			}
		}
		if !in {
			afterFilterNodes = append(afterFilterNodes, node)
		}
	}
	fmt.Println("[notwith filter]result is ", afterFilterNodes)
	return afterFilterNodes
}

func WithFilter(nodes []int, with string, allNodesInfo []*def.NodeInfoSchedulerCache) []int {
	var afterFilterNodes []int
	if with == "" {
		return nodes
	}
	for _, node := range nodes {
		in := false
		for _, instance := range allNodesInfo[node].PodInstanceList {
			if instance.PodName == with {
				in = true
				break
			}
		}
		if in {
			afterFilterNodes = append(afterFilterNodes, node)
		}
	}
	return afterFilterNodes
}

func ResourceFilter(etcdClient *clientv3.Client, nodes []int, CPU int, memory int, allNodesInfo []*def.NodeInfoSchedulerCache) []int {
	var afterFilterNodes []int
	for _, node := range nodes {
		info := GetNodeResourceInfo(etcdClient, allNodesInfo[node].NodeID)
		if info.Validate == false ||
			(info.CPUNum >= CPU && (info.MemoryTotal-info.MemoryUsage) >= uint64(memory)) {
			afterFilterNodes = append(afterFilterNodes, node)
		}
	}
	return afterFilterNodes
}

func ChooseNode(etcdClient *clientv3.Client, nodes []int, allNodesInfo []*def.NodeInfoSchedulerCache) int {
	var chose int
	minInstances := math.MaxInt
	for _, node := range nodes {
		// TODO: 测试正确性
		if utils.GetNodeByID(etcdClient, allNodesInfo[node].NodeID).Status == def.NotReady {
			continue
		}
		instances := len(utils.GetPodInstanceIDListOfNode(etcdClient, allNodesInfo[node].NodeID))
		if instances < minInstances {
			minInstances = instances
			chose = node
		}
	}
	return chose
}

// TODO: CPU的单位问题

func MemoryToByte(memString string) int {
	if memString == `0` || memString == `` {
		return 0
	}
	memByte := 0
	for _, c := range memString {
		if c >= '0' && c <= '9' {
			memByte = memByte*10 + int(c-'0')
		} else if c == 'K' || c == 'k' {
			return memByte * 1024
		} else if c == 'M' || c == 'm' {
			return memByte * 1024 * 1024
		} else if c == 'G' || c == 'g' {
			return memByte * 1024 * 1024 * 1024
		}
	}
	return 0
}

func PodResourceRequest(podInstance *def.PodInstance) (int, int) {
	CPU := 0
	memory := 0
	for _, container := range podInstance.Pod.Spec.Containers {
		//requestCPU := container.Resources.ResourceRequest.CPU
		//if requestCPU != 0 && requestCPU > CPU {
		//	CPU = requestCPU
		//}
		memory += MemoryToByte(container.Resources.ResourceRequest.Memory)
	}
	return CPU, memory
}
