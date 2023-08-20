package models

type Lockon struct {
	Id             int    `gorm:"primary_key;type:int" json:"id"`
	Datetime_logon string `gorm:"type:nvarchar(200)" json:"datetime_logon"`
	Image          string `gorm:"type:nvarchar(200)" json:"image"`
	Id_employee    int    `gorm:"type:int(11)" json:"id_employee"`
}
