package function

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/utils"
)

func MakeFunctionalImage(function *def.Function) {
	pyString := utils.ReadFile(function.Function)
	requirementsString := utils.ReadFile(function.Requirements)
	container := def.Container{
		Image: def.PyFunctionTemplateImage,
	}
	imageName := fmt.Sprintf("image_%s_%d", function.Name, function.Version)
	function.Image = imageName
	containerID := docker.CreateContainer(container, imageName)
	docker.StartContainer(containerID)
	docker.CopyToContainer(containerID, def.PyHandlerParentDirPath, def.PyHandlerFileName, pyString)
	docker.CopyToContainer(containerID, def.RequirementsParentDirPath, def.RequirementsFileName, requirementsString)
	docker.DockerExec(containerID, []string{def.StartBash, def.PyFunctionPrepareFile})

	docker.CommitContainer(containerID, imageName)
	docker.PushImage(imageName)
	docker.StopContainer(containerID)
	_, _ = docker.RemoveContainer(containerID)
}
