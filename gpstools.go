package main

import (
	"fmt"
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
func recoveryAppFunction(troubleName string) {
	if recoveryMessage := recover(); recoveryMessage != nil {
		fmt.Println(recoveryMessage)
	}
	fmt.Println(getDT(), troubleName, "Application restored after Error...")
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
		defer recoveryAppFunction("convertIntVal Error " + strColumn)
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
		defer recoveryAppFunction("convertInt64Val Error " + strColumn)
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
		defer recoveryAppFunction("hex2dec Error " + hex)
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
		defer recoveryAppFunction("string_int32 Error " + snumeric)
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
		defer recoveryAppFunction("strToFloat64 Error " + strFloat)
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
func readGBRlist(isGBRlist int) string {
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
//[{"id_gbr":1,"region":"\u041a\u0438\u0435\u0432","number":"71"},{"id_gbr":2,"region":"\u041a\u0438\u0435\u0432","number":"72"},{"id_gbr":3,"region":"\u041a\u0438\u0435\u0432","number":"73"},{"id_gbr":4,"region":"\u041a\u0438\u0435\u0432","number":"74"},{"id_gbr":5,"region":"\u041a\u0438\u0435\u0432","number":"75"},
func getGBRlist() string {
	s_json := ""
	if len(gbrListID) > 0 {
		for i := 0; i < len(gbrListID); i++ {
			s_json += "{" + getQuatedJSON("id_gbr", strconv.Itoa(gbrListID[i]), 1) + "," +
				getQuatedJSON("region", prepareStringForGBR(gbrListREG[i]), 1) + "," +
				getQuatedJSON("number", prepareStringForGBR(gbrListNUM[i]), 1) + "},"
		}
	}

	return s_json
}

//------------------------------------------------------------------------------
func getConnectionPosition(alarmPosition int) int {
	for i := 0; i < len(websock_gbr_uin); i++ {
		if websock_gbr_uin[i] == gbrListID[alarmPosition] {
			return i
		}
	}
	return -1
}

//------------------------------------------------------------------------------
func updateGBRstatus(update_GBR_id, update_CMD, update_GEO string, update_GBR_status int) {
	/*
		update_GBR_status
		0 - Login USER
		1 - Start Alarm
		2 - Point Alarm
		3 - Break Alarm
		4 - Stop Alarm
		5 - Geo Alarm
	*/
	i_iden := convertIntVal(update_GBR_id)
	i_pos := getGBRposition(update_GBR_id)
	if i_pos > -1 && i_pos < len(gbrAlarmReceived) {
		i_iden = convertIntVal(gbrAlarmLast[i_pos])
		if update_GBR_status == 1 {
			gbrAlarmReceived[i_pos] = true //SET RECEIVED ALARM TRUE
			gbrAlarmWait[i_pos] = 0        //NOT WAIT NEXT CONNECT
			s_alarm := gbrAlarmLast[i_pos]
			for i := 0; i < len(gbrAlarmLast); i++ {
				if i_pos != i && s_alarm == gbrAlarmLast[i] {
					fmt.Println("Canceled on Status 1", gbrListID[i], gbrAlarmLast[i], gbrListNUM[i])
					gbrAlarmReceived[i] = false //SET RECEIVED ALARM TRUE
					gbrAlarmType[i] = 0         //SET CANCELLED ALARM, NOT SEND AGAIN
					gbrAlarmWait[i] = 0         //NOT WAIT NEXT CONNECT
					gbrAlarmCanceled[i] = 0
					gbrAlarmLast[i] = ""
					gbrAlarmCard[i] = ""
				}
			}
		} else if update_GBR_status == 2 {
			gbrAlarmReceived[i_pos] = true //SET RECEIVED ALARM TRUE
			gbrAlarmCanceled[i_pos] = 0
			gbrAlarmWait[i_pos] = 0 //NOT WAIT NEXT CONNECT
		} else if update_GBR_status == 3 {
			gbrAlarmReceived[i_pos] = false //SET RECEIVED ALARM TRUE
			gbrAlarmType[i_pos] = 2         //SET CANCELLED ALARM, NOT SEND AGAIN
			gbrAlarmWait[i_pos] = 0         //NOT WAIT NEXT CONNECT
			gbrAlarmCanceled[i_pos] = 30
			for i := 0; i < len(alarmList_ID); i++ {
				if strconv.Itoa(i_iden) == alarmList_ID[i] {
					alarmList_WAS_SEND[i] = false
					break
				}
			}
			//gbrAlarmLast[i_pos] = ""
			//gbrAlarmCard[i_pos] = ""
		} else if update_GBR_status == 4 {
			gbrAlarmLast[i_pos] = "" //Clear Last Alarm ID
			gbrAlarmPult[i_pos] = ""
			gbrAlarmName[i_pos] = ""
			gbrAlarmAddr[i_pos] = ""
			gbrAlarmReceived[i_pos] = false
			gbrAlarmSend[i_pos] = false
			gbrAlarmType[i_pos] = 0 //CLEAR STATUS ALARM TYPE
			gbrAlarmType[i_pos] = 0
			gbrAlarmWait[i_pos] = 0
			gbrAlarmCanceled[i_pos] = 0
			gbrAlarmCard[i_pos] = ""
		} else if update_GBR_status == 5 {
			//NOTHING
		}
	}
	s_update_GEO := update_GEO
	if update_GBR_status < 3 {
		if strings.Index(update_GEO, "X") != 0 || strings.Index(update_GEO, "Y") < 1 {
			s_update_GEO = "X0Y0"
		}
	}
	s_QUERY := "https://api-cs.ohholding.com.ua/api/set-status?status=" + update_CMD +
		"&param=" + s_update_GEO + "&id=" + strconv.Itoa(i_iden) //+ "application/json"
	fmt.Println("Post send", time.Now(), s_QUERY)
	postUpdateParams(s_QUERY)
}
//------------------------------------------------------------------------------
func getObjGeneral(objUIN int) string {
	s_json := ""
	fmt.Println("getObjGeneral", objUIN, gbrListID[objUIN], gbrAlarmLast[objUIN], gbrAlarmPult[objUIN], gbrAlarmName[objUIN])
	if len(gbrAlarmLast) > objUIN {
		for i := 0; i < len(alarmList_NUM_PULT); i++ {
			fmt.Println("getObjGeneral Search", alarmList_ID[i], gbrAlarmLast[objUIN])
			if alarmList_ID[i] == gbrAlarmLast[objUIN] {
				s_json = "{" + string(0x0D) + string(0x0A)
				s_json = s_json + getQuatedJSON("id", alarmList_ID[i], 1) + "," + string(0x0D) + string(0x0A)
				s_json = s_json + getQuatedJSON("name", alarmList_OBJ_NAME[i], 1) + "," + string(0x0D) + string(0x0A)
				s_json = s_json + getQuatedJSON("addr", alarmList_OBJ_ADDR[i], 1) + "," + string(0x0D) + string(0x0A)
				if len(alarmList_NUM_PULT[i]) >= 0 {
					s_json = s_json + getQuatedJSON("pult", alarmList_NUM_PULT[i], 1) + string(0x0D) + string(0x0A)
				} else {
					s_json = s_json + getQuatedJSON("pult", "no_pult", 1) + string(0x0D) + string(0x0A)
				}
				s_json = s_json + "}" + string(0x0D) + string(0x0A)
				//fmt.Println("getObjGeneral Find", s_json)
				return s_json
			}
		}
	}

	s_json = "{" + string(0x0D) + string(0x0A)
	if len(gbrAlarmLast[objUIN]) > 0 {
		s_json = s_json + getQuatedJSON("id", gbrAlarmLast[objUIN], 1) + "," + string(0x0D) + string(0x0A)
	} else {
		s_json = s_json + getQuatedJSON("id", strconv.Itoa(gbrListID[objUIN]), 1) + "," + string(0x0D) + string(0x0A)
	}

	s_json = s_json + getQuatedJSON("name", gbrAlarmName[objUIN], 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("addr", gbrAlarmAddr[objUIN], 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("pult", gbrAlarmPult[objUIN], 1) + string(0x0D) + string(0x0A)
	s_json = s_json + "}" + string(0x0D) + string(0x0A)
	return s_json
}

//=====================AUTO UPDATE MODULES =====================================
func checkUpdateView() {
	if uptList {
		updateSockList()
	}
	uptList = false
}
