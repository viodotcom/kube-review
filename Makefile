prune-build:
	cd src/prune; go build -o prune

deploy-build:
	docker build  . -t kube-review:dev