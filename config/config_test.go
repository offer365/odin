package config

import (
	"fmt"
	"testing"
)

func TestConfig_LoadYaml(t *testing.T) {
	Cfg.LoadYaml("../odin.yaml")
	fmt.Println(Cfg.Peers)
	fmt.Println(Cfg)
}
