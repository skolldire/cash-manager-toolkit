package db_connection

type Config struct {
	DbDriver        string `json:"db_driver"`
	DbName          string `json:"db_name"`
	DbUser          string `json:"db_user"`
	DbPassword      string `json:"db_password"`
	DbHost          string `json:"db_host"`
	DbPort          string `json:"db_port"`
	DbDns           string `json:"db_dns"`
	MaxOpenCons     int    `json:"max-open-cons"`
	MaxIdleCons     int    `json:"max-idle-cons"`
	ConnMaxLifetime uint   `json:"conn-max-lifetime"`
}
