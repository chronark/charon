export PATH := $(shell go env GOPATH)/src:$(PATH)
export PATH := $(shell go env GOPATH)/bin:$(PATH)

build:
	docker build -t chronark/charon-service-filecache ./service/filecache
	docker build -t chronark/charon-service-gateway ./service/gateway
	docker build -t chronark/charon-service-geocoding ./service/geocoding
	docker build -t chronark/charon-service-tiles ./service/tiles
	docker build -t chronark/charon-client-geocoding ./client/geocoding
	docker build -t chronark/charon-client-tiles ./client/tiles
	docker build -t chronark/atlas ./service/map/

fmt:
	./terraform fmt
	go fmt ./...
	go vet ./...
	go mod tidy

init:
	[ ! -f ./terraform ] && make get-terraform || true
	./terraform init

plan:
	./terraform plan -out tfplan

apply:
	./terraform apply "tfplan" || echo "If you are missing docker images, please run 'make build' and try again."

prune:
	./terraform destroy -auto-approve
	@docker rm -f $$(docker ps -aq) || true 
	@docker image rm -f $$(docker image ls -aq) || true
	@docker volume rm -f $$(docker volume ls -q) || true



get-terraform:
	curl -o terraform.zip https://releases.hashicorp.com/terraform/0.12.19/terraform_0.12.19_linux_amd64.zip
	unzip -o terraform.zip
	rm terraform.zip

proto:
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


update:
	rm ./**/**/go.mod
	rm ./**/**/go.sum
	cd ./service/filecache && go clean && go mod init github.com/chronark/charon/service/filecache && go get
	cd ../gateway && go clean && go mod init github.com/chronark/charon/service/gateway && go get
	cd ../geocoding && go clean && go mod init github.com/chronark/charon/service/geocoding && go get
	cd ../tiles && go clean && go mod init github.com/chronark/charon/service/tiles && go get


	cd ../../client/geocoding && go mod init github.com/chronark/charon/client/geocoding && go get
	cd ../tiles && go mod init github.com/chronark/charon/client/tiles && go get

