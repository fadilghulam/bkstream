package structs

const TableNameWebSocketState = "websocket_state"

// WebSocketState mapped from table <websocket_state>
type WebSocketState struct {
	ID        int64   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	AppName   string  `gorm:"column:app_name;not null" json:"appName"`
	UserID    int32   `gorm:"column:user_id;not null" json:"userId"`
	Route     string  `gorm:"column:route;not null" json:"route"`
	State     string  `gorm:"column:state;not null" json:"state"`
	Url       *string `gorm:"column:url;default:null" json:"url"`
	Timestamp string  `gorm:"column:timestamp;default: now()" json:"timestamp"`
	Body      *string `gorm:"column:body;default:null" json:"body"`
	Error     *string `gorm:"column:error;default:null" json:"error"`
}

// TableName WebSocketState's table name
func (*WebSocketState) TableName() string {
	return TableNameWebSocketState
}
