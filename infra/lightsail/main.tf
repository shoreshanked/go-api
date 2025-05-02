provider "aws" {
  region = var.aws_region
}

# Lightsail container service creation
resource "aws_lightsail_container_service" "go_api_service" {
  name  = "go-api-container"
  power = "nano" # Adjust power (CPU/Memory) as needed
  scale = 1      # Number of containers to run
  deployment {
    containers {
      name  = "go-api"
      image = "229418028078.dkr.ecr.us-west-2.amazonaws.com/go-api:latest"
      ports = ["80"] # Port mapping for your Go API

      environment = {
        INFLUX_URL    = var.influx_url
        INFLUX_TOKEN  = var.influx_token
        INFLUX_ORG    = var.influx_org
        INFLUX_BUCKET = var.influx_bucket
        API_KEY       = var.octopus_api_key
        BASE_URI      = var.octopus_base_uri
        ELEC_MPAN     = var.octopus_elec_mpan
        ELEC_SERIAL   = var.octopus_elec_serial
        GAS_MPRN      = var.octopus_gas_mprn
        GAS_SERIAL    = var.octopus_gas_serial
      }

    }

    containers {
      name  = "influxdb"
      image = "influxdb:latest" # InfluxDB image
      environment = {
        DOCKER_INFLUXDB_INIT_MODE        = "setup"
        DOCKER_INFLUXDB_INIT_USERNAME    = var.influx_admin_user # e.g. "admin"
        DOCKER_INFLUXDB_INIT_PASSWORD    = var.influx_admin_pass # e.g. secure secret
        DOCKER_INFLUXDB_INIT_ORG         = var.influx_org
        DOCKER_INFLUXDB_INIT_BUCKET      = var.influx_bucket
        DOCKER_INFLUXDB_INIT_ADMIN_TOKEN = var.influx_token
      }
      ports = ["8086"] # Default port for InfluxDB
    }

  }
}