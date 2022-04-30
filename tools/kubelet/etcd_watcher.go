package kubelet

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/pod"
	"sort"
)

type PodInstances struct {
	Instances []string
}

func EtcdWatcher(e *echo.Echo) {
	key := fmt.Sprintf("Node%d_podInstances", node.NodeID)
	watch := etcd.Watch(node.EtcdClient, key)
	for wc := range watch {
		for _, w := range wc.Events {
			var instances []string
			_ = json.Unmarshal(w.Kv.Value, &instances)
			handlePodInstancesChange(instances, e)
		}
	}
}

func handlePodInstancesChange(instances []string, e *echo.Echo) {
	adds, deletes := comparePodList(instances)
	for _, add := range adds {
		resp := etcd.Get(node.EtcdClient, add)
		podInstance := pod.PodInstance{}
		jsonString := ``
		for _, ev := range resp.Kvs {
			jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
		}
		jsonString = fmt.Sprintf(`{%s}`, jsonString)
		_ = json.Unmarshal([]byte(jsonString), &podInstance)
		podInstance.CreateAndStartPod()
		node.PodInstances = append(node.PodInstances, &podInstance)
	}
	for _, index := range deletes {
		node.PodInstances[index].StopAndRemovePod()
	}
}

func comparePodList(instancesNew []string) (added []string, deleted []int) {
	var instancesCurrent []string
	sort.Strings(instancesNew)
	for index, instance := range node.PodInstances {
		if sort.SearchStrings(instancesNew, instance.Name) == len(instancesNew) {
			deleted = append(deleted, index)
		}
		instancesCurrent = append(instancesCurrent, instance.Name)
	}
	sort.Strings(instancesCurrent)
	for _, instance := range instancesNew {
		if sort.SearchStrings(instancesCurrent, instance) == len(instancesCurrent) {
			added = append(added, instance)
		}
	}
	return
}