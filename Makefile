PIPELINE?=build
STAGE?=stg

setup:
	brew tap codefresh-io/cli
	brew install gomplate codefresh
	codefresh auth create-context --api-key ${CFTOKEN}

pull:
	codefresh get pipelines --name ${PIPELINE} -o yaml

spec:
	echo 'STAGE: ${STAGE}' | gomplate -f .codefresh/${PIPELINE}/spec.yml -d data=stdin:///foo.yaml > /tmp/${PIPELINE}-spec.yml

create-pipeline: spec
	codefresh create -f /tmp/${PIPELINE}-spec.yml

update-pipeline: spec
	codefresh replace -f /tmp/${PIPELINE}-spec.yml

create:
	PIPELINE=build make create-pipeline
	PIPELINE=auto-prune make create-pipeline

update:
	PIPELINE=build make update-pipeline
	PIPELINE=auto-prune make update-pipeline