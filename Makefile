prune-build:
	cd src/prune; go build -o prune

api-build:
	cd src/api; go build -o api

deploy-build:
	docker build  . -t kube-review:dev