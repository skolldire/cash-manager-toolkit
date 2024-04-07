package db_connection

type Config struct {
	DbName             string `json:"db_name"`
	DbUser             string `json:"db_user"`
	DbPassword         string `json:"db_password"`
	DbHost             string `json:"db_host"`
	DbPort             string `json:"db_port"`
	DbDns              string `json:"db_dns"`
	MaxOpenCons        int    `json:"max-open-cons"`
	SetMaxIdleCons     int    `json:"set-max-idle-cons"`
	SetConnMaxLifetime int    `json:"set-conn-max-lifetime"`
}
