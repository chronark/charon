export PATH := $(shell go env GOPATH)/src:$(PATH)
export PATH := $(shell go env GOPATH)/bin:$(PATH)
export DOCKER_BUILDKIT := 1

test:
	go test -covermode=atomic ./...

build: build-filecache build-api build-geocoding build-tiles build-rsyslog


build-rsyslog:
	docker build -t chronark/charon.rsyslog ./service/rsyslog

build-filecache:
	docker build \
	-t chronark/charon.srv.filecache \
	-f ./service/Dockerfile \
	--build-arg SERVICE=filecache \
	.

build-api:
	docker build \
	-t chronark/charon.srv.api \
	-f ./service/Dockerfile \
	--build-arg SERVICE=api \
	.

build-geocoding:
	docker build \
	-t chronark/charon.srv.geocoding \
	-f ./service/Dockerfile \
	--build-arg SERVICE=geocoding \
	.

build-tiles:
	docker build \
	-t chronark/charon.srv.tiles \
	-f ./service/Dockerfile \
	--build-arg SERVICE=tiles \
	.

build-clients:
	docker build -t chronark/charon-client-geocoding ./client/geocoding
	docker build -t chronark/charon-client-tiles ./client/tiles
	


build-map:
	git clone https://github.com/chronark/atlas.git
	docker build \
	-t chronark/atlas:A \
	--build-arg CHARON_URL=http://localhost \
	--build-arg TEST_DISPLAY_ALWAYS=true \
	./atlas/

	docker build \
	-t chronark/atlas:B \
	--build-arg CHARON_URL=http://localhost \
	./atlas/
	rm -rf atlas
	
fmt:
	./terraform fmt
	go fmt ./...
	go vet ./...
	go mod tidy
	~/go/bin/golangci-lint run ./... -v
	docker run --rm -i hadolint/hadolint < ./service/Dockerfile
	docker run --rm -i hadolint/hadolint < ./service/rsyslog/Dockerfile

init:
	[ ! -f ./terraform ] && make get-terraform || true
	./terraform init

plan: init
	./terraform plan -out tfplan

apply: plan
	./terraform apply "tfplan"
	
update: 
	git checkout master
	git pull
	make apply

purge:
	./terraform destroy -auto-approve ||true
	docker rm -f $$(docker ps -aq) || true 
	docker image rm -f $$(docker image ls -aq) || true
	docker volume rm -f $$(docker volume ls -q) || true
	docker network prune -f
	rm ./terraform ||true
	rm -rf ./volumes


netdata:
	docker run -d --name=netdata \
	-p 19999:19999 \
	-v /proc:/host/proc:ro \
	-v /sys:/host/sys:ro \
	-v /var/run/docker.sock:/var/run/docker.sock:ro \
	--cap-add SYS_PTRACE \
	--security-opt apparmor=unconfined \
	netdata/netdata

get-terraform:
	curl -o terraform.zip https://releases.hashicorp.com/terraform/0.12.19/terraform_0.12.19_linux_amd64.zip
	unzip -o terraform.zip
	rm terraform.zip

proto:
	go get github.com/micro/protoc-gen-micro/v2
	export PATH
	protoc \
		--micro_out=. \
		--go_out=. \
		./service/filecache/proto/filecache/filecache.proto

	protoc \
		--micro_out=. \
		--go_out=. \
		./service/geocoding/proto/geocoding/geocoding.proto
	
	protoc \
		--micro_out=. \
		--go_out=. \
		./service/tiles/proto/tiles/tiles.proto


