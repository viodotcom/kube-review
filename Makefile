PIPELINE?=ci

setup:
	npm install -g codefresh

pull:
	codefresh get pipelines --name ${PIPELINE} -o yaml > .codefresh/${PIPELINE}/spec.yml

push:
	codefresh replace -f .codefresh/${PIPELINE}/spec.yml