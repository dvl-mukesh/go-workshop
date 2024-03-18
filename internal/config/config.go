package config

type Environment struct {
	DbUserName string `env:"DB_USERNAME,required"`
	DbPassword string `env:"DB_PASSWORD,required"`
	DbHost     string `env:"DB_HOST,required"`
	DbPort     string `env:"DB_PORT,required"`
	DbName     string `env:"DB_NAME,required"`
	Port       string `env:"COMMENT_SERVICE_PORT,required"`
}
