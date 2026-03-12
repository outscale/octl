package main

import (
	"fmt"

	"github.com/outscale/octl/pkg/config"
	"github.com/samber/lo"
)

func main() {
	cfg := config.For("iaas")
	for call := range cfg.Calls {
		if lo.ContainsBy(cfg.Aliases, func(a config.Alias) bool {
			return a.AliasTo == call
		}) {
			continue
		}
		fmt.Println(call, "is missing")
	}
}
