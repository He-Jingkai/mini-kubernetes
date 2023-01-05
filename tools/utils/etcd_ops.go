package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/httpget"
)

func GetPodInstance(podInstanceID string, cli *clientv3.Client) def.PodInstance {
	resp := etcd.Get(cli, podInstanceID)
	podInstance := def.PodInstance{}
	EtcdUnmarshal(resp, &podInstance)
	return podInstance
}

func EtcdUnmarshal(resp *clientv3.GetResponse, v interface{}) {
	kv := resp.Kvs
	value := make([]byte, 0)
	if len(kv) != 0 {
		value = kv[0].Value
		err := json.Unmarshal(value, v)
		if err != nil {
			panic(err)
		}
	}
}

func PersistPodInstance(podInstance def.PodInstance, cli *clientv3.Client) {
	byts, _ := json.Marshal(podInstance)
	etcd.Put(cli, podInstance.ID, string(byts))
}

func GetDeploymentNameList(etcdClient *clientv3.Client) []string {
	var deploymentNameList []string
	EtcdUnmarshal(etcd.Get(etcdClient, def.DeploymentListName), &deploymentNameList)
	return deploymentNameList
}

func GetHorizontalPodAutoscalerNameList(etcdClient *clientv3.Client) []string {
	var horizontalPodAutoscalerNameList []string
	EtcdUnmarshal(etcd.Get(etcdClient, def.HorizontalPodAutoscalerListName), &horizontalPodAutoscalerNameList)
	return horizontalPodAutoscalerNameList
}

func GetDeploymentByName(etcdClient *clientv3.Client, deploymentName string) *def.ParsedDeployment {
	key := def.GetKeyOfDeployment(deploymentName)
	deployment := def.ParsedDeployment{}
	EtcdUnmarshal(etcd.Get(etcdClient, key), &deployment)
	return &deployment
}

func GetHorizontalPodAutoscalerByName(etcdClient *clientv3.Client, horizontalPodAutoscalerName string) *def.ParsedHorizontalPodAutoscaler {
	key := def.GetKeyOfAutoscaler(horizontalPodAutoscalerName)
	horizontalPodAutoscaler := def.ParsedHorizontalPodAutoscaler{}
	EtcdUnmarshal(etcd.Get(etcdClient, key), &horizontalPodAutoscaler)
	return &horizontalPodAutoscaler
}

func RemoveAllReplicasOfPod(etcdClient *clientv3.Client, podName string) {
	// remove from instance list, scheduler will remove it from node
	fmt.Println("podName is:  ", podName)
	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
	podInstanceIDList := GetReplicaNameListByPodName(etcdClient, podName)
	fmt.Println("GetReplicaNameListByPodName: ", key)
	fmt.Println("podInstanceIDList: ", podInstanceIDList)
	for _, instanceID := range podInstanceIDList {
		fmt.Println("try to get by instanceID:  ", instanceID)
		instance := GetPodInstance(instanceID, etcdClient)
		//RemovePodInstance(etcdClient, &instance)
		RemovePodInstanceByID(instance.ID)
	}
	// remove it's pod-replica entry
	etcd.Delete(etcdClient, key)
}

func GetReplicaNameListByPodName(etcdClient *clientv3.Client, podName string) []string {
	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
	var deploymentNameList []string
	EtcdUnmarshal(etcd.Get(etcdClient, key), &deploymentNameList)
	return deploymentNameList
}

func NewReplicaNameListByPodName(etcdClient *clientv3.Client, podName string) {
	// add empty list
	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
	var deploymentNameList []string
	newJsonString, _ := json.Marshal(deploymentNameList)
	etcd.Put(etcdClient, key, string(newJsonString))
}

func GetPodInstanceResourceUsageByName(etcdClient *clientv3.Client, podInstanceID string) *def.ResourceUsage {
	key := def.GetKeyOfResourceUsageByPodInstanceID(podInstanceID)
	resourceUsage := def.ResourceUsage{}
	EtcdUnmarshal(etcd.Get(etcdClient, key), &resourceUsage)
	return &resourceUsage
}

func GetGPUJobByName(li *clientv3.Client, jobName string) def.GPUJob {
	gpuJob := def.GPUJob{}
	key := def.GetGPUJobKeyByName(jobName)
	EtcdUnmarshal(etcd.Get(li, key), &gpuJob)
	return gpuJob
}

func GetPodReplicaListByPodName(li *clientv3.Client, podName string) []string {
	var list []string
	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
	EtcdUnmarshal(etcd.Get(li, key), &list)
	return list
}

func GetPodInstanceByID(li *clientv3.Client, id string) def.PodInstance {
	instance := def.PodInstance{}
	EtcdUnmarshal(etcd.Get(li, id), &instance)
	return instance
}

func PersistPod(li *clientv3.Client, pod_ def.Pod) {
	key := def.GetKeyOfPod(pod_.Metadata.Name)
	value, err := json.Marshal(pod_)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func GetPodByPodName(li *clientv3.Client, podName string) def.Pod {
	pod_ := def.Pod{}
	key := def.GetKeyOfPod(podName)
	EtcdUnmarshal(etcd.Get(li, key), &pod_)
	return pod_
}

//func GetPodReplicaIDListByPodName(li *clientv3.Client, podName string) []string {
//	idList := make([]string, 0)
//	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
//	utils.EtcdUnmarshal(etcd.Get(li, key), &idList)
//	return idList
//}

func PersistStateMachine(li *clientv3.Client, stateMachine def.StateMachine) {
	key := def.GetKeyOfStateMachine(stateMachine.Name)
	value, err := json.Marshal(stateMachine)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func PersistService(li *clientv3.Client, service def.Service) {
	key := def.GetKeyOfService(service.Name)
	value, err := json.Marshal(service)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func PersistGPUJob(li *clientv3.Client, job def.GPUJob) {
	key := def.GetGPUJobKeyByName(job.Name)
	value, err := json.Marshal(job)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func PersistFunction(li *clientv3.Client, function def.Function) {
	key := def.GetKeyOfFunction(function.Name)
	value, err := json.Marshal(function)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func AddFunctionNameToList(li *clientv3.Client, functionName string) {
	key := def.FunctionNameListKey
	var list []string
	EtcdUnmarshal(etcd.Get(li, key), &list)
	list = append(list, functionName)
	value, err := json.Marshal(list)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func AddPodInstanceIDToList(li *clientv3.Client, id string) {
	key := def.PodInstanceListID
	var list []string
	EtcdUnmarshal(etcd.Get(li, key), &list)
	list = append(list, id)
	value, err := json.Marshal(list)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func GetNodeList(li *clientv3.Client) []int {
	var list []int
	EtcdUnmarshal(etcd.Get(li, def.NodeListID), &list)
	return list
}

func GetNodeByID(li *clientv3.Client, nodeID int) def.Node {
	node := def.Node{}
	key := def.GetKeyOgNodeByNodeID(nodeID)
	EtcdUnmarshal(etcd.Get(li, key), &node)
	return node
}

func PersistNode(li *clientv3.Client, node def.Node) {
	value, err := json.Marshal(node)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, def.GetKeyOgNodeByNodeID(node.NodeID), string(value))
}

func GetFunctionByName(li *clientv3.Client, functionName string) def.Function {
	key := def.GetKeyOfFunction(functionName)
	function := def.Function{}
	EtcdUnmarshal(etcd.Get(li, key), &function)
	return function
}

func GetPodByName(li *clientv3.Client, podName string) def.Pod {
	key := def.GetKeyOfPod(podName)
	pod_ := def.Pod{}
	EtcdUnmarshal(etcd.Get(li, key), &pod_)
	return pod_
}

func GetFunctionNameList(etcdClient *clientv3.Client) []string {
	var functionNameList []string
	EtcdUnmarshal(etcd.Get(etcdClient, def.FunctionNameListKey), &functionNameList)
	return functionNameList
}

func GetStateMachineByName(etcdClient *clientv3.Client, name string) *def.StateMachine {
	stateMachine := def.StateMachine{}
	EtcdUnmarshal(etcd.Get(etcdClient, def.GetKeyOfStateMachine(name)), &stateMachine)
	return &stateMachine
}

func GetPodReplicaIDListByPodName(etcdClient *clientv3.Client, podName string) []string {
	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
	var instanceIDList []string
	EtcdUnmarshal(etcd.Get(etcdClient, key), &instanceIDList)
	return instanceIDList
}

func GetServiceByName(etcdClient *clientv3.Client, serviceName string) *def.Service {
	key := def.GetKeyOfService(serviceName)
	service := def.Service{}
	EtcdUnmarshal(etcd.Get(etcdClient, key), &service)
	return &service
}

func GetAllPodInstancesOfANode(nodeID int, etcdClient *clientv3.Client) []string {
	var replicaNameList []string
	EtcdUnmarshal(etcd.Get(etcdClient, def.PodInstanceListKeyOfNodeID(nodeID)), &replicaNameList)
	return replicaNameList
}

func GetAllPodInstancesID(etcdClient *clientv3.Client) []string {
	var allReplicas []string
	EtcdUnmarshal(etcd.Get(etcdClient, def.PodInstanceListID), &allReplicas)
	return allReplicas
}

func GetPodInstanceByName(etcdClient *clientv3.Client, replicaName string) def.PodInstance {
	podInstance := def.PodInstance{}
	EtcdUnmarshal(etcd.Get(etcdClient, replicaName), &podInstance)
	return podInstance
}

func GetResourceUsageSequenceByNodeID(etcdClient *clientv3.Client, nodeID int) def.ResourceUsage {
	// TODO: 注册node时添加空ResourceUsage(valid = false)
	nodeResource := def.ResourceUsage{}
	EtcdUnmarshal(etcd.Get(etcdClient, def.KeyNodeResourceUsage(nodeID)), &nodeResource)
	return nodeResource
}

func GetAllNodesID(etcdClient *clientv3.Client) []int {
	var nodeIDList []int
	EtcdUnmarshal(etcd.Get(etcdClient, def.NodeListID), &nodeIDList)
	return nodeIDList
}

func DeletePodInstanceFromNode(etcdClient *clientv3.Client, nodeID int, instanceName string) {
	replicaNameList := GetAllPodInstancesOfANode(nodeID, etcdClient)
	for index, replicaName := range replicaNameList {
		if replicaName == instanceName {
			replicaNameList = append(replicaNameList[:index], replicaNameList[index+1:]...)
			break
		}
	}
	PersistPodInstanceListOfNode(etcdClient, replicaNameList, nodeID)
}
func AddPodInstanceToNode(etcdClient *clientv3.Client, nodeID int, instance *def.PodInstance) {
	instance.NodeID = nodeID
	PersistPodInstance(*instance, etcdClient)
	replicaNameList := GetAllPodInstancesOfANode(nodeID, etcdClient)
	replicaNameList = append(replicaNameList, instance.ID)
	PersistPodInstanceListOfNode(etcdClient, replicaNameList, nodeID)
}

func PersistPodInstanceListOfNode(etcdClient *clientv3.Client, replicaNameList []string, nodeID int) {
	newJsonString, _ := json.Marshal(replicaNameList)
	etcd.Put(etcdClient, def.PodInstanceListKeyOfNodeID(nodeID), string(newJsonString))
}

func GetPodInstanceIDListOfNode(etcdClient *clientv3.Client, nodeID int) []string {
	key := def.GetKeyOfPodInstanceListKeyOfNodeByID(nodeID)
	var replicaIDList []string
	EtcdUnmarshal(etcd.Get(etcdClient, key), &replicaIDList)
	return replicaIDList
}

func AddNPodInstance(podName string, num int) {
	//apiServer add a podInstance
	for i := 0; i < num; i++ {
		request2 := podName
		response2 := ""
		body2, _ := json.Marshal(request2)
		err, status := httpget.Post("http://" + GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/create/replicasPodInstance").
			ContentType("application/json").
			Body(bytes.NewReader(body2)).
			GetString(&response2).
			Execute()
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
		}
		fmt.Printf("create_funcPodInstance is %s and response is: %s\n", status, response2)
	}
}

func RemovePodInstance(podName string, num int) {
	//apiServer delete a podInstance
	for i := 0; i < num; i++ {
		response4 := ""
		err, status := httpget.DELETE("http://" + GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/delete/funcPodInstance/" + podName).
			ContentType("application/json").
			GetString(&response4).
			Execute()
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
		}

		fmt.Printf("delete funcPodInstance status is %s\n", status)
		if status == "200" {
			fmt.Printf("delete funcPodInstance %s successfully and the response is: %v\n", podName, response4)
		} else {
			fmt.Printf("funcPodInstance %s doesn't exist\n", podName)
		}
	}
}

func RemovePodInstanceByID(podInstanceID string) {
	//apiServer delete a podInstance
	response4 := ""
	err, status := httpget.DELETE("http://" + GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/delete/replicasPodInstance/" + podInstanceID).
		ContentType("application/json").
		GetString(&response4).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("delete funcPodInstance status is %s\n", status)
	if status == "200" {
		fmt.Printf("delete funcPodInstance %s successfully and the response is: %v\n", podInstanceID, response4)
	} else {
		fmt.Printf("funcPodInstance %s doesn't exist\n", podInstanceID)
	}
}
