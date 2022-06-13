package model

import "gorm.io/gorm"

type Curd struct {
	Id    uint   `gorm:"primarykey"`
	Name  string `gorm:"name"`
	Email string `gorm:"email"`
	Phone string `gorm:"phone"`
	City  string `gorm:"city"`
	gorm.Model
}

func (Curd) TableName() string {
	return "curd"
}