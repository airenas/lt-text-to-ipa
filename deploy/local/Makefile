-include Makefile.options
#####################################################################################
start: .env 
	docker-compose up -d

stop:
	docker-compose stop && docker-compose rm -f

logs:
	docker-compose logs	

clean: stop
#####################################################################################
.env: .env.in Makefile.options
	cat $< | envsubst > $@
#####################################################################################
test:
	curl http://localhost:$(port)/ipa/Jonas -X GET  | jq

.EXPORT_ALL_VARIABLES:
	

