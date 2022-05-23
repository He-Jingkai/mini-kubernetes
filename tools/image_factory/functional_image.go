package image_factory

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/util"
	"os/exec"
)

func MakeFunctionalImage(function *def.Function) {
	pyString := util.ReadFile(function.Function)
	requirementsString := util.ReadFile(function.Requirements)
	container := def.Container{
		Image: def.PyFunctionTemplateImage,
	}
	imageName := fmt.Sprintf("image_%s_%d", function.Name, function.Version)
	function.Image = imageName
	containerID := docker.CreateContainer(container, imageName)
	docker.CopyToContainer(containerID, def.PyHandlerParentDirPath, def.PyHandlerFileName, pyString)
	docker.CopyToContainer(containerID, def.RequirementsParentDirPath, def.RequirementsFileName, requirementsString)
	cmd := exec.Command("docker", "exec", containerID, "/bin/bash", "-c", fmt.Sprintf("'%s'", def.PyFunctionPrepareCmd)).String()
	WriteCmdToFile(def.TemplateCmdFilePath, cmd)
	command := fmt.Sprintf(`%s .`, def.TemplateCmdFilePath)
	err := exec.Command("/bin/bash", "-c", command).Run()
	if err != nil {
		fmt.Println(err)
	}
	docker.CommitContainer(containerID, imageName)
	docker.PushImage(imageName)
	docker.StopContainer(containerID)
	_, _ = docker.RemoveContainer(containerID)
}
