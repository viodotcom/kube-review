prune-build:
	cd prune; go build -o prune

deploy-build:
	docker build  . -t kube-review:dev