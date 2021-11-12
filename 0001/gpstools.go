package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"sort"

	"math"
	"strconv"
	"strings"
	"time"
)

//===================== Delphi convert date functions ==========================
// 4x faster than dateutils version
func UNIXTimeToDateTimeFAST(intUnixTime int64) float64 {
	//2021-06-08 16:34:19.3509522 +0300 EEST
	s := time.Now().UTC().Local().String()
	i_pos := strings.Index(s, "+")
	sec_local := 3600 * convertInt64Val(s[i_pos+1:i_pos+3])
	return float64(intUnixTime+sec_local)/86400 + 25569
}

//------------------------------------------------------------------------------
// 10x faster than dateutils version
func DateTimeToUNIXTimeFAST(DelphiTime float64) int64 {
	return int64(math.Round((DelphiTime - 25569) * 86400))
}

//------------------------------------------------------------------------------
func floatToStr64(f_data float64) string {
	return strconv.FormatFloat(f_data, 'f', -1, 64)
}

//------------------------------------------------------------------------------
func delphiDateToSQL(unixConvDT int64) string {
	f_date := UNIXTimeToDateTimeFAST(unixConvDT)
	s_date := floatToStr64(f_date)
	return doReplaceStr(s_date, ",", ".")
}

//======================== General utils ==================================
func recoveryAppFunction() {
	if recoveryMessage := recover(); recoveryMessage != nil {
		fmt.Println(recoveryMessage)
	}
	fmt.Println(getDT(), "Application restored after Error...")
}

//-----------------------------------------------------------------------------
func getDT() string {
	currentTime := time.Now()
	return currentTime.Format("2006.01.02 15:04:05") + " > "
	//return currentTime.String() + " > "
}

//-------------------------------------------------------------------------
func getQuatedJSON(datJSON string, valJSON string, tpJSON int) string {
	var rez_json string
	rez_json = ""
	rez_json = string(34) + datJSON + string(34) + ": "
	if tpJSON == 1 {
		rez_json = rez_json + string(34) + valJSON + string(34)
	} else {
		rez_json = rez_json + valJSON
	}
	return rez_json
}

//------------------------------------------------------------------------
func convertIntVal(strColumn string) int {
	var i_convert int
	i_convert = 0
	j, err := strconv.Atoi(strColumn)
	if err != nil {
		defer recoveryAppFunction()
		fmt.Println(getDT(), err)
	} else {
		i_convert = j
	}
	return i_convert
}

//--------------------------------------------------------------------------
func convertInt64Val(strColumn string) int64 {
	var i_convert int64
	i_convert = 0
	j, err := strconv.ParseInt(strColumn, 10, 64)
	if err != nil {
		fmt.Println(getDT(), err)
		defer recoveryAppFunction()
	}
	i_convert = j
	return i_convert
}

//--------------------------------------------------------------------------
func dec2hex(dec int) string {
	//hexStr := dec * 255 / 100
	//return fmt.Sprintf("%02x", hexStr)
	hexStr := fmt.Sprintf("%02x", dec)
	return strings.ToUpper(hexStr)

}

//--------------------------------------------------------------------------
func hex2dec(hex string) int64 {
	dec, err := strconv.ParseInt("0x"+hex, 0, 16)
	if err != nil {
		fmt.Println(getDT(), err)
		defer recoveryAppFunction()
	}
	dec = dec * 100 / 255
	return dec

}

//---------------- convert -----------------------
func hex2int(hexStr string) int {
	result, err := strconv.ParseUint(hexStr, 16, 32)
	if err != nil {
		return -1
	}
	return int(result)
}

//--------------------------------------------------------------------------
func string_int32(snumeric string, defnimeric int64) int64 {
	var i_result int64
	i_result = defnimeric
	i, err := strconv.ParseInt(snumeric, 10, 64)
	if err != nil {
		defer recoveryAppFunction()
		fmt.Println(getDT(), "Error Convertion type:"+snumeric, err)
		//panic(err)
	}
	i_result = i
	return i_result
}

//------------------------------------------------------------------------------
func strToFloat64(strFloat string) float64 {
	value, err := strconv.ParseFloat(strFloat, 64)
	if err != nil {
		defer recoveryAppFunction()
		return 0
	}
	return float64(value)
}

//=================== Slice Utils =========================================
func removeFromSliceInt(sliceSource []int, slicePos int) []int {
	sliceListTemp := make([]int, 0)
	for i := 0; i < len(sliceSource); i++ {
		if i != slicePos {
			sliceListTemp = append(sliceListTemp, sliceSource[i])
		}
	}
	return sliceListTemp
}

//------------------------------------------------------------------------------
func removeFromSliceString(sliceSource []string, slicePos int) []string {
	sliceListTemp := make([]string, 0)
	for i := 0; i < len(sliceSource); i++ {
		if i != slicePos {
			sliceListTemp = append(sliceListTemp, sliceSource[i])
		}
	}
	return sliceListTemp
}

//------------------------------------------------------------------------------
func removeFromSliceBool(sliceSource []bool, slicePos int) []bool {
	sliceListTemp := make([]bool, 0)
	for i := 0; i < len(sliceSource); i++ {
		if i != slicePos {
			sliceListTemp = append(sliceListTemp, sliceSource[i])
		}
	}
	return sliceListTemp
}

//----------------------------------------------------------------------------
func getGBRuin(userid string) bool {
	user_iden := convertIntVal(userid)
	for i := 0; i < len(gbrListID); i++ {
		if gbrListID[i] == user_iden {
			return true
		}
	}
	return false
}

//----------------------------------------------------------------------------
func getGBRposition(userid string) int {
	user_iden := convertIntVal(userid)
	for i := 0; i < len(gbrListID); i++ {
		if gbrListID[i] == user_iden {
			return i
		}
	}
	return -1
}

//=================== String functions =====================================//
func doReplaceStr(whole_String, old_Str, new_Str string) string {
	var s_result string
	s_result = whole_String
	for ok := true; ok; ok = (strings.Index(s_result, old_Str) > 0) {
		s_result = strings.Replace(s_result, old_Str, new_Str, -1)
	}
	return s_result
}

//---------------------------------------------------------------------------
func dbQuatedString(queryString string) string {
	return string(39) + queryString + string(39)
}

//---------------------------------------------------------------------------
func jsonQuatedString(queryString string) string {
	return string(34) + queryString + string(34)
}

//======================= JSON FUNCTIONS =====================================
func checkValidJson(strChecked string) bool {
	var is_valid bool
	is_valid = true
	if len(strChecked) < 16 {
		return false
	}
	if strings.Index(strChecked, "{") < 0 {
		return false
	}
	if strings.Index(strChecked, "id") < 1 {
		return false
	}

	if strings.Index(strChecked, "cmnd") < 1 {
		return false
	}
	if strings.Index(strChecked, "name") < 1 {
		return false
	}
	if strings.Index(strChecked, "param") < 1 {
		return false
	}
	if strings.Index(strChecked, "}") < 1 {
		return false
	}
	return is_valid
}

//===================== GET DATABASE TYPES =====================================
func getRegions() []string {
	var (
		regionList []string
		s_check    string
	)
	s_check = ""
	regionList = make([]string, 0)
	for i := 0; i < len(gbrListREG); i++ {
		if strings.Index(s_check, gbrListREG[i]) < 0 {
			s_check += gbrListREG[i] + ";"
			//fmt.Println(gbrListREG[i])
			regionList = append(regionList, gbrListREG[i])
		}
	}
	sort.Strings(regionList)
	return regionList
}

//------------------------------------------------------------------------------
func prepareStringForGBR(strGBRorREG string) string {
	s_check := strings.ReplaceAll(strGBRorREG, string(34), " ")
	s_check = strings.ReplaceAll(s_check, string(13), " ")
	s_check = strings.ReplaceAll(s_check, string(10), " ")
	return strings.TrimSpace(s_check)
}

//------------------------------------------------------------------------------
func getGBRlist(isGBRlist int) string {
	var (
		jsonGbrList []string
		s_json      string
		i_len       int
	)
	jsonGbrList = make([]string, 0)
	s_json = ""
	i_len = 0
	if isGBRlist == 0 { //REGION LIST
		jsonGbrList = getRegions()
		i_len = len(jsonGbrList)
	} else if isGBRlist == 1 { //GBR LIST
		i_len = len(gbrListID)
		jsonGbrList = gbrListNUM
	} else { //PERSONAL LIST
		//RESERVED
	}

	if i_len > 0 {
		for i := 0; i < len(jsonGbrList); i++ {
			s_json += "{" //+ string(0x0D) + string(0x0A)
			if isGBRlist == 1 {
				s_json += getQuatedJSON("id", strconv.Itoa(gbrListID[i]), 1) + "," //+ string(0x0D) + string(0x0A)
			} else {
				s_json += getQuatedJSON("id", strconv.Itoa(i+1), 1) + "," //+ string(0x0D) + string(0x0A)
			}
			s_json += getQuatedJSON("name", prepareStringForGBR(jsonGbrList[i]), 1) //+ string(0x0D) + string(0x0A)
			s_json += "},"
		}
	} else {
		s_json = "{}"
	}

	return s_json
}

//------------------------------------------------------------------------------
func updateGBRstatus(update_GBR_id, update_GEO, update_REPORT string, update_GBR_status int) {
	/*
		update_GBR_status
		0 - Login USER
		1 - Start Alarm
		2 - Point Alarm
		3 - Break Alarm
		4 - Stop Alarm
		5 - Geo Alarm
	*/

	i_pos := getGBRposition(update_GBR_id)
	if i_pos > -1 {
		if update_GBR_status == 1 {
			gbrAlarmReceived[i_pos] = true //SET RECEIVED ALARM TRUE
			gbrAlarmWait[i_pos] = 0        //NOT WAIT NEXT CONNECT
		} else if update_GBR_status == 2 {
			gbrAlarmReceived[i_pos] = true //SET RECEIVED ALARM TRUE
			gbrAlarmWait[i_pos] = 0        //NOT WAIT NEXT CONNECT
		} else if update_GBR_status == 3 {
			gbrAlarmReceived[i_pos] = true //SET RECEIVED ALARM TRUE
			gbrAlarmType[i_pos] = 2        //SET CANCELLED ALARM, NOT SEND AGAIN
			gbrAlarmWait[i_pos] = 0        //NOT WAIT NEXT CONNECT
		} else if update_GBR_status == 4 {
			gbrAlarmLast[i_pos] = "" //Clear Last Alarm ID
			gbrAlarmPult[i_pos] = ""
			gbrAlarmName[i_pos] = ""
			gbrAlarmReceived[i_pos] = false
			gbrAlarmSend[i_pos] = false
			gbrAlarmType[i_pos] = 0 //CLEAR STATUS ALARM TYPE
			gbrAlarmType[i_pos] = 0
			gbrAlarmWait[i_pos] = 0
			gbrAlarmCard[i_pos] = ""
		} else if update_GBR_status == 5 {
			//NOTHING
		}

	}

	if update_GBR_status == 5 {

		data := []byte(`{
		"status":"alarmstart",
		"param":"X0Y0",
		"id":"123456"
	}`)
		var testJsonUnmarshal sendStatusOfAlarm
		if err := json.Unmarshal(data, &testJsonUnmarshal); err != nil {
			panic(err)
		}
		r := bytes.NewReader(data)
		resp, err := http.Post("http://api-cs.ohholding.com.ua/api/set-status?status="+testJsonUnmarshal.Status+"&param="+testJsonUnmarshal.Param+"&id="+update_GBR_id, "application/json", r)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v %v", err, resp)
	}

	if update_GBR_status == 4 {

		data := []byte(`{
		"status":"alarmstop",
		"param":"X0Y0",
		"id":"123456"
	}`)
		var testJsonUnmarshal sendStatusOfAlarm
		if err := json.Unmarshal(data, &testJsonUnmarshal); err != nil {
			panic(err)
		}
		r := bytes.NewReader(data)
		resp, err := http.Post("http://api-cs.ohholding.com.ua/api/set-status?status="+testJsonUnmarshal.Status+"&param="+update_REPORT+"&id="+update_GBR_id, "application/json", r)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v %v", err, resp)
	}

	if update_GBR_status == 3 {

		data := []byte(`{
		"status":"alarmbreak",
		"param":"X0Y0",
		"id":"123456"
	}`)

		var testJsonUnmarshal sendStatusOfAlarm
		if err := json.Unmarshal(data, &testJsonUnmarshal); err != nil {
			panic(err)
		}
		r := bytes.NewReader(data)
		resp, err := http.Post("http://api-cs.ohholding.com.ua/api/set-status?status="+testJsonUnmarshal.Status+"&param="+update_REPORT+"&id="+update_GBR_id, "application/json", r)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v %v", err, resp)
	}

	if update_GBR_status == 2 {
		data := []byte(`{
		"status":"alarmpoint",
		"param":"X0Y0",
		"id":"123456"
	}`)

		var testJsonUnmarshal sendStatusOfAlarm
		if err := json.Unmarshal(data, &testJsonUnmarshal); err != nil {
			panic(err)
		}
		r := bytes.NewReader(data)
		resp, err := http.Post("http://api-cs.ohholding.com.ua/api/set-status?status="+testJsonUnmarshal.Status+"&param="+update_REPORT+"&id="+update_GBR_id, "application/json", r)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v %v", err, resp)
	}
}

//=====================AUTO UPDATE MODULES =====================================
func checkUpdateView() {
	if uptList {
		updateSockList()
	}
	uptList = false
}
