package structs

const TableNameFrmCategory = "frm.category"

// FrmCategory mapped from table <frm.category>
type FrmCategory struct {
	ID          string `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name        string `gorm:"column:name;not null" json:"name"`
	Description string `gorm:"column:description" json:"description"`
	Icon        string `gorm:"column:icon" json:"icon"`
}

// TableName FrmCategory's table name
func (*FrmCategory) TableName() string {
	return TableNameFrmCategory
}
