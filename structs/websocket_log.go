package structs

import "time"

const TableNameWebSocketLog = "websocket_log"

// WebSocketLog mapped from table <websocket_log>
type WebSocketLog struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID     string    `gorm:"column:user_id;not null" json:"userId"`
	SocketID   string    `gorm:"column:socket_id;not null" json:"socketId"`
	DeviceID   string    `gorm:"column:device_id;not null" json:"deviceId"`
	AppName    string    `gorm:"column:app_name;not null" json:"appName"`
	AppVersion string    `gorm:"column:app_version;default:null" json:"appVersion"`
	Datetime   time.Time `gorm:"column:datetime;default: now()" json:"datetime"`
}

// TableName WebSocketLog's table name
func (*WebSocketLog) TableName() string {
	return TableNameWebSocketLog
}
