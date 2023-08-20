package models

type Lockout struct {
	Id               int    `gorm:"primary_key;type:int" json:"id"`
	Datetime_lockout string `gorm:"type:nvarchar(200)" json:"datetime_lockout"`
	Image            string `gorm:"type:nvarchar(200)" json:"image"`
	Id_employee      int    `gorm:"type:int(11)" json:"id_employee"`
}
