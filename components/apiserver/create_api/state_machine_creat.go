package create_api

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/utils"
)

// CreateStateMachine kubectl可以不解析stateMachine的定义而直接以文件读出的字符串传输
func CreateStateMachine(cli *clientv3.Client, stateMachine def.StateMachine) {
	utils.PersistStateMachine(cli, stateMachine)
}
