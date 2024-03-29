package gpujob

import (
	"fmt"
	"mini-kubernetes/tools/def"
)

func GenerateGpuJobUploaderPod(job *def.GPUJob) def.Pod {
	generateImage(job)
	defualtResource := def.Resource{
		ResourceLimit: def.Limit{
			CPU:    `1`,
			Memory: `500M`,
		},
		ResourceRequest: def.Request{
			CPU:    `1`,
			Memory: `500M`,
		},
	}
	containerName := fmt.Sprintf("gpuUploader_container_%s_name", job.Name)
	podName := fmt.Sprintf("gpuUploader_pod_%s_name", job.Name)
	podLabel := fmt.Sprintf("gpuUploader_pod_%s_label", job.Name)

	job.PodName = podName

	return def.Pod{
		ApiVersion: `v1`,
		Kind:       `Pod`,
		Metadata: def.PodMeta{
			Name:  podName,
			Label: podLabel,
		},
		Spec: def.PodSpec{
			Containers: []def.Container{
				{
					Name:  containerName,
					Image: def.RgistryAddr + job.ImageName,
					PortMappings: []def.PortMapping{{
						Name:          "port_mapping_80",
						ContainerPort: 80,
						//HostPort:      80,
						Protocol: "TCP",
					}, {
						Name:          "port_mapping_22",
						ContainerPort: 22,
						//HostPort:      22,
						Protocol: "TCP",
					}},
					Resources: defualtResource,
					Commands:  []string{def.StartBash},
					Args:      []string{def.GPUJobUploaderRunArgs},
				},
			},
			Volumes: []def.Volume{},
		},
	}
}
