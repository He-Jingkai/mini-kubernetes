package get_api

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"strconv"
	"time"
)

func GetAllPodInstanceStatus(cli *clientv3.Client) ([]def.PodInstanceBrief, bool) {
	podInstanceList, flag := GetAllPodInstance(cli)
	resultList := make([]def.PodInstanceBrief, 0)
	for _, podInstance := range podInstanceList {
		brief := def.PodInstanceBrief{
			Name:     podInstance.Metadata.Name,
			Status:   podInstance.Status,
			Restarts: podInstance.RestartCount,
		}
		containers := podInstance.ContainerSpec
		count := 0
		for _, container := range containers {
			if container.Status == def.RUNNING {
				count++
			} else {
				fmt.Println("wadadawdawdwa")
				fmt.Println(container.Status)
			}
		}
		brief.Ready = strconv.Itoa(count) + "/" + strconv.Itoa(len(containers))
		t := time.Now()
		brief.Age = t.Sub(podInstance.StartTime)

		resultList = append(resultList, brief)
	}

	return resultList, flag
}