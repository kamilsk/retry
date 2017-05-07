PRUNE_AVAILABLE := $(shell echo "1.13.0\n$(DOCKER_VERSION)" | sort -ct. -k1,1n -k2,2n && echo true)

.PHONY: docker-clean
docker-clean: docker-clean-invalid-common
docker-clean: docker-clean-invalid-golang
docker-clean: docker-clean-invalid-custom
docker-clean: docker-clean-invalid-tools
docker-clean:
	if [ '${PRUNE}' != '' ] && [ '${PRUNE_AVAILABLE}' == 'true' ]; then docker system prune $(strip $(PRUNE)); fi

.PHONY: docker-clean-invalid-common
docker-clean-invalid-common:
	docker images --all \
	| grep '^<none>\s\+' \
	| awk '{print $$3}' \
	| xargs docker rmi -f &>/dev/null || true

.PHONY: docker-clean-invalid-golang
docker-clean-invalid-golang:
	docker images --all \
	| grep '^golang\s\+' \
	| awk '{print $$2 "\t" $$3}' \
	| grep '^<none>\s\+' \
	| awk '{print $$2}' \
	| xargs docker rmi -f &>/dev/null || true

.PHONY: docker-clean-invalid-custom
docker-clean-invalid-custom:
	docker images --all \
	| grep '^kamilsk\/golang\s\+' \
	| awk '{print $$2 "\t" $$3}' \
	| grep '^<none>\s\+' \
	| awk '{print $$2}' \
	| xargs docker rmi -f &>/dev/null || true

.PHONY: docker-clean-invalid-tools
docker-clean-invalid-tools:
	docker images --all \
	| grep '^kamilsk\/go-tools\s\+' \
	| awk '{print $$2 "\t" $$3}' \
	| grep '^<none>\s\+' \
	| awk '{print $$2}' \
	| xargs docker rmi -f &>/dev/null || true
