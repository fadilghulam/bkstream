package structs

import "time"

const TableNameFrmThread = "frm.thread"

// FrmThread mapped from table <frm.thread>
type FrmThread struct {
	ID                int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID            int32     `gorm:"column:user_id;not null" json:"userId"`
	CategoryID        int32     `gorm:"column:category_id;not null" json:"categoryId"`
	Title             string    `gorm:"column:title;not null" json:"title"`
	Body              string    `gorm:"column:body;not null" json:"body"`
	IsPinned          int16     `gorm:"column:is_pinned;default: null" json:"isPinned"`
	PinnedAt          time.Time `gorm:"column:pinned_at;default:null" json:"pinnedAt"`
	IsDuplicate       int16     `gorm:"column:is_duplicate;default: null" json:"isDuplicate"`
	DuplicateThreadID int64     `gorm:"column:duplicate_thread_id;default:null" json:"duplicateThreadId"`
	IsClose           int16     `gorm:"column:is_close;default: null" json:"isClose"`
	CloseAt           time.Time `gorm:"column:close_at;default:null" json:"closeAt"`
	Tags              string    `gorm:"column:tags;default:null" json:"tags"`
	TotalView         int64     `gorm:"column:total_view;default:0" json:"totalView"`
	CreatedAt         time.Time `gorm:"column:created_at;not null;default:now()" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null;default:now()" json:"updatedAt"`
}

// TableName FrmThread's table name
func (*FrmThread) TableName() string {
	return TableNameFrmThread
}
