SHELL := /bin/sh

IMAGE ?= cms-demo
TAG ?= latest
REGISTRY_IMAGE ?= batty1039.startdedicated.net:5000/cms-demo

.PHONY: help ui backend demo production run test check docker docker-push docker-run clean

help:
	@echo "make ui         build the React app into cmd/server/web/dist"
	@echo "make backend    run the Go server on :8080 (serves the last UI build)"
	@echo "make demo       run the in-memory demo profile"
	@echo "make production run the SSO production profile"
	@echo "make run        build the UI, then run the server"
	@echo "make test       run the Go test suite"
	@echo "make check      gofmt, vet, and tests"
	@echo "make docker     build the all-in-one image"
	@echo "make docker-push build and push the image to the private registry"
	@echo "make docker-run run the image on :8080"
	@echo "make clean      remove build output"

# The Go server embeds the compiled frontend, so the UI is built first and copied
# into the embed directory that cmd/server/main.go points at.
ui:
	@if [ ! -d ui/app/node_modules ]; then \
		npm --prefix ui/app ci; \
	fi
	npm --prefix ui/app run build
	rm -rf cmd/server/web/dist
	mkdir -p cmd/server/web/dist
	cp -R ui/dist/. cmd/server/web/dist/

backend:
	set -a; [ ! -f .env ] || . ./.env; set +a; \
	go run ./cmd/server

demo:
	CMS_PROFILE=demo $(MAKE) backend

production:
	CMS_PROFILE=production $(MAKE) backend

run: ui backend

test:
	go test ./...

check:
	gofmt -l ./cmd ./internal
	go vet ./...
	go test ./...

docker:
	docker build -t $(IMAGE):$(TAG) .

docker-push: docker
	docker tag $(IMAGE):$(TAG) $(REGISTRY_IMAGE):$(TAG)
	docker push $(REGISTRY_IMAGE):$(TAG)

docker-run:
	docker run --rm -p 8080:8080 --env-file .env $(IMAGE):$(TAG)

clean:
	rm -rf cmd/server/web/dist ui/dist
	mkdir -p cmd/server/web/dist
	printf 'Placeholder so the Go embed directive resolves before the frontend is built.\nRun `make ui` to replace this directory with the real React build.\n' > cmd/server/web/dist/placeholder.txt
