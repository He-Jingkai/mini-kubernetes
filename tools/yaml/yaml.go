package yaml

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	Pod_t int = iota
	ClusterIP_t
	Nodeport_t
	Dns_t
	Deployment_t
	Autoscaler_t
	Gpujob_t
	Activity_t
	Unknown_t
)

type YmlSpec struct {
	Type string `yaml:"type" json:"type"`
}

type YmlMeta struct {
	Name string `yaml:"name" json:"name"`
}

type YmlObj struct {
	Kind     string  `yaml:"kind" json:"kind"`
	Metadata YmlMeta `yaml:"metadata" json:"metadata"`
	Spec     YmlSpec `yaml:"spec" json:"spec"`
}

func ReadType(path string) (int, error) {
	yml_ := &YmlObj{}
	if f, err := os.Open(path); err != nil {
		return -1, err
	} else {
		_ = yaml.NewDecoder(f).Decode(yml_)
		if (*yml_).Kind == "Pod" {
			return Pod_t, err
		} else if (*yml_).Kind == "Service" {
			if (*yml_).Spec.Type == "ClusterIP" {
				return ClusterIP_t, err
			} else if (*yml_).Spec.Type == "NodePort" {
				return Nodeport_t, err
			}
		} else if (*yml_).Kind == "DNS" {
			return Dns_t, err
		} else if (*yml_).Kind == "Deployment" {
			return Deployment_t, err
		} else if (*yml_).Kind == "HorizontalPodAutoscaler" {
			return Autoscaler_t, err
		} else if (*yml_).Kind == "GPUJob" {
			return Gpujob_t, err
		} else if (*yml_).Kind == "activity" {
			return Activity_t, err
		}
		return -1, err
	}
}

func ReadTypeAndName(path string) (int, string, error) {
	yml_ := &YmlObj{}
	if f, err := os.Open(path); err != nil {
		return -1, "", err
	} else {
		_ = yaml.NewDecoder(f).Decode(yml_)
		if (*yml_).Kind == "Pod" {
			return Pod_t, (*yml_).Metadata.Name, err
		} else if (*yml_).Kind == "Service" {
			if (*yml_).Spec.Type == "ClusterIP" {
				return ClusterIP_t, (*yml_).Metadata.Name, err
			} else if (*yml_).Spec.Type == "NodePort" {
				return Nodeport_t, (*yml_).Metadata.Name, err
			}
		} else if (*yml_).Kind == "DNS" {
			return Dns_t, (*yml_).Metadata.Name, err
		} else if (*yml_).Kind == "Deployment" {
			return Deployment_t, (*yml_).Metadata.Name, err
		} else if (*yml_).Kind == "HorizontalPodAutoscaler" {
			return Autoscaler_t, (*yml_).Metadata.Name, err
		} else if (*yml_).Kind == "GPUJob" {
			gpu_, err := ReadGPUJobConfig(path)
			if err != nil {
				fmt.Println("wrong config")
				return -1, "", err
			}
			return Gpujob_t, (*gpu_).Name, err
		} else if (*yml_).Kind == "activity" {
			act_, err := ReadFunctionConfig(path)
			if err != nil {
				fmt.Println("wrong config")
				return -1, "", err
			}
			return Activity_t, (*act_).Name, err
		}
		return -1, "", err
	}
}

func ReadPodYamlConfig(path string) (*def.Pod, error) {
	pod_ := &def.Pod{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(pod_)
		if err != nil {
			return nil, err
		}
		if (*pod_).ApiVersion != "v1" {
			fmt.Println("apiVersion should be v1!")
			return nil, err
		} else if (*pod_).Kind != "Pod" {
			fmt.Println("kind should be Pod!")
			return nil, err
		}
	}
	return pod_, nil
}

func ReadServiceClusterIPConfig(path string) (*def.ClusterIPSvc, error) {
	service_c := &def.ClusterIPSvc{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(service_c)
		if err != nil {
			return nil, err
		}
		if (*service_c).ApiVersion != "v1" {
			fmt.Println("apiVersion should be v1!")
			return nil, err
		} else if (*service_c).Kind != "Service" {
			fmt.Println("kind should be Pod!")
			return nil, err
		} else if (*service_c).Spec.Type != "ClusterIP" {
			fmt.Println("spec type should be ClusterIP!")
			return nil, err
		}
	}
	return service_c, nil
}

func ReadServiceNodeportConfig(path string) (*def.NodePortSvc, error) {
	service_n := &def.NodePortSvc{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(service_n)
		if err != nil {
			return nil, err
		}
		if (*service_n).ApiVersion != "v1" {
			fmt.Println("apiVersion should be v1!")
			return nil, err
		} else if (*service_n).Kind != "Service" {
			fmt.Println("kind should be Pod!")
			return nil, err
		} else if (*service_n).Spec.Type != "NodePort" {
			fmt.Println("spec type should be NodePort!")
			return nil, err
		}
	}
	return service_n, nil
}

func ReadDNSConfig(path string) (*def.DNS, error) {
	dns_ := &def.DNS{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(dns_)
		if err != nil {
			return nil, err
		}
		if (*dns_).Kind != "DNS" {
			fmt.Println("kind should be DNS!")
			return nil, err
		}
	}
	return dns_, nil
}

func ReadDeploymentConfig(path string) (*def.Deployment, error) {
	dep_ := &def.Deployment{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(dep_)
		if err != nil {
			return nil, err
		}
		if (*dep_).Kind != "Deployment" {
			fmt.Println("kind should be Deployment!")
			return nil, err
		}
	}
	return dep_, nil
}

func ReadAutoScalerConfig(path string) (*def.Autoscaler, error) {
	auto_ := &def.Autoscaler{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(auto_)
		if err != nil {
			return nil, err
		}
		if (*auto_).Kind != "HorizontalPodAutoscaler" {
			fmt.Println("kind should be HorizontalPodAutoscaler!")
			return nil, err
		}
	}
	return auto_, nil
}

func ReadGPUJobConfig(path string) (*def.GPUJob, error) {
	gpu_ := &def.GPUJob{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(gpu_)
		if err != nil {
			return nil, err
		}
		if (*gpu_).Kind != "GPUJob" {
			fmt.Println("kind should be GPUJob!")
			return nil, err
		}
	}
	return gpu_, nil
}

func ReadFunctionConfig(path string) (*def.Function, error) {
	func_ := &def.Function{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(func_)
		if err != nil {
			return nil, err
		}
		if (*func_).Kind != "activity" {
			fmt.Println("kind should be activity!")
			return nil, err
		}
	}
	return func_, nil
}
