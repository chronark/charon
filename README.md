# Charon

[![CodeFactor](https://www.codefactor.io/repository/github/chronark/charon/badge)](https://www.codefactor.io/repository/github/chronark/charon)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fchronark%2Fcharon.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fchronark%2Fcharon?ref=badge_shield)

Backend services for caching and displaying geospatial data

## Getting Started

These instructions will get you a copy of the project up and running on your system. This is intended to run on an ubuntu machine but can be run anywhere you have docker installed.

### Dependencies

- make
- unzip (to unzip terraform)
- docker

Make is most likely already installed on your system.
To install docker follow the instructions [here](https://docs.docker.com/install/linux/docker-ce/ubuntu/) or use the install script:

```sh
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

### Deployment

There are a couple of steps to run the services.

0. Clone the repository to your machine.

    ```sh
    git clone https://github.com/chronark/charon.git
    cd charon
    ```

1. Build all necessary docker images if you don't have them alrady.

    ```sh
    make build
    ```

2. Get and initialize Terraform, a infrastructure management and provisioning tool.

    ```sh
    make init
    ```

3. Create a terraform plan

    ```sh
    make plan
    ```

4. Apply the plan

    ```sh
    make apply
    ```

Visit `http://localhost` to see the map.

## Services

![Architecture](https://raw.githubusercontent.com/chronark/charon/master/architecture.svg?sanitize=true)

### Filecache

Discovery name: `charon.srv.filecache`

Filecache is a basic cache that writes or reads bytes to and from the disk.
The cache is mounted on the host as `./volumes/filecache` and can be backed up or modified by removing or adding files on the host machine.

### API

The HTTP API to interact with all microservices. Exposes various routes to fetch data.
Listens on port 52000
Routes:

- `/geocoding/forward?query=XYZ` calls the geocoding service.
- `/geocoding/reverse?lat=1.0&lon=1.0` calls the geocoding service.
- `/tile?x=0&y=0&z=0`calls the tiles service.

### Geocoding

Discovery name: `charon.srv.geocoding.nominatim` or `charon.srv.geocoding.mapbox`

Geocoding has two handlers:  [Nominatim](https://nominatim.org/) and [Mapbox](https://www.mapbox.com/)
They both receive a forward or reverse geocoding search and search the cache. In case of a miss they will query their respective 3rd party APIs, write to cache and return the result.
The cache is mounted on the host as `./volumes/geocoding` and can be backed up or modified by removing or adding files on the host machine.

### Tiles

Discovery name: `charon.srv.tiles.osm` or `charon.srv.geocoding.mapbox`

Tiles has two handlers:  [osm](https://openstreetmap.org/) and [Mapbox](https://www.mapbox.com/)
They both receive a tile request and search the cache. In case of a miss they will query their respective 3rd party APIs, write to cache and return the result.

## Built With

- [go](https://golang.org/)
- [docker](https://www.docker.com/)
- [gRPC](https://grpc.io/) - Internal communication
- [terraform](https://www.terraform.io/) - Infrastructure management
- [go-micro](https://github.com/micro/go-micro/v2/v2) - Microservice framework

## Development

To get up and running as fast as possible you can use Vagrant to start and provision a local ubuntu:18.04 machine

```sh
vagrant up
```

Charon makes use of protobuf for the internal communication. In order to compile the protobuf definition you need to install:

```sh
go get google.golang.org/grpc
go get github.com/golang/protobuf/protoc-gen-go
go get github.com/micro/protoc-gen-micro
```

Then run `make proto` to compile all protobuf definitions at once.

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fchronark%2Fcharon.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fchronark%2Fcharon?ref=badge_large)
