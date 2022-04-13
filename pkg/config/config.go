package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func LoadEnv(env string) {
	if strings.Contains(env, "/") || strings.Contains(env, `\`) {
		err := godotenv.Load(env)
		log.WithError(err)
	} else {
		err := godotenv.Load("./config/" + env)
		log.WithError(err)
	}
}

func ShowListEnvs() {
	log.Debug("list of envs:")
	for _, each := range os.Environ() {
		log.Debug("  ==> ", each)
	}
}
