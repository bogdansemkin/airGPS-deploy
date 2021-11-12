package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	//"strconv"
	//"strings"
	"time"
)

//=========================== JSON STRUCTURES ==================================

//------------------------ ALARM STRUCTURE -------------------------------------
type Alarms struct {
	//Alarms []Alarm `json:"alarms"`
	Alarms []gbrNowActiveWorkers `json:"alarms"`
}

type Alarm struct {
	ID          int    `json:"id_workings"`
	NUM_PULT    string `json:"f_object_number_pult"`
	OBJ_ADDR    string `json:"f_object_adress"`
	OBJ_NAME    string `json:"f_object_name"`
	OBJ_REGION  string `json:"f_object_name"`
	GBR_NAME    string `json:"f_gbr_number"`
	GBR_RESERVE string `json:"f_gbr_number_rezerv"`
}

type gbrNowActiveWorkers struct {
	Id_workings      string `json:"id_workings"`
	ObjectNumberPult string `json:"f_object_number_pult"`
	ObjectAdress     string `json:"f_object_adress"`
	ObjectName       string `json:"f_object_name"`
	Region           string `json:"f_region"`
	GbrNumber        string `json:"f_gbr_number"`
	GbrNumberRezerv  string `json:"f_gbr_number_rezerv"`
	IdGBR            string `json:"id_gbr"`
}

//--------------------------- GBR STRUCTURE ------------------------------------
type GbrListFull struct {
	GbrListFull []gbrListAll `json:"gbrlist"`
}

type gbrListAll struct {
	Id_gbr int    `json:"id_gbr"`
	Region string `json:"region"`
	Numgbr string `json:"number"`
}

//------------------------------------------------------------------------------
type sendStatusOfAlarm struct {
	Status string `json:"status"`
	Param  string `json:"param"`
	Id     string `json:"id"`
}

type CardBase struct {
	ID                       string   `json:"id"`
	CARD_TYPE                string   `json:"type_object_cart"`
	CARD_AFFILATION          string   `json:"affiliation"`
	CARD_INSTALLER           string   `json:"installer"`
	CARD_CLIENT              string   `json:"field_client"`
	CARD_PULTNUM             string   `json:"field_pult_number"`
	CARD_RADIO_CHANEL        string   `json:"field_radio_chanel"`
	CARD_RADIO_CHANEL_RESERV string   `json:"field_radio_chanel_rezerv"`
	CARD_REGION              string   `json:"field_region"`
	CARD_PEREZVON            string   `json:"perezvon"`
	CARD_GBR_ACTION          string   `json:"gbr_action"`
	CARD_CALL                string   `json:"field_call"`
	CARD_CALL_RESERV         string   `json:"field_call_rezerv"`
	CARD_CALL_RESERV2        string   `json:"field_call_rezerv2"`
	CARD_TIME_RESPONSE       string   `json:"field_time_response"`
	CARD_CONTROL_PANEL       string   `json:"field_contol_panel"`
	CARD_COMPANY             string   `json:"field_company"`
	CARD_ALIENT_PULT         string   `json:"field_alien_pult"`
	CARD_NAME                string   `json:"field_name"`
	CARD_ADRES               string   `json:"field_adress"`
	CARD_MODE                string   `json:"field_mode"`
	CARD_TYPE_OBJECT         string   `json:"field_type_object"`
	CARD_EXTRACT_ADDRESS     string   `json:"exact_address"`
	CARD_STOREYS             string   `json:"storeys"`
	CARD_FLOOR               string   `json:"floor"`
	CARD_KEY_PRESENT         string   `json:"key_availability"`
	CARD_HAVE_DOG            string   `json:"having_dog"`
	CARD_OUT_INTO            string   `json:"build_out_or_into"`
	CARD_WINDOW_DOOR         string   `json:"window_and_dor"`
	CARD_SECURITY            string   `json:"security_in_object"`
	CARD_WAYMARK             string   `json:"waymark"`
	CARD_PORCH               string   `json:"field_porch"`
	CARD_VULNER              string   `json:"field_vulnerabilities"`
	CARD_INFO2               string   `json:"field_description_2"`
	CARD_EQUIP               string   `json:"field_equipment"`
	CARD_WHOSE_EQUIP         string   `json:"field_whose_equipment"`
	CARD_AUTHOR              string   `json:"field_author"`
	CARD_MANAGER             string   `json:"field_manager"`
	CARD_DOGOVOR             string   `json:"field_dogovor"`
	CARD_SUM_MONTH           string   `json:"field_summ_in_month"`
	CARD_PEOPLE              []People `json:"field_people"`
	CARD_SHLEYF              []Zone   `json:"field_shleif"`

	CARD_DATE_ENTER string `json:"field_date_enter_object"`
	CARD_START_SEC  string `json:"field_date_start_security"`
	CARD_WARNING    string `json:"field_warning"`
	CARD_LAT        string `json:"lat"`
	CARD_LON        string `json:"lng"`

	CARD_FILES []string `json:"files"`
}

type AirQuery struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Cmnd  string `json:"cmnd"`
	Param string `json:"param"`
}

//------------------------------------------------------------------------------

type AppSend struct {
	Command string `json:"cmnd"`
	ID      string `json:"id"`
}

type People struct {
	MAN_NUM    string `json:"number_people_line"`
	MAN_NAME   string `json:"field_people_name"`
	MAN_ADDR   string `json:"field_adress"`
	MAN_PHONE  string `json:"field_phone"`
	MAN_ACCESS string `json:"field_access"`
	MAN_NOTE   string `json:"field_note"`
	//	Peoples []string `json:"users"`
}
type Peoples struct {
	Peoples []People `json:"field_people"`
}

//------------------------------------------------------------------------------
type Zone struct {
	ZONE_NUM   string `json:"number_shleif_line"`
	ZONE_NAME  string `json:"field_shleif_name"`
	ZONE_PLACE string `json:"field_shleif_place"`
	//	Zoness []string `json:"zones"`
}

type Zones struct {
	Zones []Zone `json:"field_shleif"`
}

//------------------------------------------------------------------------------

type Done struct{}

//=========================== JSON DECODING ====================================
func addJsonHeader(header, body string) string {
	return "{" + string(34) + header + string(34) + ": " + body + "}"
}

//=========================== DECODE GBR FUNC ==================================
func decodeGbrList(jsonIncoming []byte) {
	var (
		gbrs GbrListFull
	)

	strJSON := []byte(addJsonHeader("gbrlist", string(jsonIncoming)))
	err := json.Unmarshal(strJSON, &gbrs)
	if err != nil {
		//		defer recoveryAppFunction()
		fmt.Println("Error decoding Json:", strJSON)
		panic(err)
	}

	uptGBRlist = len(gbrs.GbrListFull)
	if uptGBRlist > 0 {
		for i := 0; i < uptGBRlist; i++ {
			gbrListID = append(gbrListID, gbrs.GbrListFull[i].Id_gbr)
			gbrListREG = append(gbrListREG, gbrs.GbrListFull[i].Region)
			gbrListNUM = append(gbrListNUM, gbrs.GbrListFull[i].Numgbr)
			gbrAlarmLast = append(gbrAlarmLast, "")
			gbrAlarmSend = append(gbrAlarmSend, "")
		}
	}

}

//------------------------------------------------------------------------------
func getNameGBR(gbrUIN int) string {
	for i := 0; i < uptGBRlist; i++ {
		if gbrListID[i] == gbrUIN {
			return gbrListNUM[i]
		}
	}
	return "Unknown"
}

//=========================== DECODE ALARMS ====================================
func decodeAlarmList(jsonIncoming []byte) {
	var (
		alarms Alarms
	)
	newAlarmList := string(jsonIncoming)

	if lastAlarmList != newAlarmList && len(newAlarmList) > 10 {
		lastAlarmList = newAlarmList
		uptAlarm = true
		strJSON := []byte(addJsonHeader("alarms", string(jsonIncoming)))
		err := json.Unmarshal(strJSON, &alarms)
		if err != nil {
			//		defer recoveryAppFunction()
			fmt.Println("Error decoding Json:", strJSON)
			panic(err)
		}
		// CLEAR SLICE
		alarmList_ID = make([]string, 0)
		alarmList_NUM_PULT = make([]string, 0)
		alarmList_OBJ_ADDR = make([]string, 0)
		alarmList_OBJ_NAME = make([]string, 0)
		alarmList_OBJ_REGION = make([]string, 0)
		alarmList_GBR_NAME = make([]string, 0)
		alarmList_GBR_RESERVE = make([]string, 0)
		alarmList_GBR_ID = make([]string, 0)
		alarmList_WAS_SEND = make([]bool, 0)
		alarmList_HAVE_CARD = make([]string, 0)
		for i := 0; i < len(alarms.Alarms); i++ {
			alarmList_ID = append(alarmList_ID, alarms.Alarms[i].Id_workings)
			alarmList_NUM_PULT = append(alarmList_NUM_PULT, alarms.Alarms[i].ObjectNumberPult)
			alarmList_OBJ_ADDR = append(alarmList_OBJ_ADDR, alarms.Alarms[i].ObjectAdress)
			alarmList_OBJ_NAME = append(alarmList_OBJ_NAME, alarms.Alarms[i].ObjectName)
			alarmList_OBJ_REGION = append(alarmList_OBJ_REGION, alarms.Alarms[i].Region)
			alarmList_GBR_NAME = append(alarmList_GBR_NAME, alarms.Alarms[i].GbrNumber)
			alarmList_GBR_RESERVE = append(alarmList_GBR_RESERVE, alarms.Alarms[i].GbrNumberRezerv)
			alarmList_GBR_ID = append(alarmList_GBR_ID, alarms.Alarms[i].IdGBR)
			alarmList_WAS_SEND = append(alarmList_WAS_SEND, false)
			alarmList_HAVE_CARD = append(alarmList_HAVE_CARD, "")

			fmt.Println("ID: ", alarms.Alarms[i].Id_workings)
			fmt.Println("NUM_PULT: " + alarms.Alarms[i].ObjectNumberPult)
			fmt.Println("OBJ_ADDR: " + alarms.Alarms[i].ObjectAdress)
			fmt.Println("OBJ_NAME: " + alarms.Alarms[i].ObjectName)
			fmt.Println("OBJ_REGION: " + alarms.Alarms[i].Region)
			fmt.Println("GBR_NAME: " + alarms.Alarms[i].GbrNumber)
			fmt.Println("GBR_RESERVE: " + alarms.Alarms[i].GbrNumberRezerv)
			fmt.Println("GBR_ID: " + alarms.Alarms[i].IdGBR)
			checkUpdateAlarms(i, alarmList_NUM_PULT[i])
			fmt.Println("==============================")

		}

	}
}

//============================= Decode CARD ====================================
func prepareStringForSemen(nopStr string) string {
	s_result := doReplaceStr(nopStr, string(32)+string(92)+string(34), " <<")
	s_result = doReplaceStr(s_result, string(92)+string(34), ">>") //+string(32)

	return s_result
}

//------------------------------------------------------------------------------
func decodeZoneList(strJSON []byte) string {

	var (
		zones  Zones
		s_json string
	)
	s_json = ""
	//strJSON = []byte(jsonIncoming)
	err := json.Unmarshal(strJSON, &zones)
	if err != nil {
		//		defer recoveryAppFunction()
		fmt.Println("Error decoding Json:", strJSON)
		panic(err)
	}
	for i := 0; i < len(zones.Zones); i++ {
		s_num := zones.Zones[i].ZONE_NUM
		s_name := zones.Zones[i].ZONE_NAME
		s_info := zones.Zones[i].ZONE_PLACE
		fmt.Println("ZONE_NUM: " + s_num)
		fmt.Println("ZONE_NAME: " + s_name)
		fmt.Println("ZONE_PLACE: " + s_info)
		s_json += "{" + getQuatedJSON("name", s_name, 1) + "," +
			getQuatedJSON("num", s_num, 1) + "," + string(34) +
			getQuatedJSON("tel", s_info, 1) + "}"

	}
	s_json = string(34) + "zonelist" + string(34) + ":" + "[" + s_json + "]"
	return s_json
}

//------------------------------------------------------------------------------
func decodePeopleList(strJSON []byte) string {

	var (
		peoples Peoples
		s_json  string
	)
	s_json = ""
	//strJSON = []byte(jsonIncoming)
	err := json.Unmarshal(strJSON, &peoples)
	if err != nil {
		//		defer recoveryAppFunction()
		fmt.Println("Error decoding Json:", strJSON)
		panic(err)
	}
	for i := 0; i < len(peoples.Peoples); i++ {
		if len(s_json) > 0 {
			s_json += ","
		}
		s_num := peoples.Peoples[i].MAN_NUM
		s_name := peoples.Peoples[i].MAN_NAME
		//s_address := peoples.Peoples[i].MAN_ADDR
		s_phone := peoples.Peoples[i].MAN_PHONE
		s_json += "{" + getQuatedJSON("name", s_name, 1) + "," +
			getQuatedJSON("num", s_num, 1) + "," + string(34) +
			getQuatedJSON("tel", s_phone, 1) + "}"

	}
	s_json = string(34) + "userlist" + string(34) + ":" + "[" + s_json + "]"
	return s_json
}

//------------------------------------------------------------------------------
func decodeImageList(strJSON []string) string {
	s_json := ""

	for i := 0; i < len(strJSON); i++ {
		if len(s_json) > 0 {
			s_json += ","
		}
		s_json += "{" + string(34) + "url" + string(34) + ":" + string(34) + strJSON[i] + string(34) + "}"
	}
	s_json = string(34) + "imagelist" + string(34) + ":" + "[" + s_json + "]"

	return s_json
}

//------------------------------------------------------------------------------
func decodeGbrCard(cardPos int, jsonIncoming []byte) {
	var (
		cardBase CardBase
	)
	card_valid := true
	if cardPos < len(alarmList_HAVE_CARD) {
		err := json.Unmarshal(jsonIncoming, &cardBase)
		if err != nil {
			card_valid = false
			defer recoveryAppFunction()
			fmt.Println("Error decoding Json:", cardPos) //jsonIncoming
			panic(err)
		}
		//2 {"obinfo":[{""id": "7238","lat": "48.515266","lon": "34.613621","obadr": "г. Каменское, проспект Свободы, 51","obname": "Банк \"Пумб\" отделение №1","pult": "2282","details": ""}]
		s_cardBase_ID := cardBase.ID
		s_cardBase_CARD_LAT := cardBase.CARD_LAT
		s_cardBase_CARD_LON := cardBase.CARD_LON
		s_cardBase_CARD_ADRES := prepareStringForSemen(cardBase.CARD_ADRES)
		s_cardBase_CARD_NAME := prepareStringForSemen(cardBase.CARD_NAME)
		s_cardBase_CARD_PULTNUM := cardBase.CARD_PULTNUM
		s_cardBase_CARD_WAYMARK := prepareStringForSemen(cardBase.CARD_WAYMARK)
		s_cardBase_CARD_FILES := cardBase.CARD_FILES

		s_json := "{" + string(34) + "obinfo" + string(34) + ":" + "[" + "{" + string(34) +
			getQuatedJSON("id", s_cardBase_ID, 1) + "," +
			getQuatedJSON("lat", s_cardBase_CARD_LAT, 1) + "," +
			getQuatedJSON("lon", s_cardBase_CARD_LON, 1) + "," +
			getQuatedJSON("obadr", s_cardBase_CARD_ADRES, 1) + "," +
			getQuatedJSON("obname", s_cardBase_CARD_NAME, 1) + "," +
			getQuatedJSON("pult", s_cardBase_CARD_PULTNUM, 1) + "," +
			getQuatedJSON("details", s_cardBase_CARD_WAYMARK, 1) + "}" + "]"
		s_json += "," + decodePeopleList(jsonIncoming) + ","
		s_json += decodeZoneList(jsonIncoming) + ","
		s_json += string(34) + "eventlist" + string(34) + ":[{}],"
		s_json += string(34) + "imagelist" + string(34) + ":" + "[" + decodeImageList(s_cardBase_CARD_FILES) + "]"
		s_json += "}"

		if card_valid {
			fmt.Println(cardPos, s_json)
		} else {
			fmt.Println(cardPos, "ERROR DECODE")
		}

	}

}

//======================== HTTP FUNCTIONS ======================================
func getQueryLink(accUIN string) (int, string) {

	if uptGBRlist == 0 {
		return 0, "http://api-cs.ohholding.com.ua/gbr_list/get"
	} else if len(accUIN) > 3 {
		return 2, "http://api-cs.ohholding.com.ua/object_cart/" + accUIN + "/get"
	}
	return 1, "http://api-cs.ohholding.com.ua/active_workings/get"

}

//------------------------------------------------------------------------------
func checkUpdateAlarms(cardPos int, accUIN string) { // Check New Alarms
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	i_query, s_guery := getQueryLink(accUIN)
	req, err := http.NewRequest("GET", s_guery, nil)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	req.Header.Add("Accept", `application/json`)
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	fmt.Printf("Body : %s \n ", body)
	fmt.Printf("Response status : %s \n", resp.Status)
	if i_query == 0 {
		decodeGbrList(body)
	} else if i_query == 2 {
		decodeGbrCard(cardPos, body)
	} else {
		decodeAlarmList(body)
	}
}
