package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	//	"strconv"
	"time"

	//"os"

	"github.com/gorilla/websocket"
)

type msgweb struct {
	Input string
	//message string
}

var (
	websock_addr_user    []*websocket.Conn //Current socket
	websock_last_connect []int64           //Last connection date/time
	websock_uin_users    []int             //Connection Index
	websock_uin_gbr      []string          //Connection GBR Name
	websock_send_repeat  []int             //Send update counter
	websock_send_device  []string          //Current tocken
	websock_addr_counter int               //Connections counter

)

func recoverySocketFunction() {
	if recoveryMessage := recover(); recoveryMessage != nil {
		fmt.Println(recoveryMessage)
	}
	fmt.Println(getDT(), "Application restored after Error...")
}

//-----------------------------------------------------------------------------

func rootHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		defer recoverySocketFunction()
		fmt.Println("Could not open file.", err)
	}
	fmt.Fprintf(w, "%s", content)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	can_update := true

	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		can_update = false
		defer recoverySocketFunction()
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	if conn.LocalAddr() == nil {
		can_update = false
	} else if len(conn.LocalAddr().String()) < 8 {
		can_update = false
	}

	if conn != nil && can_update {
		go echosock(conn)
	}

}

func conRemoval(conn *websocket.Conn) {
	var (
		websock_temp_user    []*websocket.Conn
		websock_temp_connect []int64
		websock_temp_uin     []int
		websock_temp_gbr     []string
		websock_temp_repeat  []int
		websock_temp_device  []string
		websock_temp_counter int
	)

	if websock_addr_counter < 1 {
		return
	}
	websock_temp_user = make([]*websocket.Conn, 0)
	websock_temp_connect = make([]int64, 0)
	websock_temp_uin = make([]int, 0)
	websock_temp_gbr = make([]string, 0)
	websock_temp_repeat = make([]int, 0)
	websock_temp_device = make([]string, 0)

	websock_temp_counter = websock_addr_counter
	for i := 0; i < websock_temp_counter; i++ {
		if i >= len(websock_addr_user) {
			if websock_addr_counter > 0 {
				websock_addr_counter--
			}
		} else if conn != websock_addr_user[i] {
			websock_temp_user = append(websock_temp_user, websock_addr_user[i])

			websock_temp_connect = append(websock_temp_connect, websock_last_connect[i])
			websock_temp_uin = append(websock_temp_uin, websock_uin_users[i])
			websock_temp_gbr = append(websock_temp_gbr, websock_uin_gbr[i])
			websock_temp_repeat = append(websock_temp_repeat, websock_send_repeat[i])
			websock_temp_device = append(websock_temp_device, websock_send_device[i])

		} else {
			if websock_addr_counter > 0 {
				for j := 0; j < len(gbrAlarmSocket); j++ {
					if gbrAlarmSocket[j] == i {
						gbrAlarmSocket[j] = -1
					}
				}
				websock_addr_counter--
			}
		}
	}
	websock_addr_user = websock_temp_user
	websock_last_connect = websock_temp_connect
	websock_uin_users = websock_temp_uin
	websock_uin_gbr = websock_temp_gbr
	websock_send_repeat = websock_temp_repeat
	websock_send_device = websock_temp_device

}

//------------------------------------------------------------------------------
func updateSockList() {
	var (
		s_json     string
		jsonResult map[string]interface{}
	)

	for i := 0; i < websock_addr_counter; i++ {
		s_json = sendUpdator(i)
		errJson := json.Unmarshal([]byte(s_json), &jsonResult)
		if errJson != nil {
			defer recoverySocketFunction()
			fmt.Println("Error message: ", s_json, errJson)
		}

		if i < len(websock_addr_user) {
			if websock_addr_user[i] != nil && errJson == nil {

				if err := websock_addr_user[i].WriteJSON(jsonResult); err != nil {
					defer recoverySocketFunction()
					fmt.Println("Error: Websocket has send", i)
					fmt.Println(err)
					defer conRemoval(websock_addr_user[i])
				}
			} else {
				fmt.Println("Error: Websocket has nil value", i)
			}
		} else {
			fmt.Println("Error: Websocket out of range", i)
		}
	}
}

//------------------------------------------------------------------------------
func sendALL(resendSTR string, conn *websocket.Conn) {

	for i := 0; i < websock_addr_counter; i++ {
		if conn != websock_addr_user[i] {
			if err := websock_addr_user[i].WriteJSON("Resend data: " + resendSTR); err != nil {
				defer recoverySocketFunction()
				fmt.Println(err)
			}
		}
	}
}

//------------------------------------------------------------------------------
func sendYeden(connIndex int, conJson string) {

	var (
		resendSTR  string
		jsonResult map[string]interface{}
	)
	if websock_addr_counter < 1 {
		return
	}
	for i := 0; i < websock_addr_counter; i++ {

		if connIndex == i {
			if websock_send_repeat[i] == 0 {
				websock_send_repeat[i] = 3
			} else if websock_send_repeat[i] == 1 {
				websock_send_repeat[i] = 0
			}
			resendSTR = conJson

			errJson := json.Unmarshal([]byte(resendSTR), &jsonResult)
			if errJson != nil {
				defer recoverySocketFunction()
				fmt.Println("Error message: ", resendSTR, errJson)
			}

			if i < len(websock_addr_user) {
				if websock_addr_user[i] != nil && errJson == nil {

					if err := websock_addr_user[i].WriteJSON(jsonResult); err != nil {
						defer recoverySocketFunction()
						fmt.Println("Error: Websocket has send", i)
						fmt.Println(err)
						defer conRemoval(websock_addr_user[i])
					}
				} else {
					fmt.Println("Error: Websocket has nil value", i)
				}

			} else {
				fmt.Println("Error: Websocket out of range", i)
			}

		}
	}
}

//------------------------------------------------------------------------------
func sendUpdateSock(connDevice, connIndex int) {

	if websock_addr_counter < 1 {
		return
	}
	for i := 0; i < websock_addr_counter; i++ {
		if i < len(websock_uin_users) {
			if connIndex == websock_uin_users[i] {
				websock_send_repeat[i] = 4
			}
		} else if websock_addr_counter > 0 {
			websock_addr_counter--
		}
	}
}

//------------------------------------------------------------------------------
func echosock(conn *websocket.Conn) {

	websock_addr_user = append(websock_addr_user, conn)
	websock_last_connect = append(websock_last_connect, time.Now().Unix())
	websock_uin_users = append(websock_uin_users, 0)
	websock_uin_gbr = append(websock_uin_gbr, "")
	websock_send_repeat = append(websock_send_repeat, 0)
	websock_send_device = append(websock_send_device, "")
	websock_addr_counter++

	for {
		con_index, m, err := conn.ReadMessage() //msg{}
		if con_index < 0 {
			//defer conn.CloseHandler()
			defer conn.Close()
			conRemoval(conn)
			fmt.Println("Event close: ", con_index, websock_addr_counter, m, conn.RemoteAddr())
			return
			//os.Exit(0)
		}

		//	m := msg{}
		//	err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading json.", err)

		}
		jsonRead := string(m)
		//fmt.Println("Read message: ", con_index, websock_addr_counter, jsonRead, conn.RemoteAddr())
		jsonSend := decodeGpsJson(jsonRead, conn)
		var jsonResult map[string]interface{}
		errJson := json.Unmarshal([]byte(jsonSend), &jsonResult)
		if errJson != nil {
			fmt.Println("Error message: ", jsonSend, errJson)
			recoverySocketFunction()
		}

		fmt.Println(jsonResult)

		if err = conn.WriteJSON(jsonResult); err != nil {
			recoverySocketFunction()
			fmt.Println(err)
		}
		//sendALL(s, conn)
	}
}

//------------------------------------------------------------------------------
func getSocketIndex(conn *websocket.Conn) int {

	for i := 0; i < websock_addr_counter; i++ {
		if i < len(websock_addr_user) {
			if conn == websock_addr_user[i] {
				return i
			}
		} else if websock_addr_counter > 0 {
			websock_addr_counter--
		}
	}
	return -1
}

//------------------------------------------------------------------------------
func broker_upgrade() {

	for i := 0; i < websock_addr_counter; i++ {
		if websock_send_repeat[i] == 4 {
			websock_send_repeat[i] = 0
			sendYeden(i, sendUpdator(i))
		} else if websock_send_repeat[i] > 1 {
			websock_send_repeat[i]--
		} else if websock_send_repeat[i] == 1 {
			sendYeden(i, sendUpdator(i))
		}
	}
}
