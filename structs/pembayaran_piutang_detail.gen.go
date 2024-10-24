// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package structs

import (
	"time"
)

const TableNamePembayaranPiutangDetail = "pembayaran_piutang_detail"

// PembayaranPiutangDetail mapped from table <pembayaran_piutang_detail>
type PembayaranPiutangDetail struct {
	ID                  string     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	PembayaranPiutangID string     `gorm:"column:pembayaran_piutang_id;not null" json:"pembayaran_piutang_id"`
	PiutangID           string     `gorm:"column:piutang_id;not null" json:"piutang_id"`
	Nominal             int32     `gorm:"column:nominal;not null" json:"nominal"`
	IsLunas             int16     `gorm:"column:is_lunas;not null;default: 0" json:"is_lunas"`
	DtmCrt              time.Time `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd              time.Time `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	SyncKey             string    `gorm:"column:sync_key;default:now()" json:"sync_key"`
	TeamleaderID        int32     `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	UserID              int32     `gorm:"column:user_id;default:null" json:"user_id"`
}

// TableName PembayaranPiutangDetail's table name
func (*PembayaranPiutangDetail) TableName() string {
	return TableNamePembayaranPiutangDetail
}