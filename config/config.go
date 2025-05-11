package config

import "os"

type DatabaseSettingType struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	DBName       string
	DBPORT       string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

type ConfigType struct {
	DatabaseSetting DatabaseSettingType
	Port            string
}

func GetDbSetting() *DatabaseSettingType {
	return &DatabaseSettingType{
		DBType:       "postgres",
		UserName:     os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		Host:         os.Getenv("DB_HOST"),
		DBName:       os.Getenv("DB_NAME"),
		DBPORT:       os.Getenv("DB_PORT"),
		Charset:      "utf8",
		ParseTime:    true,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	}
}

func GetConfig() *ConfigType {
	return &ConfigType{
		DatabaseSetting: *GetDbSetting(),
		Port:            os.Getenv("PORT"),
	}
}
