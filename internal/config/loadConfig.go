package config

import (
	_ "embed"

	"github.com/Anacardo89/mailer_sender/internal/mail"
	"github.com/Anacardo89/mailer_sender/pkg/logger"
	"github.com/Anacardo89/mailer_sender/pkg/rabbit"
	"gopkg.in/yaml.v2"
)

//go:embed mailConfig.yaml
var mailYaml []byte

//go:embed rabbitConfig.yaml
var rabbitYaml []byte

func LoadMailConfig() *mail.Config {
	var config *mail.Config
	err := yaml.Unmarshal(mailYaml, &config)
	if err != nil {
		logger.Error.Fatal(err)
		return nil
	}
	return config
}

func LoadRabbitConfig() *rabbit.Config {
	var config rabbit.Config
	err := yaml.Unmarshal(rabbitYaml, &config)
	if err != nil {
		logger.Error.Fatal(err)
		return nil
	}
	return &config
}
