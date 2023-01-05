package gpujob

import (
	"fmt"
	"github.com/jakehl/goid"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/utils"
)

func apiServerIPFileGenerator() string {
	//tempFilePath := def.TemplateFileDir + goid.NewV4UUID().String()
	//fp, err := os.Create(tempFilePath)
	//if err != nil {
	//	return ""
	//}
	//defer func(fp *os.File) {
	//	_ = fp.Close()
	//}(fp)
	//_, _ = fp.WriteString(fmt.Sprintf("%s:%d\n", utils.GetLocalIP().String(), def.MasterPort))
	//return tempFilePath
	return fmt.Sprintf("%s:%d", utils.GetLocalIP().String(), def.MasterPort)
}

func generateImage(job *def.GPUJob) {
	newImageName := fmt.Sprintf("gpuJob-%s-%s", job.Name, goid.NewV4UUID().String())
	slurmContent := slurmGenerator(job.Slurm)
	apiServerIPFileContent := apiServerIPFileGenerator()
	sourceCodeContent := utils.ReadFile(job.SourceCodePath)
	makefileContent := utils.ReadFile(job.MakefilePath)
	container := def.Container{
		Image: def.GPUJobUploaderImage,
	}
	containerID := docker.CreateContainer(container, newImageName)
	docker.CopyToContainer(containerID, def.GPUSlurmScriptParentDirPath, def.GPUSlurmScriptFileName, slurmContent)
	docker.CopyToContainer(containerID, def.GPUApiServerIpAndPortFileParentDirPath, def.GPUApiServerIpAndPortFileFileName, apiServerIPFileContent)
	docker.CopyToContainer(containerID, def.GPUJOBMakefileParentDirPath, def.GPUJOBMakefileFileName, makefileContent)
	docker.CopyToContainer(containerID, def.GPUJobSourceCodeParentDirPath, def.GPUJobSourceCodeFileName, sourceCodeContent)
	docker.CopyToContainer(containerID, def.GPUJobNameParentDirName, def.GPUJobNameFileName, job.Name)

	docker.CommitContainer(containerID, newImageName)
	docker.PushImage(newImageName)
	docker.StopContainer(containerID)
	_, _ = docker.RemoveContainer(containerID)
	job.ImageName = newImageName
}
