package config

type AppDbConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

type AppConfig struct {
	DbConfig AppDbConfig
}

func InitFron(pathToJson string) (result *AppConfig, err error) {

	return
}
