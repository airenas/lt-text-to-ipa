#####################################################################################
dist_dir=$(CURDIR)/data
executable_name=text-to-ipa
main_dir=../cmd
service=airenas/text-to-ipa
version=0.1
commit_count=$(shell git rev-list --count HEAD)
go_build_cmd=CGO_ENABLED=0 go build -ldflags "-X main.version=$(version).$(commit_count)" 
#####################################################################################
$(dist_dir):
	mkdir -p $@
$(dist_dir)/$(executable_name): | $(dist_dir)
	cd $(main_dir)/$(executable_name) && $(go_build_cmd) -o $@
#####################################################################################
build: $(dist_dir)/$(executable_name)
########### DOCKER ##################################################################
tag=$(service):$(version).$(commit_count)

dbuild: $(dist_dir)/$(executable_name)
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
#####################################################################################
clean:
	rm -rf $(dist_dir)

.PHONY:
	clean build dbuild dpush
