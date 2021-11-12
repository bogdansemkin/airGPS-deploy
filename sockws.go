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
	websock_gbr_addr     []*websocket.Conn //Current socket
	websock_gbr_last_con []int64           //Last connection date/time
	websock_gbr_uin      []int             //Connection GBR UIN
	websock_gbr_name     []string          //Connection GBR Name
	websock_gbr_repeat   []int             //Send update counter
	websock_gbr_tocken   []string          //Current tocken
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
		websock_gbr_addr_temp     []*websocket.Conn //Current socket
		websock_gbr_last_con_temp []int64           //Last connection date/time
		websock_gbr_uin_temp      []int             //Connection Index
		websock_gbr_name_temp     []string          //Connection GBR Name
		websock_gbr_repeat_temp   []int             //Send update counter
		websock_gbr_tocken_temp   []string          //Current tocken
		websock_addr_counter_temp int               //Connections counter

	)

	if websock_addr_counter < 1 {
		return
	}
	websock_gbr_addr_temp = make([]*websocket.Conn, 0) //Current socket
	websock_gbr_last_con_temp = make([]int64, 0)       //Last connection date/time
	websock_gbr_uin_temp = make([]int, 0)              //Connection Index
	websock_gbr_name_temp = make([]string, 0)          //Connection GBR Name
	websock_gbr_repeat_temp = make([]int, 0)           //Send update counter
	websock_gbr_tocken_temp = make([]string, 0)        //Current tocken
	websock_addr_counter_temp = websock_addr_counter   //Connections counter

	for i := 0; i < websock_addr_counter_temp; i++ {
		if i >= len(websock_gbr_addr) {
			if websock_addr_counter > 0 {
				websock_addr_counter--
			}
		} else if conn != websock_gbr_addr[i] {
			websock_gbr_addr_temp = append(websock_gbr_addr_temp, websock_gbr_addr[i])
			websock_gbr_last_con_temp = append(websock_gbr_last_con_temp, websock_gbr_last_con[i])
			websock_gbr_uin_temp = append(websock_gbr_uin_temp, websock_gbr_uin[i])
			websock_gbr_name_temp = append(websock_gbr_name_temp, websock_gbr_name[i])
			websock_gbr_repeat_temp = append(websock_gbr_repeat_temp, websock_gbr_repeat[i])
			websock_gbr_tocken_temp = append(websock_gbr_tocken_temp, websock_gbr_tocken[i])
		} else {
			if websock_addr_counter > 0 {

				websock_addr_counter--
			}
		}
	}
	websock_gbr_addr = websock_gbr_addr_temp
	websock_gbr_last_con = websock_gbr_last_con_temp
	websock_gbr_uin = websock_gbr_uin_temp
	websock_gbr_name = websock_gbr_name_temp
	websock_gbr_repeat = websock_gbr_repeat_temp
	websock_gbr_tocken = websock_gbr_tocken_temp
	websock_addr_counter = len(websock_gbr_addr)

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
			fmt.Println("Error updateSockList message: ", s_json, errJson)
		}

		if i < len(websock_gbr_addr) {
			if websock_gbr_addr[i] != nil && errJson == nil {

				if err := websock_gbr_addr[i].WriteJSON(jsonResult); err != nil {
					defer recoverySocketFunction()
					fmt.Println("Error: Websocket has send", i)
					fmt.Println(err)
					defer conRemoval(websock_gbr_addr[i])
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
		if conn != websock_gbr_addr[i] {
			if err := websock_gbr_addr[i].WriteJSON("Resend data: " + resendSTR); err != nil {
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
			if websock_gbr_repeat[i] == 0 {
				websock_gbr_repeat[i] = 3
			} else if websock_gbr_repeat[i] == 1 {
				websock_gbr_repeat[i] = 0
			}
			resendSTR = conJson

			errJson := json.Unmarshal([]byte(resendSTR), &jsonResult)
			if errJson != nil {
				defer recoverySocketFunction()
				fmt.Println("Error sendYeden message: ", resendSTR, errJson)
			}

			if i < len(websock_gbr_addr) {
				if websock_gbr_addr[i] != nil && errJson == nil {

					if err := websock_gbr_addr[i].WriteJSON(jsonResult); err != nil {
						defer recoverySocketFunction()
						fmt.Println("Error: Websocket has send", i)
						fmt.Println(err)
						defer conRemoval(websock_gbr_addr[i])
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
		if i < len(websock_gbr_uin) {
			if connIndex == websock_gbr_uin[i] {
				websock_gbr_repeat[i] = 4
			}
		} else if websock_addr_counter > 0 {
			websock_addr_counter--
		}
	}
}

//------------------------------------------------------------------------------
func echosock(conn *websocket.Conn) {
	websock_gbr_addr = append(websock_gbr_addr, conn)
	websock_gbr_last_con = append(websock_gbr_last_con, time.Now().Unix())
	websock_gbr_uin = append(websock_gbr_uin, 0)
	websock_gbr_name = append(websock_gbr_name, "")
	websock_gbr_repeat = append(websock_gbr_repeat, 0)
	websock_gbr_tocken = append(websock_gbr_tocken, "")

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
			fmt.Println("Error echosock message: ", jsonSend, errJson)
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

	for i := 0; i < len(websock_gbr_addr); i++ {
		if conn == websock_gbr_addr[i] {
			return i
		}
	}

	// for i := 0; i < websock_addr_counter; i++ {
	// 	if i < len(websock_gbr_addr) {
	// 		if conn == websock_gbr_addr[i] {
	// 			return i
	// 		}
	// 	} else if websock_addr_counter > 0 {
	// 		websock_addr_counter--
	// 	}
	// }
	return -1
}

//------------------------------------------------------------------------------
func broker_upgrade() {

	for i := 0; i < websock_addr_counter; i++ {
		if websock_gbr_repeat[i] == 4 {
			websock_gbr_repeat[i] = 0
			sendYeden(i, sendUpdator(i))
		} else if websock_gbr_repeat[i] > 1 {
			websock_gbr_repeat[i]--
		} else if websock_gbr_repeat[i] == 1 {
			sendYeden(i, sendUpdator(i))
		}
	}
}
