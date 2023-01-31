GO = go
GOBIN ?= $(PWD)/bin
PATH := $(GOBIN):$(PATH)

.PHONY:run
run:
	$(GO) run cmd/main.go


.PHONY: db-up
db-up:
	@docker-compose -f docker-compose.yml up --detach

