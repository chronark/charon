provider "docker" {
}



###########################
#         Networks
###########################

resource "docker_network" "private_network" {
  name     = "internal_network"
  internal = false
}
resource "docker_network" "logging" {
  name     = "logging"
  internal = true
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
  restart = "on-failure"

  command = ["./gateway"]
  env     = ["SERVICE_ADDRESS=0.0.0.0:52000"]
  ports {
    internal = 52000
    external = 52000
  }
  networks_advanced {
    name = docker_network.private_network.name
  }
  
}

resource "docker_container" "filecache" {
  name    = "charon.service.filecache"
  image   = docker_image.filecache.latest
  restart = "on-failure"

  command = ["./filecache"]
  networks_advanced {
    name = docker_network.private_network.name
  }
  volumes {
    host_path      = "${path.cwd}/volumes/filecache"
    container_path = "/cache"
  }
}

resource "docker_container" "tiles" {
  name    = "charon.service.tiles"
  image   = docker_image.tiles.latest
  restart = "on-failure"

  command = ["./tiles"]
  env     = ["TILE_PROVIDER=osm"]
  networks_advanced {
    name = docker_network.private_network.name
  }
}

resource "docker_container" "nominatim" {
  name    = "charon.service.geocoding.nominatim"
  image   = docker_image.geocoding.latest
  restart = "on-failure"

  command = ["./geocoding"]
  env = [
    "GEOCODING_PROVIDER=nominatim",
    "DATASTORE_HOST=mongodb:27017",
  ]
  networks_advanced {
    name = docker_network.private_network.name
  }
}

resource "docker_container" "rsyslog" {
  name  = "rsyslog"
  image = "chronark/rsyslog"
  restart = "on-failure"

  networks_advanced {
    name = docker_network.logging.name
  }
  volumes {
    host_path      = "${path.cwd}/volumes/syslog"
    container_path = "/var/logs"

  }
  ports {
    internal = 514
    external = 514
    protocol = "udp"
  }

}

resource "docker_container" "logspout" {
  name  = "logspout"
  image = "gliderlabs/logspout"
  restart = "on-failure"
  networks_advanced {
    name = docker_network.logging.name
  }
  volumes {
    host_path      = "/var/run/docker.sock"
    container_path = "/var/run/docker.sock"
  }
  command = ["udp://rsyslog:514"]
}






# resource "docker_container" "geocodingclient" {
#   name  = "charon.client.geocoding"
#   image = "chronark/charon-client-geocoding"
#   networks_advanced {
#     name = docker_network.private_network.name
#   }
# }

# resource "docker_container" "tilesclient" {
#   name  = "charon.client.tiles"
#   image = "chronark/charon-client-tiles"
#   networks_advanced {
#     name = docker_network.private_network.name
#   }
# }

resource "docker_container" "atlas" {
  name  = "atlas"
  image = docker_image.atlas.latest
  ports {
    internal = 5000
    external = 80
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
}

