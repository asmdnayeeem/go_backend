package models

import "gorm.io/gorm"

type User struct {
	ID 	 uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func MigrateUser(db *gorm.DB) error{
	err:=db.AutoMigrate(&User{})
	return err
}