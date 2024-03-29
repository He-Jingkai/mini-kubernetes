package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/utils"
)

func GetAllPodInstance(cli *clientv3.Client) ([]def.PodInstance, bool) {
	flag := false
	podInstancePrefix := "/podInstance/"
	kvs := etcd.GetWithPrefix(cli, podInstancePrefix).Kvs
	podInstanceValue := make([]byte, 0)
	podInstanceList := make([]def.PodInstance, 0)
	if len(kvs) != 0 {
		for _, kv := range kvs {
			podInstance := def.PodInstance{}
			podInstanceValue = kv.Value
			err := json.Unmarshal(podInstanceValue, &podInstance)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Println("podInstance.Metadata.Name is " + podInstance.Metadata.Name)

			// add for heartbeat
			if utils.GetNodeByID(cli, podInstance.NodeID).Status == def.NotReady {
				podInstance.Status = def.UNKNOWN
			}

			podInstanceList = append(podInstanceList, podInstance)
		}
		flag = true
	}

	return podInstanceList, flag
}
