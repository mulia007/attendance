package models

type Employee struct {
	Id_employee   int    `gorm:"primary_key;type:int" json:"id_employee"`
	Name_employee string `gorm:"type:nvarchar(200)" json:"name_employee"`
	Photo         string `gorm:"type:nvarchar(200)" json:"photo"`
	Islogin       int8   `gorm:"type:tinyint" json:"islogin"`
}
