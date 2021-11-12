package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
)

//========================= MAIN LOGIC =======================================
func decodeGpsJson(jsonIncoming string, conn *websocket.Conn) string {
	var (
		airDecoding AirQuery
		strJSON     []byte
		//		i_con       int
		js_result string
		js_iden   string
		js_cmnd   string
		js_param  string
		js_name   string
	)
	js_result = "{" + string(0x0D) + string(0x0A)
	js_result += getQuatedJSON("param", "Status error", 1) + string(0x0D) + string(0x0A)
	js_result += "}" + string(0x0D) + string(0x0A)

	fmt.Println(jsonIncoming)

	//Error json format
	if checkValidJson(jsonIncoming) == true { //Check valid data
		js_result = "Status 8" //Error data
		strJSON = []byte(jsonIncoming)
		err := json.Unmarshal(strJSON, &airDecoding)
		if err != nil {
			defer recoveryAppFunction()
			fmt.Println(getDT(), "Error decoding Json:"+jsonIncoming)
			panic(err)
		}
		js_iden = airDecoding.ID
		js_name = airDecoding.Name
		js_cmnd = airDecoding.Cmnd
		js_param = airDecoding.Param

		switch js_cmnd {
		case "start": //First start
			js_result = startGBR(js_iden, js_name, js_param, getSocketIndex(conn))
		case "login": //Loging for user
			js_result = logGBR(js_iden, js_name, js_param, getSocketIndex(conn))
		case "connect": //Check New Connection
			js_result = setUnknown(js_iden, js_name, js_cmnd, "Connect_OK")
		case "alarmget": //Receive alarm

		case "alarmstart": //GBR starts trip
			js_result = procAlarm(js_iden, js_cmnd, js_name, js_param)
		case "alarmpoint": //GBR at point
			js_result = procAlarm(js_iden, js_cmnd, js_name, js_param)
		case "alarmbreak": //Problem with GBR
			js_result = procAlarm(js_iden, js_cmnd, js_name, js_param)
		case "alarmstop": //Set reaction
			js_result = procAlarm(js_iden, js_cmnd, js_name, js_param)
		case "alarminfo": //Read updates

		default:
			js_result = setUnknown(js_iden, js_name, js_cmnd, "STATUS_ERROR")
		}
	}
	return js_result
}

//------------------------------------------------------------------------------
func sendUpdator(userid int) string {
	s_json := "{" + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("id", strconv.Itoa(websock_gbr_uin[userid]), 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("cmnd", "update", 1) + string(0x0D) + string(0x0A)
	s_json = s_json + "}" + string(0x0D) + string(0x0A)
	return s_json
}

//------------------------------------------------------------------------------
func setUnknown(userid, js_name, js_command, js_param string) string {
	fmt.Println(getDT(), "Command unknown"+js_command)
	s_json := "{" + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("id", userid, 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("name", js_name, 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("cmnd", js_command, 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("param", js_param, 1) + string(0x0D) + string(0x0A)
	s_json = s_json + "}" + string(0x0D) + string(0x0A)

	return s_json
}

//------------------------------------------------------------------------------
func startGBR(userid, js_name, js_param string, conPosition int) string {
	//{"cmnd":"start","id":"0","name":"token","param":"semen2021"}
	s_json := ""
	if start_pass == js_param { //Device loging
		//userid - IMEI of device; js_name - google accound addr; js_param - password
		s_json = "{" + string(0x0D) + string(0x0A)
		s_json += getQuatedJSON("gbrlist", "[", 0) + string(0x0D) + string(0x0A)
		s_json += getGBRlist() + "]" //GBR LIST
		s_json += "}" + string(0x0D) + string(0x0A)

		s_json = doReplaceStr(s_json, "},]", "}]")

		if conPosition >= 0 && conPosition < websock_addr_counter {
			websock_gbr_tocken[conPosition] = js_name //SET TOCKEN
		}

	} else {
		s_json = "{" + string(0x0D) + string(0x0A)
		s_json += getQuatedJSON("id", userid, 1) + "," + string(0x0D) + string(0x0A)
		s_json += getQuatedJSON("name", js_name, 1) + "," + string(0x0D) + string(0x0A)
		s_json += getQuatedJSON("cmnd", "start", 1) + "," + string(0x0D) + string(0x0A)
		s_json += getQuatedJSON("param", "START_ERR", 1) + string(0x0D) + string(0x0A)
		s_json += "}" + string(0x0D) + string(0x0A)
	}
	return s_json
}

//------------------------------------------------------------------------------

func logGBR(userid, js_name, js_param string, conPosition int) string {
	s_json := ""
	gbrvalid := getGBRuin(userid)
	if gbrvalid == false || js_name != "-2" || js_param != "-111" {
		s_json = "{" + string(0x0D) + string(0x0A)
		s_json = s_json + getQuatedJSON("id", userid, 1) + "," + string(0x0D) + string(0x0A)
		s_json = s_json + getQuatedJSON("name", js_name, 1) + "," + string(0x0D) + string(0x0A)
		s_json = s_json + getQuatedJSON("cmnd", "login", 1) + "," + string(0x0D) + string(0x0A)
		s_json = s_json + getQuatedJSON("param", "GBR_ERR", 1) + string(0x0D) + string(0x0A)
		s_json = s_json + "}" //+ string(0x0D) + string(0x0A)
	} else if len(js_name) < 1 || len(js_param) < 1 { //Input data is empty
		s_json = "{" + string(0x0D) + string(0x0A)
		s_json = s_json + getQuatedJSON("id", userid, 1) + "," + string(0x0D) + string(0x0A)
		s_json = s_json + getQuatedJSON("name", js_name, 1) + "," + string(0x0D) + string(0x0A)
		s_json = s_json + getQuatedJSON("cmnd", "login", 1) + "," + string(0x0D) + string(0x0A)
		s_json = s_json + getQuatedJSON("param", "LOG_EMPTY", 1) + string(0x0D) + string(0x0A)
		s_json = s_json + "}" //+ string(0x0D) + string(0x0A)
	} else { //Not EMpty data
		s_psw := "-111"
		if len(s_psw) < 1 { //NOT LOGIN ENABLE
			s_json = "{" + string(0x0D) + string(0x0A)
			s_json = s_json + getQuatedJSON("id", userid, 1) + "," + string(0x0D) + string(0x0A)
			s_json = s_json + getQuatedJSON("name", js_name, 1) + "," + string(0x0D) + string(0x0A)
			s_json = s_json + getQuatedJSON("cmnd", "login", 1) + "," + string(0x0D) + string(0x0A)
			s_json = s_json + getQuatedJSON("param", "LOG_FALSE_L", 1) + string(0x0D) + string(0x0A)
			s_json = s_json + "}" //+ string(0x0D) + string(0x0A)
		} else if s_psw == js_param { //ALL OK

			s_json = "{" + string(0x0D) + string(0x0A)
			s_json = s_json + getQuatedJSON("id", userid, 1) + "," + string(0x0D) + string(0x0A)
			s_json = s_json + getQuatedJSON("name", js_name, 1) + "," + string(0x0D) + string(0x0A)
			s_json = s_json + getQuatedJSON("cmnd", "login", 1) + "," + string(0x0D) + string(0x0A)
			s_json = s_json + getQuatedJSON("param", "LOG_OK", 1) + string(0x0D) + string(0x0A)
			s_json = s_json + "}" //+ string(0x0D) + string(0x0A)
			s_tocken := ""

			if conPosition >= 0 && conPosition < websock_addr_counter {
				s_tocken = websock_gbr_tocken[conPosition]
				websock_gbr_uin[conPosition] = convertIntVal(userid)
			}
			fmt.Println("Try update tocken", userid, js_name, s_tocken)
			updateGBRstatus(userid, "", s_tocken, 0)
		} else { //NOT PASSWORD ENABLE
			s_json = "{" + string(0x0D) + string(0x0A)
			s_json = s_json + getQuatedJSON("id", userid, 1) + "," + string(0x0D) + string(0x0A)
			s_json = s_json + getQuatedJSON("name", js_name, 1) + "," + string(0x0D) + string(0x0A)
			s_json = s_json + getQuatedJSON("cmnd", "login", 1) + "," + string(0x0D) + string(0x0A)
			s_json = s_json + getQuatedJSON("param", "LOG_FALSE_x", 1) + string(0x0D) + string(0x0A)
			s_json = s_json + "}" //+ string(0x0D) + string(0x0A)
		}
	}
	return s_json
}

//------------------------------------------------------------------------------
func procAlarm(userid, js_cmnd, js_name, js_param string) string {
	s_status := "ALARM_ERR"
	switch js_cmnd {
	case "alarmstart": //GBR starts trip
		s_status = "START_OK"
	case "alarmpoint": //GBR at point
		s_status = "POINT_OK"
	case "alarmbreak": //Problem with GBR
		s_status = "BREAK_OK"
	case "alarmstop": //Set reaction
		s_status = "STOP_OK"
	}

	switch js_cmnd {
	case "alarmstart": //GBR starts trip
		updateGBRstatus(userid, "", js_param, 1)
	case "alarmpoint": //GBR at point
		updateGBRstatus(userid, "", js_param, 2)
	case "alarmbreak": //Problem with GBR
		updateGBRstatus(userid, "", js_param, 3)
	case "alarmstop": //Set reaction
		updateGBRstatus(userid, "", js_param, 4)
	}

	s_json := "{" + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("id", userid, 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("name", js_name, 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("cmnd", js_cmnd, 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("param", s_status, 1) + string(0x0D) + string(0x0A)
	s_json = s_json + "}" + string(0x0D) + string(0x0A)

	return s_json
}

//------------------------------------------------------------------------------
func procAlarmList() {
	for i := 0; i < len(alarmList_ID); i++ { //READ ALARMS
		//fmt.Println("Send alarms:", i, alarmList_ID[i], alarmList_OBJ_NAME[i], alarmList_GBR_NAME[i], alarmList_GBR_RESERVE[i])
		if alarmList_WAS_SEND[i] == false { //CHECK WAS SEND
			for j := 0; j < len(gbrListID); j++ { //SEARCH GBR
				s_gbr := strconv.Itoa(gbrListID[j])
				if alarmList_GBR_NAME[i] == s_gbr || alarmList_GBR_RESERVE[i] == s_gbr { //GBR VALID
					isGBRcon := false
					gbrAlarmReceived[j] = false
					gbrAlarmType[j] = 1
					gbrAlarmLast[j] = alarmList_ID[i]
					gbrAlarmPult[j] = alarmList_NUM_PULT[i]
					gbrAlarmName[j] = alarmList_OBJ_NAME[i]
					gbrAlarmCard[j] = alarmList_HAVE_CARD[i]
					for k := 0; k < websock_addr_counter; k++ { //SEARCH GBR IN CONNECT LIST
						if websock_gbr_name[k] == "" { //NOT NAMED GBR
							websock_gbr_name[k] = getNameGBR(websock_gbr_uin[k])
						}
						//TODO: CHECK GBR IS NOT BUSY
						//fmt.Println("Search connection:", k, "GBR", websock_gbr_name[k], "UIN", gbrListID[j], websock_gbr_uin[k])
						if websock_gbr_uin[k] == gbrListID[j] { //GBR CONNECTED
							isGBRcon = true
							gbrAlarmSend[j] = true
							fmt.Println("Find connection:", k, websock_gbr_name[k], websock_gbr_uin[k], alarmList_ID[i], alarmList_NUM_PULT[i], alarmList_OBJ_NAME[i], gbrAlarmCard[j])
							if len(gbrAlarmCard[j]) < 10 { //No card enable
								gbrAlarmCard[j] = "{" + string(34) + "no_pult" + string(34) + ":" + "[" + "{" +
									string(34) + "f_object_adress" + string(34) + ":" + string(34) + alarmList_OBJ_ADDR[i] + string(34) + "," +
									string(34) + "f_object_name" + string(34) + ":" + string(34) + alarmList_OBJ_NAME[i] + string(34) + "," +
									string(34) + "f_region" + string(34) + ":" + string(34) + alarmList_OBJ_REGION[i] + string(34) + "}" + "]" + "}"
							}
							gbrAlarmSocket[j] = k
							gbrAlarmWait[j] = 0
							alarmList_WAS_SEND[i] = true
							sendYeden(k, gbrAlarmCard[j])
						}
					} //SEARCH GBR FOR SEND PUSH
					if isGBRcon == false { //GBR NOT CONNECTED TO SOCKET
						if len(gbrAlarmTocken[j]) > 8 { //GBR TOCKEN IS VALID
							if getTokenList(gbrAlarmTocken[j], gbrAlarmPult[j], gbrAlarmName[j], 1) {
								gbrAlarmSend[j] = true
								gbrAlarmLast[j] = alarmList_ID[i]

							}
						}
					} //END - SEARCH GBR IN CONNECT LIST
				} //END - GBR VALID
			} //END - SEARCH GBR
		} //END - CHECK WAS SEND
	} //END - READ ALARMS
}

//------------------------------------------------------------------------------
func sendAlarmToGbr() {
	for i := 0; i < len(gbrListID); i++ {
		if len(gbrAlarmLast[i]) > 0 { //GROUP HAS ACTIVE ALARM
			if gbrAlarmReceived[i] == false { //GBR NOT RESPONSED
				gbrAlarmWait[i]++         //INCREMENT WAIT PERIOD
				if gbrAlarmWait[i] > 60 { //ALARM WAIT PERIOD FINISHED
					gbrAlarmSend[i] = false // NEED REPEAT SEND ALARM
				}
			}
			if gbrAlarmSend[i] == false && gbrAlarmType[i] < 2 { //ALARM NOT SENDED
				gbrAlarmWait[i] = 0         //RESET WAIT TIMER
				if gbrAlarmSocket[i] > -1 { //SOCKET ENABLE
					sendYeden(gbrAlarmSocket[i], gbrAlarmCard[i])
				} else if len(gbrAlarmTocken[i]) > 8 { //GBR FCM TOCKEN IS VALID
					if getTokenList(gbrAlarmTocken[i], gbrAlarmPult[i], gbrAlarmName[i], 1) {
						gbrAlarmSend[i] = true
					}
				}
			}
		}

	}
}
