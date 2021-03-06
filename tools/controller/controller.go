package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron"
	"math"
	"mini-kubernetes/tools/controller/controller_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
	"os"
	"time"
)

var controllerMeta = def.ControllerMeta{
	ParsedDeployments:                []*def.ParsedDeployment{},
	DeploymentNameList:               []string{},
	ParsedHorizontalPodAutoscalers:   []*def.ParsedHorizontalPodAutoscaler{},
	HorizontalPodAutoscalersNameList: []string{},
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	etcdClient, err := etcd.Start("", def.EtcdPort)
	controllerMeta.EtcdClient = etcdClient
	if err != nil {
		e.Logger.Error("Start etcd error!")
		os.Exit(0)
	}
	ControllerMetaInit()
	go EtcdDeploymentWatcher()
	go EtcdHorizontalPodAutoscalerWatcher()

	go ReplicaChecker()
	go HorizontalPodAutoscalerChecker()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", def.ControllerPort)))
}

func ControllerMetaInit() {
	deploymentList := controller_utils.GetDeploymentNameList(controllerMeta.EtcdClient)
	horizontalPodAutoscalerNameList := controller_utils.GetHorizontalPodAutoscalerNameList(controllerMeta.EtcdClient)
	for _, name := range deploymentList {
		controllerMeta.ParsedDeployments = append(controllerMeta.ParsedDeployments, controller_utils.GetDeploymentByName(controllerMeta.EtcdClient, name))
	}
	for _, name := range horizontalPodAutoscalerNameList {
		controllerMeta.ParsedHorizontalPodAutoscalers = append(controllerMeta.ParsedHorizontalPodAutoscalers, controller_utils.GetHorizontalPodAutoscalerByName(controllerMeta.EtcdClient, name))
	}
	controllerMeta.DeploymentNameList = deploymentList
	controllerMeta.HorizontalPodAutoscalersNameList = horizontalPodAutoscalerNameList
}

func HandleDeploymentListChange(deploymentList []string) {
	controllerMeta.DeploymentLock.Lock()
	defer controllerMeta.DeploymentLock.Unlock()

	added, deleted := util.DifferTwoStringList(controllerMeta.DeploymentNameList, deploymentList)
	fmt.Println("added:   ", added)
	fmt.Println("deleted:   ", deleted)
	for _, name := range added {
		fmt.Println("name:  ", name)
		fmt.Println("def.GetKeyOfDeployment(deploymentName):  ", def.GetKeyOfDeployment(name))
		deployment := controller_utils.GetDeploymentByName(controllerMeta.EtcdClient, name)
		fmt.Println(deployment)
		fmt.Println("deployment.Name:  ", deployment.Name)
		controllerMeta.ParsedDeployments = append(controllerMeta.ParsedDeployments, deployment)
		controller_utils.NewReplicaNameListByPodName(controllerMeta.EtcdClient, deployment.PodName)
		//controller_utils.NewNPodInstance(controllerMeta.EtcdClient, deployment.PodName, deployment.ReplicasNum)
		util.AddNPodInstance(deployment.PodName, deployment.ReplicasNum)
	}
	for _, name := range deleted {
		DeleteADeployment(name)
	}
	controllerMeta.DeploymentNameList = deploymentList
}

func HandleHorizontalPodAutoscalerListChange(horizontalPodAutoscalerList []string) {
	controllerMeta.HorizontalPodAutoscalersLock.Lock()
	defer controllerMeta.HorizontalPodAutoscalersLock.Unlock()

	added, deleted := util.DifferTwoStringList(controllerMeta.HorizontalPodAutoscalersNameList, horizontalPodAutoscalerList)
	fmt.Println("added:  ", added)
	fmt.Println("deleted:  ", deleted)
	for _, name := range added {
		horizontalPodAutoscaler := controller_utils.GetHorizontalPodAutoscalerByName(controllerMeta.EtcdClient, name)
		controllerMeta.ParsedHorizontalPodAutoscalers = append(controllerMeta.ParsedHorizontalPodAutoscalers, horizontalPodAutoscaler)
		controller_utils.NewReplicaNameListByPodName(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName)
		//controller_utils.NewNPodInstance(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName, horizontalPodAutoscaler.MinReplicas)
	}
	for _, name := range deleted {
		DeleteAHorizontalPodAutoscaler(name)
	}
	controllerMeta.HorizontalPodAutoscalersNameList = horizontalPodAutoscalerList
}

func DeleteAHorizontalPodAutoscaler(name string) {
	controller_utils.RemoveAllReplicasOfPod(controllerMeta.EtcdClient, def.GetPodNameOfAutoscaler(name))
	// sync cache
	for index, horizontalPodAutoscaler := range controllerMeta.ParsedHorizontalPodAutoscalers {
		if horizontalPodAutoscaler.Name == name {
			controllerMeta.ParsedHorizontalPodAutoscalers = append(controllerMeta.ParsedHorizontalPodAutoscalers[:index], controllerMeta.ParsedHorizontalPodAutoscalers[index+1:]...)
			break
		}
	}
}

func DeleteADeployment(name string) {
	controller_utils.RemoveAllReplicasOfPod(controllerMeta.EtcdClient, def.GetPodNameOfDeployment(name))
	// sync cache
	for index, deployment := range controllerMeta.ParsedDeployments {
		if deployment.Name == name {
			controllerMeta.ParsedDeployments = append(controllerMeta.ParsedDeployments[:index], controllerMeta.ParsedDeployments[index+1:]...)
			break
		}
	}
}

func EtcdHorizontalPodAutoscalerWatcher() {
	watch := etcd.Watch(controllerMeta.EtcdClient, def.HorizontalPodAutoscalerListName)
	for wc := range watch {
		for _, w := range wc.Events {
			var nameList []string
			_ = json.Unmarshal(w.Kv.Value, &nameList)
			HandleHorizontalPodAutoscalerListChange(nameList)
		}
	}
}

func EtcdDeploymentWatcher() {
	watch := etcd.Watch(controllerMeta.EtcdClient, def.DeploymentListName)
	for wc := range watch {
		for _, w := range wc.Events {
			var nameList []string
			_ = json.Unmarshal(w.Kv.Value, &nameList)
			HandleDeploymentListChange(nameList)
		}
	}
}

func ReplicaChecker() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/5 * * * * *", CheckAllReplicas)
	if err != nil {
		fmt.Println(err)
	}
	cron2.Start()
	defer cron2.Stop()
	for {
		if controllerMeta.ShouldStop {
			break
		}
	}
}

func HorizontalPodAutoscalerChecker() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/15 * * * * *", CheckAllHorizontalPodAutoscalers)
	if err != nil {
		fmt.Println(err)
	}
	cron2.Start()
	defer cron2.Stop()
	for {
		if controllerMeta.ShouldStop {
			break
		}
	}
}

func CheckAllReplicas() {
	controllerMeta.DeploymentLock.Lock()
	defer controllerMeta.DeploymentLock.Unlock()

	for _, deployment := range controllerMeta.ParsedDeployments {
		pod := controller_utils.GetPodByName(controllerMeta.EtcdClient, deployment.PodName)
		instancelist := controller_utils.GetReplicaNameListByPodName(controllerMeta.EtcdClient, pod.Metadata.Name)
		health := 0
		for _, instanceID := range instancelist {
			podInstance := util.GetPodInstance(instanceID, controllerMeta.EtcdClient)
			if podInstance.Status != def.FAILED {
				health++
			} else {
				//controller_utils.RemovePodInstance(controllerMeta.EtcdClient, &podInstance)
				util.RemovePodInstanceByID(podInstance.ID)
			}
		}
		fmt.Printf("[controller replica checker]%s has %d health, %d health", deployment.PodName, len(instancelist), health)
		if health < deployment.ReplicasNum {
			//controller_utils.NewNPodInstance(controllerMeta.EtcdClient, pod.Metadata.Name, deployment.ReplicasNum-health)
			util.AddNPodInstance(pod.Metadata.Name, deployment.ReplicasNum-health)
		}
	}
}

func CheckAllHorizontalPodAutoscalers() {
	controllerMeta.HorizontalPodAutoscalersLock.Lock()
	defer controllerMeta.HorizontalPodAutoscalersLock.Unlock()

	for _, horizontalPodAutoscaler := range controllerMeta.ParsedHorizontalPodAutoscalers {
		pod := controller_utils.GetPodByName(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName)
		instancelist := controller_utils.GetReplicaNameListByPodName(controllerMeta.EtcdClient, pod.Metadata.Name)
		cpu := float64(0)
		memory := int64(0)
		minCPUUsagePodInstanceID := ""
		minCPUUsage := math.MaxFloat64
		minMemoryUsagePodInstanceID := ""
		minMemoryUsage := int64(math.MaxInt64)
		activeNum := 0
		tooShort := false
		fmt.Println("[controller instancelist] ", instancelist)
		for _, instanceID := range instancelist {
			podInstance := util.GetPodInstance(instanceID, controllerMeta.EtcdClient)
			if podInstance.Status == def.FAILED {
				fmt.Println(podInstance)
				fmt.Println(podInstance.ID, "is failed")
				util.RemovePodInstanceByID(podInstance.ID)
				continue
			} else if util.TimeToSecond(time.Now())-util.TimeToSecond(podInstance.StartTime) < 15 {
				tooShort = true
				break
			} else {
				activeNum++
				fmt.Println(podInstance.ID, "is health")
				podInstanceResourceUsage := controller_utils.GetPodInstanceResourceUsageByName(controllerMeta.EtcdClient, instanceID)
				//if podInstanceResourceUsage.Valid {
				instanceCPUUsage := float64(podInstanceResourceUsage.CPULoad) / 1000
				cpu += instanceCPUUsage
				if instanceCPUUsage < minCPUUsage {
					minCPUUsage = instanceCPUUsage
					minCPUUsagePodInstanceID = podInstance.ID
				}

				instanceMemoryUsage := int64(podInstanceResourceUsage.MemoryUsage)
				memory += instanceMemoryUsage
				if instanceMemoryUsage < minMemoryUsage {
					minMemoryUsage = instanceMemoryUsage
					minMemoryUsagePodInstanceID = podInstance.ID
				}
			}
			//}
		}
		if tooShort {
			continue
		}
		fmt.Println("activeNum is ", activeNum, " cpu is ", cpu, " memory is ", memory)
		//calculate avg

		if activeNum < horizontalPodAutoscaler.MinReplicas {
			//controller_utils.NewNPodInstance(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName, horizontalPodAutoscaler.MinReplicas-activeNum)
			util.AddNPodInstance(horizontalPodAutoscaler.PodName, 1)
		} else {
			if activeNum == 0 {
				return
			}
			cpuAvg := cpu / float64(activeNum)
			memAvg := float64(memory) / float64(activeNum)
			if cpuAvg < 0.8*controller_utils.CPUToMCore(horizontalPodAutoscaler.CPUMinValue) {
				//CPU???????????????, ????????????
				fmt.Println("cpu avg too small")
				if activeNum > horizontalPodAutoscaler.MinReplicas {
					//controller_utils.RemovePodInstance(controllerMeta.EtcdClient, &minCPUUsagePodInstance)
					fmt.Println("cpu avg too small, shrink")
					util.RemovePodInstanceByID(minCPUUsagePodInstanceID)
				}
			} else if memAvg < 0.8*float64(controller_utils.MemoryToByte(horizontalPodAutoscaler.MemoryMinValue)) {
				//mem???????????????, ????????????
				fmt.Println("mem avg too small")
				if activeNum > horizontalPodAutoscaler.MinReplicas {
					//controller_utils.RemovePodInstance(controllerMeta.EtcdClient, &minMemoryUsagePodInstance)
					fmt.Println("mem avg too small, shrink")
					util.RemovePodInstanceByID(minMemoryUsagePodInstanceID)
				}
			} else if cpuAvg > 1.2*controller_utils.CPUToMCore(horizontalPodAutoscaler.CPUMaxValue) {
				//CPU???????????????, ????????????
				fmt.Println("cpu avg too large")
				if activeNum < horizontalPodAutoscaler.MaxReplicas {
					//controller_utils.NewNPodInstance(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName, 1)
					fmt.Println("cpu avg too large, expand")
					util.AddNPodInstance(horizontalPodAutoscaler.PodName, 1)
				}
			} else if memAvg > 1.2*float64(controller_utils.MemoryToByte(horizontalPodAutoscaler.MemoryMaxValue)) {
				//memory???????????????, ????????????
				fmt.Println("mem avg too large")
				if activeNum < horizontalPodAutoscaler.MaxReplicas {
					//controller_utils.NewNPodInstance(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName, 1)
					fmt.Println("mem avg too large, expand")
					util.AddNPodInstance(horizontalPodAutoscaler.PodName, 1)
				}
			}
		}

		//	if activeNum < horizontalPodAutoscaler.MinReplicas {
		//		controller_utils.NewNPodInstance(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName, horizontalPodAutoscaler.MinReplicas-activeNum)
		//	} else if cpu < 0.8*controller_utils.CPUToMCore(horizontalPodAutoscaler.CPUMinValue) || float64(memory) < 0.8*float64(controller_utils.MemoryToByte(horizontalPodAutoscaler.MemoryMinValue)) {
		//		if activeNum < horizontalPodAutoscaler.MaxReplicas {
		//			controller_utils.NewNPodInstance(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName, 1)
		//		}
		//	} else if cpu > 1.2*controller_utils.CPUToMCore(horizontalPodAutoscaler.CPUMaxValue) {
		//		if activeNum > horizontalPodAutoscaler.MinReplicas {
		//			controller_utils.RemovePodInstance(controllerMeta.EtcdClient, &minCPUUsagePodInstance)
		//		}
		//	} else if float64(memory) > 1.2*float64(controller_utils.MemoryToByte(horizontalPodAutoscaler.MemoryMaxValue)) {
		//		if activeNum > horizontalPodAutoscaler.MinReplicas {
		//			controller_utils.RemovePodInstance(controllerMeta.EtcdClient, &minMemoryUsagePodInstance)
		//		}
		//	}
	}
}
