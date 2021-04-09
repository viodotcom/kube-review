example:
	mkdir -p resources/.kustomize/ && cp -R resources/example/* resources/.kustomize/
	cd resources/.kustomize/ && \
	kustomize edit set label app.kubernetes.io/instance:re-website-proxy-xxx && \
	cd - && kustomize build resources/.kustomize/
