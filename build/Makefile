-include ../version
#####################################################################################
service=airenas/text-to-ipa
########### DOCKER ##################################################################
tag=$(service):$(version)

dbuild:
	cd .. && docker build -t $(tag) --build-arg BUILD_VERSION=$(version) -f build/Dockerfile .

dpush: dbuild
	docker push $(tag)
#####################################################################################
.PHONY: dbuild dpush
#####################################################################################
