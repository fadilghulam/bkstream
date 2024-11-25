package structs

import "time"

const TableNameFrmComment = "frm.comment"

// FrmComment mapped from table <frm.comment>
type FrmComment struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID         int32     `gorm:"column:user_id;not null" json:"userId"`
	ThreadID       int64     `gorm:"column:thread_id;not null" json:"threadId"`
	CommentQuoteID int64     `gorm:"column:comment_quote_id;default:null" json:"commentQuoteId"`
	Body           string    `gorm:"column:body;not null" json:"body"`
	IsPinned       int16     `gorm:"column:is_pinned;default: null" json:"isPinned"`
	PinnedAt       time.Time `gorm:"column:pinned_at;default:null" json:"pinnedAt"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:now()" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;default:now()" json:"updatedAt"`
}

// TableName FrmComment's table name
func (*FrmComment) TableName() string {
	return TableNameFrmComment
}
