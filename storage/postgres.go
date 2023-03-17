package storage

import(
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
type Config struct{
	Host 			string
	Port	string
	Password	string
	User	string
	DBName	string
	SSlMode	string

	}	

func NewConnection(config *Config) (*gorm.DB, error){
	dsn:=fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",config.Host,config.User,config.Password,config.DBName,config.Port,config.SSlMode)
	db,err:=gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil{
		return nil,err
	}
	return db,nil
}