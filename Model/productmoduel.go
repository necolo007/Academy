package Model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `gorm:"varchar(50)"`
	Description string  `gorm:"type:text"`
	Picture     string  `gorm:"varchar(225)"`
	Price       float64 `gorm:"type:decimal(10,2)"`
	Categories  string  `gorm:"type:text"` //json字符串
}
