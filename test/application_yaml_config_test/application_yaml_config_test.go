package application_yaml_config_test

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/gateway"
	"testing"
)

func Test(t *testing.T) {
	detail := def.DNSDetail{
		Kind: "dns",
		Name: "test",
		Host: "test.example.com",
		Paths: []def.PathPairDetail{
			def.PathPairDetail{
				Path: "testPath",
				Port: 8080,
				Service: def.ClusterIPSvc{
					Spec: def.Spec{
						ClusterIP: "1.1.1.1",
					},
				},
			},
		},
	}
	fmt.Println(fmt.Sprintf(
		"echo -e \"%s\" > %s",
		gateway.GenerateApplicationYaml(detail),
		def.GatewayRoutesConfigPathInImage))
	t.Log("test finished\n")
}
