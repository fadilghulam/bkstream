package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"gorm.io/gorm/clause"

	db "bkstream/config"
	"bkstream/structs"

	// "bkstream/controllers"
	// mobile "bkstream/controllers/mobile"

	"bkstream/helpers"

	"github.com/zishang520/socket.io/v2/socket"
)

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Landing Page!")
	})

	// 	app.Post("login", controllers.Login)
	// 	app.Post("sendOtp", controllers.SendOtp)

	app.Get("/getUtcTime", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "success",
			"data":    time.Now().UTC(),
		})
	})
}

type NotificationPayload struct {
	UserId        *string `json:"user_id"`
	ReferenceName *string `json:"reference_name"`
	ReferenceId   *string `json:"reference_id"`
	AppName       *string `json:"app_name"`
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
			"event": "find_my",
		}

		// Emit the map as JSON to the client
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println("Error marshalling data:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Success: false,
				Message: "Server error",
			})
		}

		socketio.To(socket.Room(wsLog.SocketID)).Emit("doEvent", string(jsonData))

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
			log.Println("Error marshalling data:", err)
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
			log.Println("Error marshalling data:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Success: false,
				Message: "Server error",
			})
		}

		socketio.To(socket.Room(wsLog.SocketID)).Emit("doEvent", string(jsonData))

		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Success: true,
			Message: "Event for logout sent successfully",
		})
	})

	_, err := db.DBPGX.Exec(context.Background(), "LISTEN onapproved")
	if err != nil {
		fmt.Println("Error setting up listener: %v\n", err.Error())
	}

	go func() {

		for {
			notification, err := db.DBPGX.WaitForNotification(context.Background())
			if err != nil {
				log.Fatalf("Error waiting for notification: %v\n", err)
			}

			if notification.Channel == "onapproved" {

				var payload NotificationPayload
				err = json.Unmarshal([]byte(notification.Payload), &payload)
				if err != nil {
					log.Fatalf("Error unmarshalling JSON: %v\n", err)
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
							log.Println("Error marshalling data:", err)
						}

						socketio.To(socket.Room(wsLog.SocketID)).Emit("doEvent", string(jsonData))
					}
				}

			}
			// fmt.Printf("Received notification on channel %s: %s\n", notification.Channel, notification.Payload)

			// var payload NotificationPayload
			// err = json.Unmarshal([]byte(notification.Payload), &payload)
			// if err != nil {
			// 	log.Fatalf("Error unmarshalling JSON: %v\n", err)
			// }

			// fmt.Println(&payload)

			// Emit notification to all connected Socket.IO clients
			// socket.BroadcastToNamespace("/", "notification", notification.Payload)
			// socketio.Local().Emit("notification", "test message from pg notify with data : "+*payload.ReferenceName+" = "+*payload.ReferenceId)
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

		client.On("disconnect", func(args ...interface{}) {
			// fmt.Println("Client disconnected: ", client.Id())
			client.Disconnect(true)

			// allClients2 := socketio.Of("/", nil).Sockets()
			// fmt.Println("all clients 2 : ", allClients2.Len())
		})
	})

	socketio.On("disconnect", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		fmt.Println("Client disconnected: ", client.Id())
		client.Disconnect(true)
	})

	socketio.Of("/connectCustomer", nil).On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		client.On("userId", func(args ...interface{}) {
			userIdClient := args[0].(string)
			result, err := helpers.RefreshUser(userIdClient)
			if err != nil {
				fmt.Println(err)
				client.Emit("returnData", "Failed to get user data")
			}
			client.Emit("returnData", result)
		})
	})
	app.Get("/socket.io", adaptor.HTTPHandler(socketio.ServeHandler(nil)))
	app.Post("/socket.io", adaptor.HTTPHandler(socketio.ServeHandler(nil)))

	app.Use("/ws/*", func(c *fiber.Ctx) error {
		socketio.ServeClient()
		return nil
	})
}
