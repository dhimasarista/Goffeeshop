package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySQLConfig menyimpan konfigurasi database
type MySQLConfig struct {
	Username  string
	Password  string
	Host      string
	Port      string
	Name      string
	ParseTime bool
}

// generateDSN membangun string DSN untuk koneksi MySQL
func generateDSN(config MySQLConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.ParseTime,
	)
}

// membaca konfigurasi database dari environment variables menggunakan viper
func readMySQLConfig() MySQLConfig {
	viper.SetConfigFile(".env") // pastikan ini membaca file .env
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading .env file: %s", err)
	}

	var config MySQLConfig

	// Membaca konfigurasi dari environment variables
	config.Username = viper.GetString("DB_USERNAME")
	config.Password = viper.GetString("DB_PASSWORD")
	config.Host = viper.GetString("DB_HOST")
	config.Port = viper.GetString("DB_PORT")
	config.Name = viper.GetString("DB_NAME")
	config.ParseTime = viper.GetBool("DB_PARSE_TIME")

	return config
}

var mysqlConfig MySQLConfig = readMySQLConfig()

// GormDB membuat koneksi *gorm.DB dan mengembalikkannya
func GormDB() *gorm.DB {
	// Membuat koneksi ke DB menggunakan konfigurasi yang dibaca
	db, err := gorm.Open(mysql.Open(generateDSN(mysqlConfig)), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return nil
	}
	return db
}
