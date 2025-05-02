variable "aws_region" {
  type    = string
  default = "eu-west-1"  # Adjust your region here
}

variable "influx_admin_user" {
  sensitive = true
  type    = string
  default = null
}

variable "influx_admin_pass" {
  sensitive = true
  type    = string
  default = null
}

variable "influx_url" {
  type    = string
  default = "http://influxdb:8086"
}

variable "influx_token" {
  sensitive = true
  type    = string
  default = null
}

variable "influx_org" {
  type    = string
  default = "myorg"
}

variable "influx_bucket" {
  type    = string
  default = "mybucket"
}

variable "octopus_api_key" {
  sensitive = true
  type    = string
  default = null
}

variable "octopus_base_uri" {
  sensitive = true
  type    = string
  default = null
}

variable "octopus_elec_mpan" {
  sensitive = true
  type    = string
  default = null
}

variable "octopus_elec_serial" {
  sensitive = true
  type    = string
  default = null
}

variable "octopus_gas_mprn" {
  sensitive = true
  type    = string
  default = null
}

variable "octopus_gas_serial" {
  sensitive = true
  type    = string
  default = null
}