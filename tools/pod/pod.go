package pod

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"log"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/utils"
	"strings"
	"time"
)

func CreateAndStartPod(podInstance *def.PodInstance, node *def.Node) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	defer func(cli *client.Client) {
		_ = cli.Close()
	}(cli)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("main here")
	containers := podInstance.Spec.Containers
	containerIDs := make([]string, 0)

	// Create the NetBridge if necessary
	networkID := docker.CreateNetBridge(node.CniIP.String())

	// Create the Pause container
	pauseContainerID := docker.CreatePauseContainer(cli, containers, podInstance.ID, networkID)
	pauserDetail, _ := docker.InspectContainer(pauseContainerID)
	podInstance.IP = pauserDetail.NetworkSettings.Networks["miniK8S-bridge"].IPAddress

	podInstance.PauseContainer = def.ContainerStatus{Status: def.RUNNING, ID: pauseContainerID}

	fmt.Println(pauserDetail.NetworkSettings)
	fmt.Println(pauserDetail.NetworkSettings.Networks["miniK8S-bridge"])
	fmt.Printf("podInstance.ClusterIP is %s\n", podInstance.IP)

	fmt.Println("podInstance.Spec.Volumes:  ", podInstance.Spec.Volumes)

	for index, con := range containers {
		config := docker.GenerateConfig(con)

		containerMode := "container:" + pauseContainerID
		hostConfig := docker.GenerateHostConfig(con, containerMode, podInstance.Spec.Volumes)

		tmpCons := make([]def.Container, 0)
		tmpCons = append(tmpCons, con)
		//exportsPort, _ := generatePort(con)
		//fmt.Println(exportsPort)
		//config.ExposedPorts = exportsPort

		networkingConfig := docker.GenerateNetworkingConfig(networkID)

		docker.ImageEnsure(con.Image)
		prefix := podInstance.ID[1:]
		prefix = strings.Replace(prefix, "/", "-", -1)
		name := prefix + `-` + con.Name
		body, err := cli.ContainerCreate(
			context.Background(),
			config, hostConfig,
			networkingConfig,
			nil,
			name)
		if err != nil {
			//if error, stop all containers has been created
			podInstance.Status = def.FAILED
			utils.PersistPodInstance(*podInstance, node.EtcdClient)
			for _, id := range containerIDs {
				docker.StopContainer(id)
				_, _ = docker.RemoveContainer(id)
			}
			log.Fatal(err)
			return
		}
		fmt.Println("created " + body.ID)
		containerIDs = append(containerIDs, body.ID)
		docker.StartContainer(body.ID)
		podInstance.ContainerSpec[index].Status = def.RUNNING
		podInstance.ContainerSpec[index].ID = body.ID
		utils.PersistPodInstance(*podInstance, node.EtcdClient)
	}
	podInstance.Status = def.RUNNING
	podInstance.StartTime = time.Now()
	utils.PersistPodInstance(*podInstance, node.EtcdClient)
	/* 暂时不使用 */
	//go podInstance.PodDaemon()
}

func StopAndRemovePod(podInstance *def.PodInstance, node *def.Node) {
	if podInstance.Status == def.RUNNING {
		podInstance.Status = def.SUCCEEDED
	}
	utils.PersistPodInstance(*podInstance, node.EtcdClient)
	for index, container := range podInstance.ContainerSpec {
		if podInstance.ContainerSpec[index].Status == def.RUNNING {
			podInstance.ContainerSpec[index].Status = def.SUCCEEDED
		}
		docker.StopContainer(container.ID)
		_, _ = docker.RemoveContainer(container.ID)
		utils.PersistPodInstance(*podInstance, node.EtcdClient)
	}
	docker.StopContainer(podInstance.PauseContainer.ID)
	_, _ = docker.RemoveContainer(podInstance.PauseContainer.ID)

	if podInstance.PauseContainer.Status == def.RUNNING {
		podInstance.PauseContainer.Status = def.SUCCEEDED
	}
	utils.PersistPodInstance(*podInstance, node.EtcdClient)
}
