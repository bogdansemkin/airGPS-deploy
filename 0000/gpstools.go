package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"log"

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

//------------------------------------------------------------------------------
func getObjGeneral(objUIN string, isKey bool) string {
	/*
		s_sql := "SELECT IDOB,OBNAME,OBADR,OBTEL,OHRANA,SVIZI,MOREINFO,INFODETAIL,"
		s_sql += "concat(PULTADR,'.0',AVATAR) AS PULTNUMBER,DOLGOTA,SHIROTA "
		s_sql += "FROM objectlist WHERE IDOB=" + objUIN + " LIMIT 1"

		db, err := sqlx.Connect("mysql", connect_db)
		s_json := ""
		if isKey == false {
			s_json += getQuatedJSON("obinfo", "[", 0) + string(0x0D) + string(0x0A)
		}

		if err != nil {
			defer recoveryAppFunction()
			fmt.Println(getDT(), "Error read value :"+s_sql, err)
			panic(err)
		}
		defer db.Close()
		rows, err := db.Query(s_sql)
		if err != nil {
			defer recoveryAppFunction()
			fmt.Println("Error read value :" + s_sql)
		}
		cols, err := rows.Columns()
		if err != nil {
			defer recoveryAppFunction()
			fmt.Println("Error read value :" + s_sql)
		}
		data := make(map[string]string)
		if rows.Next() {
			columns := make([]string, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}
			rows.Scan(columnPointers...)

			for i, colName := range cols {
				data[colName] = columns[i]
				if i == 0 { // 0 - Get IDOB
					s_json += "{" + string(0x0D) + string(0x0A)
					s_json += getQuatedJSON("id", columns[i], 1) + "," + string(0x0D) + string(0x0A)
				}
				if i == 1 { // 1 - OBNAME
					s_json += getQuatedJSON("obname", columns[i], 1) + "," + string(0x0D) + string(0x0A)
				}
				if i == 2 { // OBADR
					s_json += getQuatedJSON("obadr", columns[i], 1) + "," + string(0x0D) + string(0x0A)
				}
				if i == 3 { // OBTEL
					s_json += getQuatedJSON("obtel", columns[i], 1) + "," + string(0x0D) + string(0x0A)
				}
				if i == 4 { // OHRANA
					s_json += getQuatedJSON("status", columns[i], 1) + "," + string(0x0D) + string(0x0A)
				}
				if i == 5 { // SVIZI
					s_json += getQuatedJSON("con", columns[i], 1) + "," + string(0x0D) + string(0x0A)
				}
				if i == 6 { // MOREINFO
					s_json += getQuatedJSON("more", columns[i], 1) + "," + string(0x0D) + string(0x0A)
				}
				if i == 7 { // INFODETAIL
					s_json += getQuatedJSON("detail", columns[i], 1) + "," + string(0x0D) + string(0x0A)
				}
				if i == 8 { // PULTNUMBER
					s_json += getQuatedJSON("pult", columns[i], 1) + "," + string(0x0D) + string(0x0A)
				}
				if i == 9 { // LONGITUDE
					s_json += getQuatedJSON("lon", columns[i], 1) + "," + string(0x0D) + string(0x0A)
				}
				if i == 10 { // LATITUDE
					s_json += getQuatedJSON("lat", columns[i], 1) + string(0x0D) + string(0x0A)
					s_json += "}" + string(0x0D) + string(0x0A)
				}
			}
		}
		if isKey == false {
			s_json += "]"
		}
		if rows != nil {
			defer rows.Close()
		}
		return s_json
	*/
	return ""
}

//------------------------------------------------------------------------------
