example-replace:
	export NAMESPACE=re-website-proxy-xxx
	export DOMAIN=shared-stg.fih.io
	@for f in ./resources/.kustomize/base/patches/*.json; do \
		cat $$f | envsubst; \
	done

example-build:
	mkdir -p resources/.kustomize/out resources/.kustomize/base &&  \
	cp -R resources/example/* resources/.kustomize && \
	cp -R resources/base/* resources/.kustomize/base && \
	cd resources/.kustomize/ && \
	kustomize edit set label app.kubernetes.io/instance:re-website-proxy-xxx && \
	rm -rf out/* && kustomize build -o out

example-apply: example-build
	kubectl create namespace kustomize-test || true
	cd resources/.kustomize/out && kubectl apply -f . -n kustomize-test --wait=true