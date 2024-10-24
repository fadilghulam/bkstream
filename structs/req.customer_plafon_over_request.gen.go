// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package structs

import (
	"time"
)

const TableNameCustomerPlafonOverRequest = "customer_plafon_over_request"

// CustomerPlafonOverRequest mapped from table <customer_plafon_over_request>
type CustomerPlafonOverRequest struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	RequestAt   time.Time `gorm:"column:request_at;not null;default:now()" json:"request_at"`
	RequestedID int64     `gorm:"column:requested_id;default:null" json:"requested_id"`
	IsApprove   int16     `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt   time.Time `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID   int64     `gorm:"column:approve_id;default:null" json:"approve_id"`
	Note        string    `gorm:"column:note;default:null" json:"note"`
	Attachment  string    `gorm:"column:attachment;default:null" json:"attachment"`
	CreatedAt   time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:now()" json:"updated_at"`
	DateStart   time.Time `gorm:"column:date_start;not null;default:now()" json:"date_start"`
	DateEnd     time.Time `gorm:"column:date_end;not null;default:now()" json:"date_end"`
	ExecutedAt  time.Time `gorm:"column:executed_at;default:null" json:"executed_at"`
	UserID 		int32 	  `gorm:"column:user_id;not null" json:"userId"`
}

// TableName CustomerPlafonOverRequest's table name
func (*CustomerPlafonOverRequest) TableName() string {
	return TableNameCustomerPlafonOverRequest
}