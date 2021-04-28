-include Makefile.options

HOME_DIR?=/home/$(USER)

SSH_CMD=ssh -t $(USER)@$(SERVER)
RSYNC_CMD=rsync -r -e 'ssh' -P

#####################################################################################
remote_service_dir=tservices
dc_cmd=docker-compose
#####################################################################################
########### SERVICES ################################################################
services/.env: services/.env.in Makefile.options
	cat services/.env.in | envsubst > $@
copy: services/.env
	$(RSYNC_CMD) services/* $(USER)@$(SERVER):$(HOME_DIR)/$(remote_service_dir)/
	$(RSYNC_CMD) services/.env* $(USER)@$(SERVER):$(HOME_DIR)/$(remote_service_dir)/

restart: copy
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && $(dc_cmd) stop && $(dc_cmd) up -d"

deploy-one: copy
	@$(if $(strip $(srv)),echo "Service = $(srv)",echo No srv set && exit 1)
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && $(dc_cmd) stop $(srv) && \
		$(dc_cmd) rm -f $(srv) && $(dc_cmd) up -d $(srv)"

redeploy-one: copy 
	@$(if $(strip $(srv)),echo "Service = $(srv)",echo No srv set && exit 1)
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && $(dc_cmd) up --detach --build $(srv)"

log:
	@$(if $(strip $(srv)),echo "Service = $(srv)",echo No srv set && exit 1)
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && $(dc_cmd) logs $(srv)"

status:
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && $(dc_cmd) ps"
ssh:
	$(SSH_CMD)
########### END SERVICES ############################################################

.PHONY:
	copy deploy clean status ssh log redeploy-one deploy-one
	
.EXPORT_ALL_VARIABLES:	