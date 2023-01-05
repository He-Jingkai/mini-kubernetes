package function_api

import (
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/utils"
)

func CreateFuncPodInstance(cli *clientv3.Client, podName string) def.PodInstance {
	pod_ := utils.GetPodByPodName(cli, podName)
	podInstance := def.PodInstance{}
	podInstance.Pod = pod_

	//将新创建的podInstance写入到etcd当中
	podInstanceKey := def.GenerateKeyOfPodInstanceReplicas(pod_.Metadata.Name)
	podInstance.ID = podInstanceKey
	podInstance.ContainerSpec = make([]def.ContainerStatus, len(pod_.Spec.Containers))

	utils.PersistPodInstance(podInstance, cli)
	utils.AddPodInstanceIDToList(cli, podInstance.ID)

	//更新ReplicasNameList
	instanceIDListkey := def.GetKeyOfPodReplicasNameListByPodName(podName)
	var instanceIDList []string
	utils.EtcdUnmarshal(etcd.Get(cli, instanceIDListkey), &instanceIDList)
	instanceIDList = append(instanceIDList, podInstance.ID)
	instanceIDListValue, _ := json.Marshal(instanceIDList)
	etcd.Put(cli, instanceIDListkey, string(instanceIDListValue))

	return podInstance
}
