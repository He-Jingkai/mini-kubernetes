GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
TARGET_ACTIVER=activer
TARGET_KUBELET=kubelet
TARGET_APISERVER=apiserver
TARGET_CONTROLLER=controller
TARGET_KUBECTL=kubectl
TARGET_SCHEDULER=scheduler
TARGET_PROXY=proxy
.DEFAULT_GOAL := default

GO_TEST_PATH= './test/yaml_test'

.PHONY:test

all: test master node

master: kubectl apiServer controller scheduler activer

node: kubelet proxy

default: master node

test:
	for path in ${GO_TEST_PATH}; do \
		$(GO_TEST) $$path -v ; \
	done

run_node:
	wget -O cadvisor https://github.com/google/cadvisor/releases/download/v0.46.0/cadvisor-v0.46.0-linux-amd64
	$(GO_BUILD) ./components/kubelet/kubelet.go -o TARGET_KUBELET
    $(GO_BUILD) ./components/kubeproxy/proxy.go -o TARGET_PROXY
    ./cadvisor -port=8090 & /
    ./proxy &

run_master:
	$(GO_BUILD) ./components/apiserver/apiserver.go -o TARGET_APISERVER
    $(GO_BUILD) ./components/controller/controller.go -o TARGET_CONTROLLER
    $(GO_BUILD) ./components/scheduler/scheduler.go -o -o TARGET_SCHEDULER
    $(GO_BUILD) ./components/activer/activer.go -o TARGET_ACTIVER
    $(GO_BUILD) ./components/kubectl/kubectl.go -o TARGET_KUBECTL
    ./master & \
    ./controller & \
    ./scheduler & \
    ./activer &
