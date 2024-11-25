package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"gorm.io/gorm/clause"

	db "bkstream/config"
	"bkstream/controllers"

	// mobile "bkstream/controllers/mobile"
	"bkstream/helpers"
	"bkstream/structs"

	"github.com/zishang520/socket.io/v2/socket"
)

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Landing Page!")
	})

	app.Get("/testRedis", controllers.GetDataCustomerRedis)
	app.Get("getDashboardOmzet", controllers.GetDashboardOmzet)

	// 	app.Post("login", controllers.Login)
	// 	app.Post("sendOtp", controllers.SendOtp)

	// 	app.Get("/getUtcTime", func(c *fiber.Ctx) error {
	// 		return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 			"message": "success",
	// 			"data":    time.Now().UTC(),
	// 		})
	// 	})

	// 	app.Get("/cronGenerateUserId", controllers.GenerateTransactionsUserId)
	// 	app.Get("/testCronGenerateUserId", controllers.TestGenerateUserId)
	// 	app.Get("/cronGenerateUserLog", controllers.GenerateUserLog)

	// 	app.Get("/generateFlag", controllers.GenerateFlag)
	// 	app.Get("/getData", controllers.GetData)
	// 	app.Get("/getDataToday", controllers.GetDataToday)

	// 	app.Post("/insertTransactions", mobile.InsertTransactions)

	// 	serviceRoute := app.Group("service")
	// 	serviceRoute.Post("doUpload", controllers.DoUpload)

	// 	officeRoute := app.Group("office")
	// 	// officeRoute.Use(AuthMiddleware)

	// 	officeRoute.Get("/getProductTrends", controllers.GetProductTrends)
	// 	officeRoute.Get("/TestQuery", controllers.TestQuery)
	// 	officeRoute.Get("/getSalesmanDaily", controllers.GetSalesmanDailySales)
	// 	officeRoute.Get("/getUserBranch", controllers.GetUserBranch)

	// 	mobileRoute := app.Group("pluto-mobile")
	// 	mobileRoute.Get("/", func(c *fiber.Ctx) error {
	// 		return c.SendString("Landing Page Pluto Mobile!")
	// 	})

	// 	mobileRoute.Get("getAppVersioning", mobile.GetAppVersioning)

	// 	mobileRoute.Get("getGudang", mobile.GetGudang)
	// 	mobileRoute.Get("getProdukGudang", mobile.GetProdukByGudang)
	// 	mobileRoute.Get("getItemGudang", mobile.GetItemByGudang)
	// 	mobileRoute.Post("confirmOrder", mobile.ConfirmOrder)

	// 	mobileRoute.Get("getListPengajuan", mobile.GetDataRequests)

	// 	//sales
	// 	mobileRoute.Get("getStokProduk", mobile.GetStokProduk)
	// 	mobileRoute.Get("getListOrder", mobile.GetListOrder)
	// 	mobileRoute.Post("postOrder", mobile.PostOrder)

	// 	//md
	// 	mobileRoute.Get("getStokItem", mobile.GetStokItem)
	// 	mobileRoute.Get("getListOrderItem", mobile.GetListOrderMD)
	// 	mobileRoute.Post("postOrderItem", mobile.PostOrderMD)

	// 	mobileRoute.Use(AuthMiddleware)
	// 	mobileRoute.Post("getRefreshUser", controllers.RefreshDataUser)
}

type NotificationPayload struct {
	UserId        *string `json:"user_id"`
	ReferenceName *string `json:"reference_name"`
	ReferenceId   *string `json:"reference_id"`
	AppName       *string `json:"app_name"`
}

type UpdateStatePayload struct {
	UserId    *int32  `json:"userId"`
	AppName   *string `json:"appName"`
	State     *string `json:"state"`
	Url       *string `json:"url"`
	Body      *string `json:"body"`
	Error     *string `json:"error"`
	Timestamp *string `json:"timestamp"`
	Route     *string `json:"route"`
}

func SocketIoSetup(app *fiber.App) {

	socketio := socket.NewServer(nil, nil)

	// socketio.On("connection", func(clients ...interface{}) {

	// socketio.Of("/", nil).On("connection", func(clients ...interface{}) {

	// })

	//send to specific user
	app.Get("/send/:id", func(c *fiber.Ctx) error {
		clientID := c.Params("id")

		socketio.To(socket.Room(clientID)).Emit("returnData", "this is return data to specific client")

		return c.SendString("Message sent to client: " + clientID)
	})

	app.Get("/checkRoom/:room", func(c *fiber.Ctx) error {
		room := c.Params("room")

		temp := 0
		socketio.In(socket.Room(room)).FetchSockets()(func(sockets []*socket.RemoteSocket, err error) {
			if err != nil {
				// handle error
				fmt.Println("Error fetching sockets:", err)
			}

			// for _, _ := range sockets {
			// 	temp++
			// }
			temp = len(sockets)
		})

		socketio.To(socket.Room(room)).Emit("returnData", "this is return data for room")

		return c.SendString("Room checked for: " + room + " sum sockets : " + strconv.Itoa(temp))
	})

	app.Get("/wsFind", func(c *fiber.Ctx) error {
		userId := c.Query("userId")
		appName := c.Query("appName")

		wsLog := new([]structs.WebSocketLog)

		// if err := db.DB.First(&wsLog, "user_id = ? AND app_name = ?", userId, appName).Order("datetime DESC").Error; err != nil {
		if err := db.DB.Where("user_id = ? AND app_name = ?", userId, appName).Find(&wsLog).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Success: false,
				Message: "Gagal mendapatkan data user",
			})
		}

		if wsLog == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Success: false,
				Message: "User not found",
			})
		}

		data := map[string]interface{}{
			"event": "find_my",
		}

		// Emit the map as JSON to the client
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshalling data:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Success: false,
				Message: "Server error",
			})
		}

		for i := 0; i < len(*wsLog); i++ {
			socketio.To(socket.Room((*wsLog)[i].SocketID)).Emit("doEvent", string(jsonData))
		}

		// socketio.To(socket.Room(wsLog.SocketID)).Emit("doEvent", string(jsonData))

		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Success: true,
			Message: "Event for find_my sent successfully",
		})
	})

	app.Get("/wsDataCenter", func(c *fiber.Ctx) error {
		userId := c.Query("userId")
		appName := c.Query("appName")

		wsLog := new(structs.WebSocketLog)

		if err := db.DB.First(&wsLog, "user_id = ? AND app_name = ?", userId, appName).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Success: false,
				Message: "Gagal mendapatkan data user",
			})
		}

		if wsLog == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Success: false,
				Message: "User not found",
			})
		}

		data := map[string]interface{}{
			"event": "data_center",
		}

		// Emit the map as JSON to the client
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshalling data:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Success: false,
				Message: "Server error",
			})
		}

		socketio.To(socket.Room(wsLog.SocketID)).Emit("doEvent", string(jsonData))

		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Success: true,
			Message: "Event for data_center sent successfully",
		})
	})

	app.Get("/wsLogout", func(c *fiber.Ctx) error {
		userId := c.Query("userId")
		appName := c.Query("appName")

		wsLog := new(structs.WebSocketLog)

		if err := db.DB.First(&wsLog, "user_id = ? AND app_name = ?", userId, appName).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Success: false,
				Message: "Gagal mendapatkan data user",
			})
		}

		if wsLog == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Success: false,
				Message: "User not found",
			})
		}

		data := map[string]interface{}{
			"event": "logout",
		}

		// Emit the map as JSON to the client
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshalling data:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Success: false,
				Message: "Server error",
			})
		}

		socketio.To(socket.Room(wsLog.SocketID)).Emit("doEvent", string(jsonData))
		// if err := socketio.To(socket.Room(wsLog.SocketID)).Emit("doEvent", string(jsonData)); err != nil {
		// 	fmt.Println("Error emitting event:", err)
		// } else {
		// 	fmt.Println("Event sent successfully")
		// }

		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Success: true,
			Message: "Event for logout sent successfully",
		})
	})

	_, err := db.DBPGX.Exec(context.Background(), "LISTEN onapproved")
	if err != nil {
		fmt.Println("Error setting up listener: %v\n", err.Error())
	}

	_, err = db.DBPGX.Exec(context.Background(), "LISTEN ontransaction")
	if err != nil {
		fmt.Println("Error setting up listener: %v\n", err.Error())
	}

	go func() {

		for {
			notification, err := db.DBPGX.WaitForNotification(context.Background())
			if err != nil {
				fmt.Printf("Error waiting for notification: %v\n", err)
			}

			if notification.Channel == "onapproved" {

				var payload NotificationPayload
				err = json.Unmarshal([]byte(notification.Payload), &payload)
				if err != nil {
					fmt.Printf("Error unmarshalling JSON: %v\n", err)
				}

				wsLog := new(structs.WebSocketLog)

				if payload.UserId != nil && payload.AppName != nil {
					if err := db.DB.First(&wsLog, "user_id = ? AND app_name = ?", *payload.UserId, *payload.AppName).Error; err != nil {
						fmt.Println("Error getting user data " + err.Error())
					}

					if wsLog != nil {
						data := map[string]interface{}{
							"event":    "approval",
							"ref_id":   *payload.ReferenceId,
							"ref_name": *payload.ReferenceName,
						}

						// Emit the map as JSON to the client
						jsonData, err := json.Marshal(data)
						if err != nil {
							fmt.Println("Error marshalling data:", err)
						}

						socketio.To(socket.Room(wsLog.SocketID)).Emit("doEvent", string(jsonData))
					}
				}

			}
			// fmt.Printf("Received notification on channel %s: %s\n", notification.Channel, notification.Payload)

			// var payload NotificationPayload
			// err = json.Unmarshal([]byte(notification.Payload), &payload)
			// if err != nil {
			// 	fmt.Printf("Error unmarshalling JSON: %v\n", err)
			// }

			// fmt.Println(&payload)

			// Emit notification to all connected Socket.IO clients
			// socket.BroadcastToNamespace("/", "notification", notification.Payload)
			// socketio.Local().Emit("notification", "test message from pg notify with data : "+*payload.ReferenceName+" = "+*payload.ReferenceId)
		}
	}()

	go func() {
		ticker := time.NewTicker(10 * time.Second) // Create a ticker for 10-second intervals
		defer ticker.Stop()                        // Ensure the ticker stops when the function exits

		for {
			select {
			case <-ticker.C: // Every 10 seconds
				fmt.Println("emit")
				datas, _ := controllers.FetchDashboardOmzet("", "", "", "a")
				socketio.To("office").Emit("upDashboard", datas)
				// Add your task logic here
			}
		}
	}()

	socketio.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)

		fmt.Println("Client connected: ", client.Id())

		userId, _ := client.Request().Query().Get("userId")
		appName, _ := client.Request().Query().Get("appName")
		appVersion, _ := client.Request().Query().Get("appVersion")
		deviceId, _ := client.Request().Query().Get("deviceId")
		rooms, _ := client.Request().Query().Get("room")

		tempRooms := strings.Split(rooms, ",")

		webSocketLog := new(structs.WebSocketLog)

		webSocketLog.UserID = userId
		webSocketLog.SocketID = fmt.Sprintf("%s", client.Id())
		webSocketLog.DeviceID = deviceId
		webSocketLog.AppName = appName
		webSocketLog.AppVersion = appVersion
		webSocketLog.Datetime = time.Now()

		tx := db.DB.Begin()

		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "app_name"}, {Name: "device_id"}, {Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"socket_id", "app_version", "datetime"}),
		}).Create(&webSocketLog).Error; err != nil {
			tx.Rollback()
			fmt.Println("Error creating WebSocketLog: ", err)
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			fmt.Println("Error commit WebSocketLog: ", err)
		}

		client.Join(socket.Room(strings.ToLower(appName)))

		for _, vRoom := range tempRooms {
			// socketio.To(socket.Room(room)).Emit("newUser", userId)
			client.Join(socket.Room(vRoom))
		}

		//all connected clients
		// allClients2 := socketio.Of("/", nil).Sockets()
		// fmt.Println("all clients 2 : ", allClients2.Len())

		//broadcast
		// socketio.Local().Emit("test", "test message from "+client.Id())

		// client.Join(socket.Room("pluto"))
		// client.Join(socket.Room("malang"))

		// client.On("/pluto/requestData", func(args ...interface{}) {
		// 	fmt.Println(args)
		// 	fmt.Println(args[0])
		// 	fmt.Println(args[1])
		// })

		client.On("setState", func(args ...interface{}) {
			var payload UpdateStatePayload
			// err = json.Unmarshal([]byte(args[0].(string)), &payload)
			switch args[0].(type) {
			case string:
				err = json.Unmarshal([]byte(args[0].(string)), &payload)
				if err != nil {
					fmt.Println("Error unmarshalling JSON: %v\n", err)
					socketio.To(socket.Room(client.Id())).Emit("setStateResult", helpers.ResponseWebSocket{
						Success: false,
						Message: "Terjadi kesalahan mengambil input data",
						Data:    nil,
						Error:   err.Error(),
					})
				}
			case map[string]interface{}:
				data, err := json.Marshal(args[0].(map[string]interface{}))
				if err != nil {
					fmt.Println("Error marshalling data: %v\n", err)
					socketio.To(socket.Room(client.Id())).Emit("setStateResult", helpers.ResponseWebSocket{
						Success: false,
						Message: "Terjadi kesalahan mengambil input data",
						Data:    nil,
						Error:   err.Error(),
					})
				}
				err = json.Unmarshal(data, &payload)
				if err != nil {
					fmt.Println("Error unmarshalling JSON: %v\n", err)
					socketio.To(socket.Room(client.Id())).Emit("setStateResult", helpers.ResponseWebSocket{
						Success: false,
						Message: "Terjadi kesalahan mengambil input data",
						Data:    nil,
						Error:   err.Error(),
					})
				}
			}
			// data, err := json.Marshal(args[0].(map[string]interface{}))
			// if err != nil {
			// 	fmt.Println("Error marshalling data: %v\n", err)
			// 	return
			// }
			// err = json.Unmarshal(data, &payload)
			// if err != nil {
			// 	fmt.Println("Error unmarshalling JSON: %v\n", err)
			// }

			wsState := new(structs.WebSocketState)
			wsState.UserID = *payload.UserId
			wsState.AppName = *payload.AppName
			wsState.Route = *payload.Route
			wsState.State = *payload.State
			if payload.Url != nil {
				wsState.Url = payload.Url
			} else {
				wsState.Url = nil
			}

			if payload.Body != nil {
				wsState.Body = payload.Body
			} else {
				wsState.Body = nil
			}

			if payload.Error != nil {
				wsState.Error = payload.Error
			} else {
				wsState.Error = nil
			}
			wsState.Timestamp = *payload.Timestamp

			tx := db.DB.Begin()
			if err := tx.Create(&wsState).Error; err != nil {
				tx.Rollback()
				fmt.Println("Error creating WebSocketState: ", err)
				socketio.To(socket.Room(client.Id())).Emit("setStateResult", helpers.ResponseWebSocket{
					Success: false,
					Message: "Terjadi kesalahan input data",
					Data:    nil,
					Error:   err.Error(),
				})

			}

			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
				fmt.Println("Error commit WebSocketState: ", err)
				socketio.To(socket.Room(client.Id())).Emit("setStateResult", helpers.ResponseWebSocket{
					Success: false,
					Message: "Terjadi kesalahan simpan data",
					Data:    nil,
					Error:   err.Error(),
				})
			}

			// socketio.To(socket.Room(client.Id())).Emit("setStateResult", wsState)
			socketio.To(socket.Room(client.Id())).Emit("setStateResult", helpers.ResponseWebSocket{
				Success: true,
				Message: "Data berhasil disimpan",
				Data:    wsState,
				Error:   "null",
			})
			// err = socketio.To(socket.Room(client.Id())).Emit("updateState", wsState)
			// if err != nil {
			// 	fmt.Println("Error while emitting to room:", err)
			// } else {
			// 	fmt.Println("Successfully emitted to room:", "test")
			// }
		})

		client.On("setTransaction", func(args ...interface{}) {

			type TemplateInputUser struct {
				Data      map[string]interface{} `json:"data"`
				DeletedID map[string]interface{} `json:"deletedId"`
			}

			var inputUser TemplateInputUser

			switch args[0].(type) {
			case string:
				err = json.Unmarshal([]byte(args[0].(string)), &inputUser)
				if err != nil {
					fmt.Println("Error unmarshalling JSON: %v\n", err)
					socketio.To(socket.Room(client.Id())).Emit("setTransactionResult", helpers.ResponseWebSocket{
						Success: false,
						Message: "Terjadi kesalahan mengambil input data",
						Data:    nil,
						Error:   err.Error(),
					})
				}
			case map[string]interface{}:
				data, err := json.Marshal(args[0].(map[string]interface{}))
				if err != nil {
					fmt.Println("Error marshalling data: %v\n", err)
					socketio.To(socket.Room(client.Id())).Emit("setTransactionResult", helpers.ResponseWebSocket{
						Success: false,
						Message: "Terjadi kesalahan mengambil input data",
						Data:    nil,
						Error:   err.Error(),
					})
					return
				}
				err = json.Unmarshal(data, &inputUser)
				if err != nil {
					fmt.Println("Error unmarshalling JSON: %v\n", err)
					socketio.To(socket.Room(client.Id())).Emit("setTransactionResult", helpers.ResponseWebSocket{
						Success: false,
						Message: "Terjadi kesalahan mengambil input data",
						Data:    nil,
						Error:   err.Error(),
					})
				}
			}

			result := make(map[string][]map[string]interface{})

			tx := db.DB.Begin()

			for tableName, records := range inputUser.DeletedID {
				instanceSliceDelete, err := structs.GetStructInstanceByTableName(tableName)
				if err != nil {
					tx.Rollback()
					fmt.Println("Gagal mendapatkan tabel data, " + err.Error())
					socketio.To(socket.Room(client.Id())).Emit("setTransactionResult", helpers.ResponseWebSocket{
						Success: false,
						Message: "Gagal mendapatkan tabel data delete",
						Data:    nil,
						Error:   err.Error(),
					})
				}

				whereIdIn := strings.Split(records.(string), ",")

				if err := tx.Clauses(clause.Returning{}).Where("id IN (?)", whereIdIn).Delete(instanceSliceDelete).Error; err != nil {
					tx.Rollback()
					fmt.Println("Gagal delete data " + err.Error())
					socketio.To(socket.Room(client.Id())).Emit("setTransactionResult", helpers.ResponseWebSocket{
						Success: false,
						Message: "Gagal delete data",
						Data:    nil,
						Error:   err.Error(),
					})
				}

				recordsValue := reflect.ValueOf(instanceSliceDelete).Elem() // dereference the pointer to slice
				for i := 0; i < recordsValue.Len(); i++ {
					record := recordsValue.Index(i).Interface() // access the individual record

					// Use reflection to get id, sync_key, and created_at fields from the record
					id := reflect.ValueOf(record).FieldByName("ID").Interface()
					createdAtField := reflect.ValueOf(record).FieldByName("CreatedAt")
					dtmCrtField := reflect.ValueOf(record).FieldByName("DtmCrt")
					syncKeyField := reflect.ValueOf(record).FieldByName("SyncKey")
					var syncKey interface{}

					if createdAtField.IsValid() {
						syncKey = createdAtField.Interface()
					}

					if dtmCrtField.IsValid() {
						syncKey = dtmCrtField.Interface()
					}

					if syncKeyField.IsValid() {
						syncKey = syncKeyField.Interface()
					}

					result[tableName] = append(result[tableName], map[string]interface{}{
						"id":       id,
						"sync_key": syncKey,
					})
				}
			}

			tx.Commit()

			tx = db.DB.Begin()
			for tableName, records := range inputUser.Data {

				instanceSlice, err := structs.GetStructInstanceByTableName(tableName)
				if err != nil {
					tx.Rollback()
					fmt.Println("Gagal mendapatkan tabel data " + err.Error())
					socketio.To(socket.Room(client.Id())).Emit("setTransactionResult", helpers.ResponseWebSocket{
						Success: false,
						Message: "Gagal mendapatkan tabel data insert",
						Data:    nil,
						Error:   err.Error(),
					})
				}

				recordsBytes, err := json.Marshal(records)
				if err != nil {
					tx.Rollback()
					fmt.Println("Gagal konversi data tabel " + err.Error())
					socketio.To(socket.Room(client.Id())).Emit("setTransactionResult", helpers.ResponseWebSocket{
						Success: false,
						Message: "Gagal memproses data insert",
						Data:    nil,
						Error:   err.Error(),
					})
				}

				if err := json.Unmarshal(recordsBytes, instanceSlice); err != nil {
					tx.Rollback()
					// return c.Status(fiber.StatusBadRequest).SendString("Failed to parse records: " + err.Error())
					fmt.Println("Gagal konversi data tabel" + err.Error())
					socketio.To(socket.Room(client.Id())).Emit("setTransactionResult", helpers.ResponseWebSocket{
						Success: false,
						Message: "Gagal memproses data insert",
						Data:    nil,
						Error:   err.Error(),
					})
				}

				var tempIds []string
				recordsValue := reflect.ValueOf(instanceSlice).Elem() // dereference the pointer to slice
				for i := 0; i < recordsValue.Len(); i++ {
					record := recordsValue.Index(i).Interface() // access the individual record

					// Use reflection to get id, sync_key, and created_at fields from the record
					id := reflect.ValueOf(record).FieldByName("ID").Interface()

					tempIds = append(tempIds, fmt.Sprintf("%v", id))
				}

				if err := tx.Clauses(clause.Returning{}).Save(instanceSlice).Error; err != nil {
					tx.Rollback()
					fmt.Println("Gagal insert data" + err.Error())
					socketio.To(socket.Room(client.Id())).Emit("setTransactionResult", helpers.ResponseWebSocket{
						Success: false,
						Message: "Gagal menyimpan data",
						Data:    nil,
						Error:   err.Error(),
					})
				}

				tx.Where("id IN (?)", tempIds).Find(instanceSlice)

				recordsValue = reflect.ValueOf(instanceSlice).Elem() // dereference the pointer to slice
				for i := 0; i < recordsValue.Len(); i++ {
					record := recordsValue.Index(i).Interface() // access the individual record

					// Use reflection to get id, sync_key, and created_at fields from the record
					id := reflect.ValueOf(record).FieldByName("ID").Interface()
					createdAtField := reflect.ValueOf(record).FieldByName("CreatedAt")
					dtmCrtField := reflect.ValueOf(record).FieldByName("DtmCrt")
					syncKeyField := reflect.ValueOf(record).FieldByName("SyncKey")
					var syncKey interface{}

					if createdAtField.IsValid() {
						syncKey = createdAtField.Interface()
					}

					if dtmCrtField.IsValid() {
						syncKey = dtmCrtField.Interface()
					}

					if syncKeyField.IsValid() {
						syncKey = syncKeyField.Interface()
					}

					result[tableName] = append(result[tableName], map[string]interface{}{
						"id":       id,
						"sync_key": syncKey,
					})
				}
			}
			tx.Commit()

			// return c.Status(fiber.StatusOK).JSON(fiber.Map{
			// 	"message": "success",
			// 	"data":    result,
			// })

			// socketio.To(socket.Room(client.Id())).Emit("setTransactionResult ", result)
			socketio.To(socket.Room(client.Id())).Emit("setTransactionResult", helpers.ResponseWebSocket{
				Success: true,
				Message: "Data berhasil disimpan",
				Data:    result,
				Error:   "null",
			})

		})

		client.On("disconnect", func(args ...interface{}) {
			// fmt.Println("Client disconnected: ", client.Id())

			db.DB.Where("client_id = ?", client.Id()).Delete(&structs.WebSocketLog{})

			client.Disconnect(true)

			// allClients2 := socketio.Of("/", nil).Sockets()
			// fmt.Println("all clients 2 : ", allClients2.Len())
		})
	})

	app.Get("/socket.io", adaptor.HTTPHandler(socketio.ServeHandler(nil)))
	app.Post("/socket.io", adaptor.HTTPHandler(socketio.ServeHandler(nil)))
}
