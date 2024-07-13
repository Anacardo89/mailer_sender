package main

import (
	_ "embed"

	"gopkg.in/yaml.v2"
)

//go:embed mailConfig.yaml
var mailYaml []byte

//go:embed rabbitConfig.yaml
var rabbitYaml []byte

type MailConfig struct {
	SmtpHost     string `yaml:"smtp_host"`
	SmtpPort     string `yaml:"smtp_port"`
	SmtpUser     string `yaml:"smtp_user"`
	SmtpPass     string `yaml:"smtp_pass"`
	PublicSender string `yaml:"public_sender"`
}

type RabbitConfig struct {
	RabbitUser    string `yaml:"rabbit_user"`
	RabbitPass    string `yaml:"rabbit_pass"`
	RabbitHost    string `yaml:"rabbit_host"`
	RabbitPort    string `yaml:"rabbit_port"`
	QueueRegister string `yaml:"queue_ragister`
}

func loadMailConfig() (*MailConfig, error) {
	var mailConfig *MailConfig
	err := yaml.Unmarshal(mailYaml, mailConfig)
	if err != nil {
		return nil, err
	}
	return mailConfig, nil
}

func loadRabbitConfig() (*RabbitConfig, error) {
	var rabbitConfig RabbitConfig
	err := yaml.Unmarshal(rabbitYaml, &rabbitConfig)
	if err != nil {
		return nil, err
	}
	return &rabbitConfig, nil
}
