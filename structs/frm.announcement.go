package structs

import (
	"time"
)

const TableNameFrmAnnouncement = "frm.accouncement"

// FrmAnnouncement mapped from table <frm.accouncement>
type FrmAnnouncement struct {
	ID         string    `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Title      string    `gorm:"column:title;not null" json:"title"`
	Body       string    `gorm:"column:body;not null" json:"body"`
	CategoryID int32     `gorm:"column:category_id" json:"category_id"`
	ThreadID   int32     `gorm:"column:thread_id" json:"thread_id"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

// TableName FrmAnnouncement's table name
func (*FrmAnnouncement) TableName() string {
	return TableNameFrmAnnouncement
}
