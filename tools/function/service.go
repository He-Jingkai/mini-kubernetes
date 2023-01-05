package function

import (
	"fmt"
	"github.com/jakehl/goid"
	"mini-kubernetes/tools/def"
)

func GenerateFunctionPodAndService(function *def.Function) (pod def.Pod, service def.ClusterIPSvc) {
	MakeFunctionalImage(function)
	defaultContainerResource := def.Resource{
		ResourceLimit: def.Limit{
			CPU:    `3`,
			Memory: `3G`,
		},
		ResourceRequest: def.Request{
			CPU:    `1`,
			Memory: `500Mi`,
		},
	}
	id := fmt.Sprintf("functional_%s_%d_%s", function.Name, function.Version, goid.NewV4UUID().String())
	containerName := fmt.Sprintf("image_%s", id)
	podName := fmt.Sprintf("pod_%s", id)
	podLabel := fmt.Sprintf("pod_%s_label", id)
	serviceName := fmt.Sprintf("service_%s", id)
	pod = def.Pod{
		ApiVersion: `v1`,
		Kind:       `Pod`,
		Metadata: def.PodMeta{
			Name:  podName,
			Label: podLabel,
		},
		Spec: def.PodSpec{
			Containers: []def.Container{
				{
					Name:         containerName,
					Image:        def.RgistryAddr + function.Image,
					Commands:     []string{def.StartBash},
					Args:         []string{def.PyFunctionStartArgs},
					WorkingDir:   "",
					VolumeMounts: []def.VolumeMount{},
					PortMappings: []def.PortMapping{
						{
							Name:          fmt.Sprintf("portmaping_%s", id),
							ContainerPort: 80,
							//HostPort:      80,
							Protocol: "TCP",
						},
					},
					Resources: defaultContainerResource,
				},
			},
			Volumes: []def.Volume{},
		},
	}

	service = def.ClusterIPSvc{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata:   def.Meta{Name: serviceName},
		Spec: def.Spec{
			Type:      "ClusterIP",
			ClusterIP: "",
			Ports: []def.PortPair{{
				Port:       80,
				TargetPort: `80`,
				Protocol:   "TCP",
			}},
			Selector: def.Selector{Name: podLabel},
		},
	}
	function.PodName = podName
	function.ServiceName = serviceName
	return
}
