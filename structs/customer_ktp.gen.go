// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package structs

import (
	"time"
)

const TableNameCustomerKtp = "customer_ktp"

// CustomerKtp mapped from table <customer_ktp>
type CustomerKtp struct {
	ID             string     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID     string     `gorm:"column:customer_id;default:null" json:"customer_id"`
	ImageKtp       string    `gorm:"column:image_ktp;default:null" json:"image_ktp"`
	CreatedAt      time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;default:now()" json:"updated_at"`
	SalesmanID     int32     `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	TeamleaderID   int32     `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	MerchandiserID int32     `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	SyncKey        string    `gorm:"column:sync_key;default:now()" json:"sync_key"`
	NameKtp        string    `gorm:"column:name_ktp;default:null" json:"name_ktp"`
	NikKtp         string    `gorm:"column:nik_ktp;default:null" json:"nik_ktp"`
}

// TableName CustomerKtp's table name
func (*CustomerKtp) TableName() string {
	return TableNameCustomerKtp
}