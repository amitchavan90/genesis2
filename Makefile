# use bash shell
SHELL:=/bin/bash
# Go parameters
PROD_USER=genesis
# PROD_HOST=x.x.x.x
PROD_VERSION=$(shell cat server/cmd/platform/main.go |grep 'const version'| grep -oP 'v[0-9]+.[0-9]+.[0-9]+')
PROD_PATH_UPLOAD=/usr/share/latitude28/genesis-$(PROD_VERSION)
PROD_PATH_ONLINE=/usr/share/latitude28/genesis-online
PROD_PATH=/usr/share/latitude28/genesis
PROD_BASE=$(shell dirname $(PROD_PATH))
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

SERVER_BINARY_PATH=../deploy/bin/genesis-server

all: getver precheck clean prepare deps build-server build-web copy-server copy-web
getver:
	@echo version is $(PROD_VERSION)
precheck:
	@if [[ "$(BRANCH)" != "master" ]]; then \
		echo "ERROR: switch to master branch first"; \
		exit 1; \
	fi
prepare: 
	mkdir deploy
build-server: 
	cd server && go build -o $(SERVER_BINARY_PATH) cmd/platform/main.go
build-web:
	cd web && npm install
	cd web && npm run build
clean:
	rm -rf deploy
copy-server:
	rsync 	 server/*.html              deploy/
	rsync    bin/migrate                deploy/bin
	rsync -r init                       deploy/
	rsync -r configs                    deploy/
	rsync -r server/migrations          deploy/
	rsync scripts/db-prod_migrate.sh    deploy/scripts/
copy-web:
	rsync -r web/dist/                  deploy/web
deps:
	cd server && go mod download
deploy-prod-full: remote-mkdir-upload rsync-prod-full
deploy-prod-frontend: remote-mkdir-upload rsync-prod-frontend
deploy-prod-backend: remote-mkdir-upload rsync-prod-backend
rsync-prod-full:
	rsync -avz --delete --exclude "*.js.map" ./deploy/ $(PROD_USER)@$(PROD_HOST):$(PROD_PATH_UPLOAD)
rsync-prod-frontend:
	rsync -avz --delete --exclude "*.js.map" ./deploy/web/ $(PROD_USER)@$(PROD_HOST):$(PROD_PATH_UPLOAD)/web
rsync-prod-backend:
	rsync -avz --delete ./deploy/bin/         $(PROD_USER)@$(PROD_HOST):$(PROD_PATH_UPLOAD)/bin/
	rsync -avz --delete ./deploy/configs/     $(PROD_USER)@$(PROD_HOST):$(PROD_PATH_UPLOAD)/configs/
	rsync -avz --delete ./deploy/init/        $(PROD_USER)@$(PROD_HOST):$(PROD_PATH_UPLOAD)/init/
	rsync -avz --delete ./deploy/migrations/  $(PROD_USER)@$(PROD_HOST):$(PROD_PATH_UPLOAD)/migrations/
	rsync -avz --delete ./deploy/scripts/     $(PROD_USER)@$(PROD_HOST):$(PROD_PATH_UPLOAD)/scripts/
	rsync -avz --delete ./deploy/*.html       $(PROD_USER)@$(PROD_HOST):$(PROD_PATH_UPLOAD)/
remote-mkdir-online:
	ssh root@$(PROD_HOST) 'mkdir -p $(PROD_PATH_ONLINE); ln -s $(PROD_PATH_ONLINE) $(PROD_PATH)'
remote-mkdir-upload:
	ssh root@$(PROD_HOST) 'mkdir -p $(PROD_PATH_UPLOAD); chown -R $(PROD_USER):$(PROD_USER) $(PROD_BASE)'
remote-sysd:
	ssh root@$(PROD_HOST) "cd /etc/systemd/system; ln -s $(PROD_PATH)/init/$(PROD_USER).service; systemctl daemon-reload"
remote-install: remote-mkdir-upload remote-mkdir-online deploy-prod-full remote-sysd
