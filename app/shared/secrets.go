package shared

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

func GetSecrets() *koanf.Koanf {
	k := koanf.New("/")
	err := k.Load(file.Provider("config/secrets/env.yaml"), yaml.Parser())
	if err != nil {
		log.Fatal("Env file not found")

	}
	return k
}
