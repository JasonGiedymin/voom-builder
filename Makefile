APP_NAME=voom-builder

NO_COLOR=\033[0m
TEXT_COLOR=\033[1m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

BIN=go
DEPS=$(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
COMMON_OPTS=
PACKAGE=github.com/JasonGiedymin/voom-builder
PATH_SRC=$(GCE_SDK_GOPATH)/src/github.com/JasonGiedymin/voom-builder
PATH_PKG=$(GCE_SDK_GOPATH)/pkg/*/github.com/JasonGiedymin/voom-builder

help:
	@echo "$(OK_COLOR)-----------------------Commands:----------------------$(NO_COLOR)"
	@echo "$(TEXT_COLOR) link:        symlinks this repo to gopath $(NO_COLOR)"
	@echo "$(TEXT_COLOR) install:     alias for link $(NO_COLOR)"
	@echo "$(TEXT_COLOR) uninstall:   uninstalls $(APP_NAME) (via go) $(NO_COLOR)"
	@echo "$(TEXT_COLOR) reinstall:   uninstalls and then installs $(APP_NAME) (via go) $(NO_COLOR)"
	@echo "$(TEXT_COLOR) dev-setup    set up developer environment $(NO_COLOR)"
	@echo "$(TEXT_COLOR) deps:        install dependencies $(NO_COLOR)"
	@echo "$(TEXT_COLOR) updatedeps:  update dependencies $(NO_COLOR)"
	@echo "$(TEXT_COLOR) format:      formats the code $(NO_COLOR)"
	@echo "$(TEXT_COLOR) lint:        lints code $(NO_COLOR)"

	@echo "$(OK_COLOR) -- Dev actions -- $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test:        tests code $(NO_COLOR)"
# @echo "$(TEXT_COLOR) serve:       start and serve devserver $(NO_COLOR)"
# @echo "$(TEXT_COLOR) stop:        stop devserver $(NO_COLOR)"
# @echo "$(TEXT_COLOR) kill:        kills devserver $(NO_COLOR)"
# @echo "$(TEXT_COLOR) deploy:      deploying app $(NO_COLOR)"
# @echo "$(TEXT_COLOR) seed-dev:    seed devserver datastore $(NO_COLOR)"

# @echo "$(OK_COLOR) -- Prod actions -- $(NO_COLOR)"
# @echo "$(TEXT_COLOR) seed-prod:   seed production datastore $(NO_COLOR)"
# @echo "$(TEXT_COLOR) status-prod: check production status $(NO_COLOR)"
	@echo "$(OK_COLOR)------------------------------------------------------$(NO_COLOR)"

link:
	@echo "$(OK_COLOR)==> Symlinking project to $(PATH_SRC) $(NO_COLOR)"
	@ln -vFfsn $(shell pwd) $(PATH_SRC)

# install-pkg:
# 	@echo "$(OK_COLOR)==> Installing $(APP_NAME) $(PATH_SRC) $(NO_COLOR)"
# 	@go install $(PACKAGE)

# install: install-pkg
install: link

uninstall:
	@echo "$(OK_COLOR)==> Uninstalling $(APP_NAME) $(PATH_PKG) $(NO_COLOR)"
	@if [ -d $(PATH_PKG) ]; then rm -R $(PATH_PKG); fi;

reinstall: uninstall install

dev-setup: deps link

deps:
	@echo "$(OK_COLOR)==> Installing dependencies $(NO_COLOR)"
	@$(BIN) get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d

updatedeps:
	@echo "$(OK_COLOR)==> Updating all dependencies $(NO_COLOR)"
	@$(BIN) get -d -v -u ./...
	@echo $(DEPS) | xargs -n1 go get -d -v -u

format:
	@echo "$(OK_COLOR)==> Formatting $(NO_COLOR)"
	$(BIN) fmt ./...

lint:
	@echo "$(OK_COLOR)==> Linting $(NO_COLOR)"
	golint .

test:
	@echo "$(OK_COLOR)==> Testing $(NO_COLOR)"
	$(BIN) test -v $(COMMON_OPTS) ./registry/...

all: format lint test

# kill:
# 	@echo "$(OK_COLOR)==> Killing $(NO_COLOR)"
# 	@ps -ef | grep -m 1 "[l]ocalhost --port 8080" | awk '{print $2}' |xargs kill -9

# stop:
# 	@echo "$(OK_COLOR)==> Killing $(NO_COLOR)"
# 	@ps -ef | grep -m 1 "[l]ocalhost --port 8080" | awk '{print $2}' |xargs kill -3

# serve:
# 	@echo "$(OK_COLOR)==> Serving app $(NO_COLOR)"
# 	@$(BIN) serve registry/ &

# deploy:
# 	@echo "$(OK_COLOR)==> Deploying app $(NO_COLOR)"
# 	@$(BIN) deploy registry/

# status-prod:
# 	@echo "$(OK_COLOR)==> Checking Production status $(NO_COLOR)"
# 	@http http://registry.amuxbit.com/status

#backup-dev:
#	@appcfg.py download_data -A dev~voom-registry-service --url=http://localhost:8080/_ah/remote_api --filename=seed/backup_now.dat -e jason.giedymin@gmail.com
