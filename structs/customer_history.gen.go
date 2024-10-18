// Code generated by gorm.io/gen;default:null. DO NOT EDIT.
// Code generated by gorm.io/gen;default:null. DO NOT EDIT.
// Code generated by gorm.io/gen;default:null. DO NOT EDIT.

package structs

import (
	"time"
)

const TableNameCustomerHistory = "customer_history"

// CustomerHistory mapped from table <customer_history>
type CustomerHistory struct {
	ID                string     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	EmployeeID        int64     `gorm:"column:employee_id;default:null" json:"employee_id"`
	CustomerID        string     `gorm:"column:customer_id;default:null" json:"customer_id"`
	Datetime          time.Time `gorm:"column:datetime;default:null" json:"datetime"`
	JSONBefore        string    `gorm:"column:json_before;default:null" json:"json_before"`
	SrID              int16     `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID           int16     `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	BranchID          int16     `gorm:"column:branch_id;default:null" json:"branch_id"`
	AreaID            int32     `gorm:"column:area_id;default:null" json:"area_id"`
	CreatedAt         time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;default:now()" json:"updated_at"`
	JSONAfter         string    `gorm:"column:json_after;default:null" json:"json_after"`
	SalesmanID        int32     `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	Note              string    `gorm:"column:note;default:null" json:"note"`
	LatitudeLongitude string    `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
}

// TableName CustomerHistory's table name
func (*CustomerHistory) TableName() string {
	return TableNameCustomerHistory
}
