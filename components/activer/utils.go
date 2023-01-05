package main

import (
	"encoding/json"
	"fmt"
	"github.com/thedevsaddam/gojsonq/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/utils"
)

func GetPartOfJsonResponce(reg string, response interface{}) string {
	//只支持如下两种
	if reg == "$" {
		bytes, _ := json.Marshal(response)
		return string(bytes)
	}
	//"$.level1.level2...."
	part := string(([]byte(reg))[2:])
	bytes, _ := json.Marshal(gojsonq.New().FromInterface(response).Find(part))
	return string(bytes)
}

func AdjustReplicaNum2Target(etcdClient *clientv3.Client, funcName string, target int) {
	function := utils.GetFunctionByName(etcdClient, funcName)
	replicaNameList := utils.GetPodReplicaIDListByPodName(etcdClient, function.PodName)
	fmt.Println("target size is:   ", target)
	fmt.Println("len(replicaNameList) is:   ", len(replicaNameList))
	if len(replicaNameList) < target {
		utils.AddNPodInstance(function.PodName, target-len(replicaNameList))
	} else if len(replicaNameList) > target {
		utils.RemovePodInstance(function.PodName, len(replicaNameList)-target)
		//if target == 0 {
		//	StopService(function.ServiceName)
		//}
	}
}
