package create_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/gpujob"
	"mini-kubernetes/tools/utils"
)

func CreateGPUJobUploader(cli *clientv3.Client, job def.GPUJob) {
	pod_ := gpujob.GenerateGpuJobUploaderPod(&job)
	utils.PersistGPUJob(cli, job)
	podInstance := def.PodInstance{}
	podInstance.Pod = pod_

	//将pod存入etcd中
	utils.PersistPod(cli, pod_)

	//将新创建的podInstance写入到etcd当中
	podInstanceKey := def.GenerateKeyOfPodInstanceReplicas(pod_.Metadata.Name)
	podInstance.ID = podInstanceKey
	podInstance.ContainerSpec = make([]def.ContainerStatus, len(pod_.Spec.Containers))

	utils.PersistPodInstance(podInstance, cli)
	replicaIDList := []string{podInstanceKey}
	value, err := json.Marshal(replicaIDList)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, def.GetKeyOfPodReplicasNameListByPodName(pod_.Metadata.Name), string(value))
	utils.AddPodInstanceIDToList(cli, podInstance.ID)
}
