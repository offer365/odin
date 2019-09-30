package config

import (
	"fmt"
	"testing"
)

func TestConfig_LoadYaml(t *testing.T) {
	Cfg.LoadYaml("../odin.yaml")
	fmt.Println(Cfg)
}

func TestConfig_LoadJson(b *testing.T) {
	Cfg.LoadJson("../odin.json")
	fmt.Println(Cfg)
}
