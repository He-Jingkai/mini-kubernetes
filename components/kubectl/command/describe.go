package command

import (
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/utils"

	"github.com/urfave/cli"
)

func NewDescribeCommand() cli.Command {
	return cli.Command{
		Name:  "describe",
		Usage: "Describe resources according name",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			describeFunc(c)
			return nil
		},
	}
}

func describeFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		wrong("Unavailable command!")
		return
	} else if len(c.Args()) < 2 {
		wrong("You need to specify resource name")
		return
	}

	if c.Args()[0] == "pod" {
		// kubectl describe pod podName
		podName := c.Args()[1]
		response := def.Pod{}
		err, status := httpget.Get("http://" + utils.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get_pod/" + podName).
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			wrong(err.Error())
			return
		}
		// fmt.Printf("get_pod status is %s\n", status)
		if status == "200" {
			if res, err := json.MarshalIndent(response, "", "   "); err == nil {
				fmt.Println("get pod successfully and the response is:\n", string(res))
			}
		} else {
			fmt.Printf("pod_ %s doesn't exist\n", podName)
		}
	} else if c.Args()[0] == "service" {
		// kubectl describe service serviceName
		// 用来获取特定名称的 service，需要发送给apiserver的参数为 serviceName(string)
		serviceName := c.Args()[1]
		response := def.Service{}
		err, status := httpget.Get("http://" + utils.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/service/" + serviceName).
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			wrong(err.Error())
			return
		}
		// fmt.Printf("get_service status is %s\n", status)
		if status == "200" {
			if res, err := json.MarshalIndent(response, "", "   "); err == nil {
				fmt.Println("get service successfully and the response is:\n", string(res))
			}
		} else {
			fmt.Printf("service %s doesn't exist\n", serviceName)
		}
	} else if c.Args()[0] == "dns" {
		// kubectl describe dns dnsName
		// 用来获取特定名称的 dns，需要发送给apiserver的参数为 dnsName(string)
		// http调用返回的json需解析转为def.DNSDetail类型，
		dnsName := c.Args()[1]
		response := def.DNSDetail{}
		err, status := httpget.Get("http://" + utils.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/dns/" + dnsName).
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			wrong(err.Error())
			return
		}
		// fmt.Printf("get_dns status is %s\n", status)
		if status == "200" {
			if res, err := json.MarshalIndent(response, "", "   "); err == nil {
				fmt.Println("get dns successfully and the response is:\n", string(res))
			}
		} else {
			fmt.Printf("dns %s doesn't exist\n", dnsName)
		}
	} else if c.Args()[0] == "deployment" {
		// kubectl describe deployment deploymentName
		// 用来获取特定名称的 deployment，需要发送给apiserver的参数为 deploymentName(string)
		// http调用返回的json需解析转为def.DeploymentDetail类型
		deploymentName := c.Args()[1]
		response := def.DeploymentDetail{}
		err, status := httpget.Get("http://" + utils.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/deployment/" + deploymentName).
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			wrong(err.Error())
			return
		}
		// fmt.Printf("get_deployment status is %s\n", status)
		if status == "200" {
			if res, err := json.MarshalIndent(response, "", "   "); err == nil {
				fmt.Println("get deployment successfully and the response is:\n", string(res))
			}
		} else {
			fmt.Printf("deployment %s doesn't exist\n", deploymentName)
		}
	} else if c.Args()[0] == "autoscaler" {
		// kubectl describe autoscaler autoscalerName
		// 用来获取特定名称的 autoscaler，需要发送给apiserver的参数为 autoscalerName(string)
		autoscalerName := c.Args()[1]
		response := def.AutoscalerDetail{}
		err, status := httpget.Get("http://" + utils.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/autoscaler/" + autoscalerName).
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			wrong(err.Error())
			return
		}
		// fmt.Printf("get_autoscaler status is %s\n", status)
		if status == "200" {
			if res, err := json.MarshalIndent(response, "", "   "); err == nil {
				fmt.Println("get autoscaler successfully and the response is:\n", string(res))
			}
		} else {
			fmt.Printf("autoscaler %s doesn't exist\n", autoscalerName)
		}
	} else if c.Args()[0] == "gpujob" {
		// kubectl describe gpujob gpuJobName
		// 用来获取特定名称的 gpuJob，需要发送给apiserver的参数为 gpuJobName(string)
		gpuJobName := c.Args()[1]
		response := def.GPUJobDetail{}
		err, status := httpget.Get("http://" + utils.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/gpuJob/" + gpuJobName).
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			wrong(err.Error())
			return
		}
		// fmt.Printf("get_gpuJob status is %s\n", status)
		if status == "200" {
			if res, err := json.MarshalIndent(response, "", "   "); err == nil {
				fmt.Printf("get gpuJob successfully and the response is: %v\n", string(res))
			}
		} else {
			fmt.Printf("gpuJob %s doesn't exist\n", gpuJobName)
		}
	} else if c.Args()[0] == "function" {
		// kubectl describe function functionName
		// 用来获取特定名称的 function，需要发送给apiserver的参数为 functionName(string)
		functionName := c.Args()[1]
		response := def.Function{}
		err, status := httpget.Get("http://" + utils.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/function/" + functionName).
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			wrong(err.Error())
			return
		}
		// fmt.Printf("get_function status is %s\n", status)
		if status == "200" {
			if res, err := json.MarshalIndent(response, "", "   "); err == nil {
				fmt.Printf("get function successfully and the response is: %v\n", string(res))
			}
		} else {
			fmt.Printf("function %s doesn't exist\n", functionName)
		}
	} else if c.Args()[0] == "statemachine" {
		// kubectl describe statemachine stateMachineName
		// 用来获取特定名称的 StateMachine，需要发送给apiserver的参数为 stateMachineName(string)
		stateMachineName := c.Args()[1]
		response := def.StateMachine{}
		err, status := httpget.Get("http://" + utils.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/stateMachine/" + stateMachineName).
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			wrong(err.Error())
			return
		}
		// fmt.Printf("get_stateMachine status is %s\n", status)
		if status == "200" {
			if res, err := json.MarshalIndent(response, "", "   "); err == nil {
				fmt.Printf("get stateMachine successfully and the response is: %v\n", string(res))
			}
		} else {
			fmt.Printf("stateMachine %s doesn't exist\n", stateMachineName)
		}
	} else {
		wrong("Wrong resource type or name")
	}
}
