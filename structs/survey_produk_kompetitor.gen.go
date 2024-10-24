// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package structs

import (
	"time"
)

const TableNameSurveyProdukKompetitor = "survey_produk_kompetitor"

// SurveyProdukKompetitor mapped from table <survey_produk_kompetitor>
type SurveyProdukKompetitor struct {
	ID                 string     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID         string     `gorm:"column:customer_id;default:null" json:"customer_id"`
	SalesmanID         int32     `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	MerchandiserID     int32     `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	TeamleaderID       int32     `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	Datetime           time.Time `gorm:"column:datetime;default:null" json:"datetime"`
	ProdukKompetitorID int32     `gorm:"column:produk_kompetitor_id;default:null" json:"produk_kompetitor_id"`
	HargaBeli          int64     `gorm:"column:harga_beli;default:null" json:"harga_beli"`
	HargaJual          int64     `gorm:"column:harga_jual;default:null" json:"harga_jual"`
	CreatedAt          time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;default:now()" json:"updated_at"`
	Note               string    `gorm:"column:note;default:null" json:"note"`
	SyncKey            time.Time `gorm:"column:sync_key;default:now()" json:"sync_key"`
	LatitudeLongitude  string    `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
	Photo              string    `gorm:"column:photo;default:null" json:"photo"`
	KompetitorID       string     `gorm:"column:kompetitor_id;default:null" json:"kompetitor_id"`
}

// TableName SurveyProdukKompetitor's table name
func (*SurveyProdukKompetitor) TableName() string {
	return TableNameSurveyProdukKompetitor
}