terraform {
  required_providers {
    outscale = {
      source  = "outscale/outscale"
      version = "1.7.0"
    }
  }
}

provider "outscale" {
}
