terraform {
  required_providers {
    outscale = {
      source  = "outscale/outscale"
      version = "1.6.0"
    }
  }
}

provider "outscale" {
}
