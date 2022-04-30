package util

import (
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/pod"
)

func PersistPodInstance(podInstance pod.PodInstance, cli *clientv3.Client) {
	byts, _ := json.Marshal(podInstance)
	etcd.Put(cli, podInstance.ID, string(byts))
}
