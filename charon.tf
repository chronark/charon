provider "docker" {
}
resource "docker_network" "tracing" {
  name = "jaeger"
}
resource "docker_network" "logging" {
  name = "syslog"
}
resource "docker_network" "data" {
  name = "data"
}

##########################
#         Images
###########################

resource "docker_image" "gateway" {
  name          = "chronark/charon-service-gateway:latest"
  pull_triggers = ["chronark/charon-service-gateway:latest.sha256_digest"]
}

resource "docker_image" "filecache" {
  name          = "chronark/charon-service-filecache:latest"
  pull_triggers = ["chronark/charon-service-filecache:latest.sha256_digest"]
}

resource "docker_image" "tiles" {
  name          = "chronark/charon-service-tiles:latest"
  pull_triggers = ["chronark/charon-service-tiles:latest.sha256_digest"]
}


resource "docker_image" "geocoding" {
  name          = "chronark/charon-service-geocoding:latest"
  pull_triggers = ["chronark/charon-service-geocoding:latest.sha256_digest"]
}

resource "docker_image" "atlas" {
  name          = "chronark/atlas:latest"
  pull_triggers = ["chronark/atlas:latest.sha256_digest"]
}

##########################
#         Services
###########################

resource "docker_container" "gateway" {
  name    = "charon.service.gateway"
  image   = docker_image.gateway.latest
  restart = "always"

  command = ["./gateway"]
  env     = ["SERVICE_ADDRESS=0.0.0.0:52000"]
  ports {
    internal = 52000
    external = 52000
  }
  networks_advanced {
    name = docker_network.data.name
  }
   networks_advanced {
    name = docker_network.tracing.name
  }


}

resource "docker_container" "filecache" {
  name    = "charon.service.filecache"
  image   = docker_image.filecache.latest
  restart = "always"

  command = ["./filecache"]

  volumes {
    host_path      = "${path.cwd}/volumes/filecache"
    container_path = "/cache"
  }
  networks_advanced {
    name = docker_network.data.name
  }
   networks_advanced {
    name = docker_network.tracing.name
  }
}

resource "docker_container" "tiles" {
  name    = "charon.service.tiles"
  image   = docker_image.tiles.latest
  restart = "always"

  command = ["./tiles"]
  env     = ["TILE_PROVIDER=osm"]
  networks_advanced {
    name = docker_network.data.name
  }
   networks_advanced {
    name = docker_network.tracing.name
  }
}

resource "docker_container" "nominatim" {
  name    = "charon.service.geocoding.nominatim"
  image   = docker_image.geocoding.latest
  restart = "always"

  command = ["./geocoding"]
  env = [
    "GEOCODING_PROVIDER=nominatim",
  ]
  networks_advanced {
    name = docker_network.data.name
  }
   networks_advanced {
    name = docker_network.tracing.name
  }
}

resource "docker_container" "rsyslog" {
  name    = "rsyslog"
  image   = "chronark/rsyslog"
  restart = "always"


  volumes {
    host_path      = "${path.cwd}/volumes/syslog"
    container_path = "/var/logs"

  }
  ports {
    internal = 514
    external = 514
    protocol = "udp"
  }

   networks_advanced {
    name = docker_network.logging.name
  }

}

resource "docker_container" "logspout" {
  name    = "logspout"
  image   = "gliderlabs/logspout"
  restart = "always"

  volumes {
    host_path      = "/var/run/docker.sock"
    container_path = "/var/run/docker.sock"
  }
  command = ["udp://rsyslog:514"]
  networks_advanced {
    name = docker_network.logging.name
  }
}





resource "docker_container" "atlas" {
  name  = "atlas"
  image = docker_image.atlas.latest
  ports {
    internal = 80
    external = 80
  }
  restart = "always"
  networks_advanced {
    name = docker_network.data.name
  }
  
}

resource "docker_container" "portainer" {
  name  = "portainer"
  image = "portainer/portainer"
  ports {
    internal = 8000
    external = 8000
  }
  ports {
    internal = 9000
    external = 9000
  }
  volumes {
    host_path      = "/var/run/docker.sock"
    container_path = "/var/run/docker.sock"
  }
  restart = "always"
}

resource "docker_container" "jaeger" {
  name  = "jaeger"
  image = "jaegertracing/all-in-one:latest"
  env   = ["COLLECTOR_ZIPKIN_HTTP_PORT=9411"]
  ports {
    internal = 5775
    external = 5775
    protocol = "udp"
  }

  ports {
    // accept jaeger.thrift in compact Thrift protocol used by most current Jaeger clients
    internal = 6831
    external = 6831
    protocol = "udp"
  }
  ports {
    internal = 6832
    external = 6832
    protocol = "udp"
  }
  ports {
    // UI
    internal = 16686
    external = 16686
  }
  ports {
    // Healthcheck at / and metrics at /metrics
    internal = 14268
    external = 14268
  }
  ports {
    internal = 9411
    external = 9411
  }
  networks_advanced {
    name = docker_network.tracing.name
  }

}
