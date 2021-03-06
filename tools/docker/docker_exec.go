package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func DockerExec(containerID string, commands []string) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer func(cli *client.Client) {
		_ = cli.Close()
	}(cli)
	copyfiletoBin, err := cli.ContainerExecCreate(context.Background(), containerID, types.ExecConfig{
		Tty:          true,
		Cmd:          commands,
		AttachStdin:  true,
		AttachStdout: true,
		Detach:       true,
	})
	if err != nil {
		fmt.Println("err:   ", err)
		panic(err)
	}
	err = cli.ContainerExecStart(context.Background(), copyfiletoBin.ID, types.ExecStartCheck{
		Tty: true,
	})
	for {
		inspect, _ := cli.ContainerExecInspect(context.Background(), copyfiletoBin.ID)
		if !inspect.Running {
			return
		}
	}
}
