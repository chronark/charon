fmt:
	go fmt ./...
	go vet ./...
	go mod tidy
	docker run --rm -i hadolint/hadolint < Dockerfile


build: fmt
	docker build -t chronark/charon-client-tiles .