package gateway

import (
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
)

//func WriteCmdToFile(filePath string, cmd string) {
//	file, err := os.Create(filePath)
//
//	//err := os.Truncate(filePath, 0)
//	if err != nil {
//		fmt.Println(err)
//	}
//	//file, _ := os.OpenFile(filePath, os.O_RDWR, os.ModeAppend)
//	_, err = file.Write([]byte(cmd))
//	if err != nil {
//		fmt.Println(err)
//	}
//	err = file.Close()
//	if err != nil {
//		fmt.Println(err)
//	}
//}

func MakeGatewayImage(dns *def.DNSDetail, nameGatewayImageName string) {
	container := def.Container{
		Image: def.GatewayImage,
	}
	containerID := docker.CreateContainer(container, nameGatewayImageName)
	fileStr := GenerateApplicationYaml(*dns)
	docker.CopyToContainer(containerID, def.RequirementsParentDirPath, def.GatewayRoutesConfigPathInImage, fileStr)
	//cmd := exec.Command("docker", "exec", containerID, "/bin/bash", "-c", fmt.Sprintf("'%s'", def.GatewayPackageCmd)).String()
	//WriteCmdToFile(def.TemplateCmdFilePath, cmd)
	//command := fmt.Sprintf(`%s .`, def.TemplateCmdFilePath)
	//err := exec.Command("/bin/bash", "-c", command).Run()
	//if err != nil {
	//	fmt.Println(err)
	//}
	docker.CommitContainer(containerID, nameGatewayImageName)
	docker.PushImage(nameGatewayImageName)
	docker.StopContainer(containerID)
	_, _ = docker.RemoveContainer(containerID)
}
