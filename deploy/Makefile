example-build:
	mkdir -p resources/.kustomize/out && cp -R resources/example/* resources/.kustomize/
	cd resources/.kustomize/ && \
	kustomize edit set label app.kubernetes.io/instance:re-website-proxy-xxx && \
	rm -rf out/* && kustomize build -o out

example-apply: example-build
	kubectl create namespace kustomize-test || true
	cd resources/.kustomize/out && kubectl apply -f . -n kustomize-test --wait=true