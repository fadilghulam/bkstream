package helpers

type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
type ResponseDataMultiple struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"datas"`
}
type ResponseWithoutData struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
type ResponseWebSocket struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}
