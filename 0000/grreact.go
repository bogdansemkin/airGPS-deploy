package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var (
	userids             []string
	alarmstartuser      []string
	logResult           string
	checker             int
	gbrWorkersArray     []int
	objectPultArray     []string
	convertId           int
	gbrPult             []string
	gbrJsonRawEscaped   json.RawMessage
	gbrJsonRawUnescaped json.RawMessage
	getWeb              *websocket.Conn
	userid              []string
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//========================= get JSON func =======================================
type StringTable []string

func (st StringTable) Get(i int) string {
	if i < 0 || i > len(st) {
		return ""
	}
	return st[i]
}

func (st StringTable) GetIndex(i int) int {
	return i
}

func getJSON(url string, target interface{}) error {
	return json.NewDecoder(bytes.NewBufferString(url)).Decode(target)
}

func JSONget(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func bla(userid string) bool {
	mytable := StringTable{
		1:  "71",
		2:  "72",
		3:  "73",
		4:  "74",
		5:  "75",
		6:  "78",
		7:  "79",
		8:  "80",
		9:  "81",
		10: "82",
		11: "83",
		12: "84",
		13: "85",
		14: "86",
		15: "88",
		16: "89",
		17: "92",
		18: "Байкал 1",
		19: "Байкал 2",
		20: "Байкал 3",
		21: "Байкал 4",
		22: "Байкал 5",
		23: "Байкал 6",
		24: "Байкал 7",
		25: "Днепр 1",
		26: "Днепр 2",
		27: "Днепр 3",
		28: "Днепр 4",
		29: "Днепр 5",
		30: "Днепр 6",
		31: "Каменск 1",
		32: "Каменск 2",
		35: "Кривбас 1",
		36: "Кривбас 2",
		37: "Кривбас 3",
		38: "Кривбас 4",
		39: "Кривбас 7",
		40: "Кривбас 6",
		43: "Сокол",
		49: "Львов 1",
		50: "Львов 2",
		51: "Львов 3",
		52: "Львов 4",
		60: "76",
		65: "Мариуполь-1",
		66: "Мариуполь-2",
		67: "Мариуполь-4",
		72: "87",
		74: "Покровск",
		76: "Энергодар",
		77: "77",
		80: "Павлоград 1",
		81: "Павлоград 2",
		82: "Каменск 3",
		83: "Каменск 4",
		84: "Днепр 7",
		85: "Байкал 8",
		86: "Львов 5",
		88: "91",
		89: "Каменск 5",
		90: "Мариуполь-3",
		91: "МОТО - 1",
	}
	intvar, _ := strconv.Atoi(userid)

	if mytable.Get(intvar) == "" {
		return false
	}
	return true
}

func dbQuatedStringt(unQuated string) string { //&quot;
	return string(39) + doReplaceStr(unQuated, "&quot;", "") + string(39)
}

func getFotoFile(cardBase *CardBase) string {
	/*	cardFoto := new(CardBase)
	 */s_json := ""
	for i := 0; i < len(cardBase.CARD_FILES); i++ {
		s_json += "{" + string(34) + "url" + string(34) + ":" + string(34) + cardBase.CARD_FILES[i] + string(34) + "}" + hasNext(len(cardBase.CARD_FILES), i)
	}
	return s_json
}

func hasNext(baseLen int, i int) string {
	if i >= 0 && baseLen-1 != i {
		return ","
	}
	return ""
}
func hasNextzone(baseLen int, i int) string {
	/*	fmt.Println("-------------------------")
		fmt.Println("METHOD I:", i)
		fmt.Println("METHOD BASELEN:",baseLen)
		fmt.Println("-------------------------")*/

	if i >= 0 && baseLen != i {
		return ","
	}
	return ""
}
func userlist(url string) string {
	incomingJSON := new(CardBase)
	s_json := ""
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	json.Unmarshal(body, &incomingJSON)
	for _, people := range incomingJSON.CARD_PEOPLE {
		s2 := string(34)
		s_json = s2 + "userlist" + s2 + ":" + "[" + "{" + s2 + "name" + s2 + ":" + s2 + people.MAN_NAME + s2 + "," + "num" + ":" + s2 + people.MAN_NUM + s2 + "," + s2 + "tel" + s2 + ":" + s2 + people.MAN_PHONE + s2 + "}" + "]"
	}

	return s_json
}

func zonelist(url string) string {
	incomingJSON := new(CardBase)
	s_json := ""
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	json.Unmarshal(body, &incomingJSON)
	i := 0
	for _, zone := range incomingJSON.CARD_SHLEYF {

		i++
		/*if i == len(incomingJSON.CARD_SHLEYF){
			i -= 1
		}*/
		/*		fmt.Println("-----------------------------------:", i)
		 */s2 := string(34)
		s_json += "{" + s2 + "name" + s2 + ":" + s2 + zone.ZONE_NAME + s2 + "," +
			s2 + "num" + s2 + ":" + s2 + zone.ZONE_NUM + s2 + "," + s2 + "tel" + s2 + ":" + s2 + zone.ZONE_PLACE + s2 + "}" + hasNextzone(len(incomingJSON.CARD_SHLEYF), i)
	}

	return s_json
}

func linearSearch(gbrArray []int, key int) bool {
	checker = -1
	for _, item := range gbrArray {
		checker++
		if item == key {
			return true
		}
	}
	return false
}

func gbrNowActiveWorkersArray(incomingJSON []gbrNowActiveWorkers) {
	for _, val := range incomingJSON {
		var newGbrNumber, _ = strconv.Atoi(val.GbrNumber)
		gbrWorkersArray = append(gbrWorkersArray, newGbrNumber)
	}
}

func comparingPultGbr(incomingJSON []gbrNowActiveWorkers) {
	for _, val := range incomingJSON {
		objectPultArray = append(objectPultArray, val.ObjectNumberPult)
	}
}

func runJSON(incomingJSON []gbrNowActiveWorkers, userid []string, conn *websocket.Conn) {
	gbrArray := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91}
	cardBase := new(CardBase)
	j := 0
	for _, val := range incomingJSON {
		comparingPultGbr(incomingJSON)
		gbrNowActiveWorkersArray(incomingJSON)
		fmt.Println("CONVERTED id IN RUNJSON: ", convertId)
		if len(userid) != 0 && j < len(userid) {
			secondUserId, _ := strconv.Atoi(userid[j])
			//var blablaf = secondUserId

			if val.ObjectNumberPult != "" {
				convertId = gbrArray[secondUserId]
				convertedGbrId, _ := strconv.Atoi(val.GbrNumber)
				println("=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!!=====")
				fmt.Println("Cheker: ", checker)
				fmt.Println("convertedID: ", convertId)
				fmt.Println("ConvertedGBRId: ", convertedGbrId)
				println("=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!=!!=====")
				fmt.Println(gbrWorkersArray)
				if linearSearch(gbrWorkersArray, convertId) == true {
					objectPultString := objectPultArray[checker]
					fmt.Println("OBJECTPULTSTRING::::::::::::::::::" + objectPultString)
					JSONget("http://api-cs.ohholding.com.ua/object_cart/"+objectPultString+"/get", cardBase)
					url := "http://api-cs.ohholding.com.ua/object_cart/" + objectPultString + "/get"
					s2 := string(34)
					s_json := "{" + s2 + "obinfo" + s2 + ":" + "[" + "{" + s2 + "id" + s2 + ":" + s2 + cardBase.ID + s2 + "," + s2 + "lat" + s2 + ":" + s2 +
						cardBase.CARD_LAT + s2 + "," + s2 + "lon" + s2 + ":" + s2 + cardBase.CARD_LON + s2 + "," + s2 + "obadr" + s2 + ":" + s2 +
						cardBase.CARD_ADRES + s2 + "," + s2 + "obname" + s2 + ":" + s2 + cardBase.CARD_NAME + s2 + "," + s2 + "pult" + s2 + ":" + s2 +
						cardBase.CARD_PULTNUM + s2 + "," + s2 + "details" + s2 + ":" + s2 + cardBase.CARD_WAYMARK + s2 + "}" + "]"
					s_json += ","
					s_json += userlist(url)
					s_json += ","
					s_json += s2 + "zonelist" + s2 + ":" + "["
					s_json += zonelist(url)
					s_json += "]"
					s_json += ","
					s_json += s2 + "eventlist" + s2 + ":" + "[" + "{" + "}" + "]"
					s_json += ","
					s_json += s2 + "imagelist" + s2 + ":" + "[" + getFotoFile(cardBase) + "]"
					s_json += "}"
					fmt.Println(s_json)
					s_jsonbyte := []byte(s_json)
					conn.WriteMessage(websocket.TextMessage, s_jsonbyte)
				}
			} else if val.ObjectNumberPult == "" {
				convertId, _ := strconv.Atoi(userid[j])
				convertedGbrId, _ := strconv.Atoi(val.GbrNumber)
				if gbrArray[convertId] == convertedGbrId {
					s_json := "{" + string(34) + "no_pult" + string(34) + ":" + "[" + "{" +
						string(34) + "f_object_adress" + string(34) + ":" + string(34) + val.ObjectAdress + string(34) + "," +
						string(34) + "f_object_name" + string(34) + ":" + string(34) + val.ObjectName + string(34) + "," +
						string(34) + "f_region" + string(34) + ":" + string(34) + val.Region + string(34) + "}" + "]" + "}"
					fmt.Println("in case PULT == 0: " + s_json)
					s_jsonbyte := []byte(s_json)
					conn.WriteMessage(websocket.TextMessage, s_jsonbyte)
				}
			}
			j++
			//	k ++
		}
	}
}

//========================= MAIN LOGIC =======================================
func decodeGpsJson(jsonIncoming string, conn *websocket.Conn) string {
	var (
		gbrlistSout gbrListAll
		airDecoding AirQuery
		strJSON     []byte
		//		i_con       int
		js_result string
		js_iden   string
		js_cmnd   string
		js_param  string
		js_name   string
	)
	getJSON("http://api-cs.ohholding.com.ua/gbr_list/get", &gbrlistSout)
	js_result = "{" + string(0x0D) + string(0x0A)
	js_result += getQuatedJSON("param", "Status error", 1) + string(0x0D) + string(0x0A)
	js_result += "}" + string(0x0D) + string(0x0A)

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

		fmt.Println(js_name)

		getJSON("http://api-cs.ohholding.com.ua/gbr_list/get", &gbrlistSout)
		fmt.Println(gbrlistSout)
		b := []byte(`{"gbrlist":[{"id_gbr":1,"region":"\u041a\u0438\u0435\u0432","number":"71"},{"id_gbr":2,"region":"\u041a\u0438\u0435\u0432","number":"72"},{"id_gbr":3,"region":"\u041a\u0438\u0435\u0432","number":"73"},{"id_gbr":4,"region":"\u041a\u0438\u0435\u0432","number":"74"},{"id_gbr":5,"region":"\u041a\u0438\u0435\u0432","number":"75"},{"id_gbr":6,"region":"\u041a\u0438\u0435\u0432","number":"78"},{"id_gbr":7,"region":"\u041a\u0438\u0435\u0432","number":"79"},{"id_gbr":8,"region":"\u041a\u0438\u0435\u0432","number":"80"},{"id_gbr":9,"region":"\u041a\u0438\u0435\u0432","number":"81"},{"id_gbr":10,"region":"\u041a\u0438\u0435\u0432","number":"82"},{"id_gbr":11,"region":"\u041a\u0438\u0435\u0432","number":"83"},{"id_gbr":12,"region":"\u041a\u0438\u0435\u0432","number":"84"},{"id_gbr":13,"region":"\u041a\u0438\u0435\u0432","number":"85"},{"id_gbr":14,"region":"\u041a\u0438\u0435\u0432","number":"86"},{"id_gbr":15,"region":"\u041a\u0438\u0435\u0432","number":"88"},{"id_gbr":16,"region":"\u041a\u0438\u0435\u0432","number":"89"},{"id_gbr":17,"region":"\u041a\u0438\u0435\u0432","number":"92"},{"id_gbr":18,"region":"\u0417\u0430\u043f\u043e\u0440\u043e\u0436\u044c\u0435","number":"\u0411\u0430\u0439\u043a\u0430\u043b 1"},{"id_gbr":19,"region":"\u0417\u0430\u043f\u043e\u0440\u043e\u0436\u044c\u0435","number":"\u0411\u0430\u0439\u043a\u0430\u043b 2"},{"id_gbr":20,"region":"\u0417\u0430\u043f\u043e\u0440\u043e\u0436\u044c\u0435","number":"\u0411\u0430\u0439\u043a\u0430\u043b 3"},{"id_gbr":21,"region":"\u0417\u0430\u043f\u043e\u0440\u043e\u0436\u044c\u0435","number":"\u0411\u0430\u0439\u043a\u0430\u043b 4"},{"id_gbr":22,"region":"\u0417\u0430\u043f\u043e\u0440\u043e\u0436\u044c\u0435","number":"\u0411\u0430\u0439\u043a\u0430\u043b 5"},{"id_gbr":23,"region":"\u0417\u0430\u043f\u043e\u0440\u043e\u0436\u044c\u0435","number":"\u0411\u0430\u0439\u043a\u0430\u043b 6"},{"id_gbr":24,"region":"\u0417\u0430\u043f\u043e\u0440\u043e\u0436\u044c\u0435","number":"\u0411\u0430\u0439\u043a\u0430\u043b 7"},{"id_gbr":25,"region":"\u0414\u043d\u0435\u043f\u0440","number":"\u0414\u043d\u0435\u043f\u0440 1"},{"id_gbr":26,"region":"\u0414\u043d\u0435\u043f\u0440","number":"\u0414\u043d\u0435\u043f\u0440 2"},{"id_gbr":27,"region":"\u0414\u043d\u0435\u043f\u0440","number":"\u0414\u043d\u0435\u043f\u0440 3"},{"id_gbr":28,"region":"\u0414\u043d\u0435\u043f\u0440","number":"\u0414\u043d\u0435\u043f\u0440 4"},{"id_gbr":29,"region":"\u0414\u043d\u0435\u043f\u0440","number":"\u0414\u043d\u0435\u043f\u0440 5"},{"id_gbr":30,"region":"\u0414\u043d\u0435\u043f\u0440","number":"\u0414\u043d\u0435\u043f\u0440 6"},{"id_gbr":31,"region":"\u041a\u0430\u043c\u0435\u043d\u0441\u043a","number":"\u041a\u0430\u043c\u0435\u043d\u0441\u043a 1"},{"id_gbr":32,"region":"\u041a\u0430\u043c\u0435\u043d\u0441\u043a","number":"\u041a\u0430\u043c\u0435\u043d\u0441\u043a 2"},{"id_gbr":35,"region":"\u041a\u0440\u0438\u0432\u043e\u0439 \u0420\u043e\u0433","number":"\u041a\u0440\u0438\u0432\u0431\u0430\u0441 1"},{"id_gbr":36,"region":"\u041a\u0440\u0438\u0432\u043e\u0439 \u0420\u043e\u0433","number":"\u041a\u0440\u0438\u0432\u0431\u0430\u0441 2\r\n"},{"id_gbr":37,"region":"\u041a\u0440\u0438\u0432\u043e\u0439 \u0420\u043e\u0433","number":"\u041a\u0440\u0438\u0432\u0431\u0430\u0441 3"},{"id_gbr":38,"region":"\u041a\u0440\u0438\u0432\u043e\u0439 \u0420\u043e\u0433","number":"\u041a\u0440\u0438\u0432\u0431\u0430\u0441 4"},{"id_gbr":39,"region":"\u041a\u0440\u0438\u0432\u043e\u0439 \u0420\u043e\u0433","number":"\u041a\u0440\u0438\u0432\u0431\u0430\u0441 7"},{"id_gbr":40,"region":"\u041a\u0440\u0438\u0432\u043e\u0439 \u0420\u043e\u0433","number":"\u041a\u0440\u0438\u0432\u0431\u0430\u0441 6"},{"id_gbr":43,"region":"\u0414\u043e\u0431\u0440\u043e\u043f\u043e\u043b\u044c\u0435","number":"\u0421\u043e\u043a\u043e\u043b"},{"id_gbr":49,"region":"\u041b\u044c\u0432\u043e\u0432","number":"\u041b\u044c\u0432\u043e\u0432 1"},{"id_gbr":50,"region":"\u041b\u044c\u0432\u043e\u0432","number":"\u041b\u044c\u0432\u043e\u0432 2"},{"id_gbr":51,"region":"\u041b\u044c\u0432\u043e\u0432","number":"\u041b\u044c\u0432\u043e\u0432 3"},{"id_gbr":52,"region":"\u041b\u044c\u0432\u043e\u0432","number":"\u041b\u044c\u0432\u043e\u0432 4"},{"id_gbr":60,"region":"\u041a\u0438\u0435\u0432","number":"76"},{"id_gbr":65,"region":"\u041c\u0430\u0440\u0438\u0443\u043f\u043e\u043b\u044c","number":"\u041c\u0430\u0440\u0438\u0443\u043f\u043e\u043b\u044c-1"},{"id_gbr":66,"region":"\u041c\u0430\u0440\u0438\u0443\u043f\u043e\u043b\u044c","number":"\u041c\u0430\u0440\u0438\u0443\u043f\u043e\u043b\u044c-2"},{"id_gbr":67,"region":"\u041c\u0430\u0440\u0438\u0443\u043f\u043e\u043b\u044c","number":"\u041c\u0430\u0440\u0438\u0443\u043f\u043e\u043b\u044c-4"},{"id_gbr":72,"region":"\u041a\u0438\u0435\u0432","number":"87"},{"id_gbr":74,"region":"\u041f\u043e\u043a\u0440\u043e\u0432\u0441\u043a\u043e\u0435","number":"\u041f\u043e\u043a\u0440\u043e\u0432\u0441\u043a"},{"id_gbr":76,"region":"\u042d\u043d\u0435\u0440\u0433\u043e\u0434\u0430\u0440","number":"\u042d\u043d\u0435\u0440\u0433\u043e\u0434\u0430\u0440"},{"id_gbr":77,"region":"\u041a\u0438\u0435\u0432","number":"77"},{"id_gbr":80,"region":"\u041f\u0430\u0432\u043b\u043e\u0433\u0440\u0430\u0434","number":"\u041f\u0430\u0432\u043b\u043e\u0433\u0440\u0430\u0434 1"},{"id_gbr":81,"region":"\u041f\u0430\u0432\u043b\u043e\u0433\u0440\u0430\u0434","number":"\u041f\u0430\u0432\u043b\u043e\u0433\u0440\u0430\u0434 2"},{"id_gbr":82,"region":"\u041a\u0430\u043c\u0435\u043d\u0441\u043a","number":"\u041a\u0430\u043c\u0435\u043d\u0441\u043a 3"},{"id_gbr":83,"region":"\u041a\u0430\u043c\u0435\u043d\u0441\u043a","number":"\u041a\u0430\u043c\u0435\u043d\u0441\u043a 4"},{"id_gbr":84,"region":"\u0414\u043d\u0435\u043f\u0440","number":"\u0414\u043d\u0435\u043f\u0440 7"},{"id_gbr":85,"region":"\u0417\u0430\u043f\u043e\u0440\u043e\u0436\u044c\u0435","number":"\u0411\u0430\u0439\u043a\u0430\u043b 8"},{"id_gbr":86,"region":"\u041b\u044c\u0432\u043e\u0432","number":"\u041b\u044c\u0432\u043e\u0432 5"},{"id_gbr":88,"region":"\u041a\u0438\u0435\u0432","number":"91"},{"id_gbr":89,"region":"\u041a\u0430\u043c\u0435\u043d\u0441\u043a","number":"\u041a\u0430\u043c\u0435\u043d\u0441\u043a 5"},{"id_gbr":90,"region":"\u041c\u0430\u0440\u0438\u0443\u043f\u043e\u043b\u044c","number":"\u041c\u0430\u0440\u0438\u0443\u043f\u043e\u043b\u044c-3"},{"id_gbr":91,"region":"\u041a\u0438\u0435\u0432","number":"\u041c\u041e\u0422\u041e - 1"}]}`)
		/*go sleepinGoopher(startpoint, stoppoint, js_iden, conn, &waitgroup)*/
		switch js_cmnd {
		case "start": //First start
			//	js_result = startGBR(js_iden, js_name, js_param, getSocketIndex(conn))
			json.Unmarshal(b, &gbrlistSout)
			err = conn.WriteMessage(websocket.TextMessage, b)
		case "login": //Loging for user
			js_result = logGBR(js_iden, js_name, js_param, getSocketIndex(conn))
			userids = append(userids, js_iden)
			fmt.Println("===============================")
			fmt.Println(userids)
			fmt.Println("===============================")
			getWebsocket(conn, userids)
		case "connect":
			fmt.Println("In connect case")
			message := []byte("{\"cmnd\":\"connect\",\"id\":\"8\",\"name\":\"-1\",\"param\":\"Connect_OK\"}")
			err = conn.WriteMessage(websocket.TextMessage, message)
			fmt.Println("Successfully connected....")
			//TODO comparation json files
			//	case "alarmlist": //Get alarm list
			//		js_result = getAlarms(js_iden, js_name, js_param)
		case "alarmget": //Receive alarm

		case "alarmstart": //GBR starts trip
			js_result = procAlarm(js_iden, js_cmnd, js_name, js_param)

			alarmstartuser = userids

			for i := 0; i < len(userid); i++ {
				userids[i] = "0"
			}
			fmt.Println("IN ALARMSTART CASE")
			fmt.Println("userID-S: ", userids)

		case "alarmpoint": //GBR at point
			js_result = procAlarm(js_iden, js_cmnd, js_name, js_param)
		case "alarmbreak": //Problem with GBR
			js_result = procAlarm(js_iden, js_cmnd, js_name, js_param)
		case "alarmstop": //Set reaction
			js_result = procAlarm(js_iden, js_cmnd, js_name, js_param)
			userids = alarmstartuser
		case "alarminfo": //Read updates

		default:
			js_result = setUnknown(js_iden, js_name, js_cmnd)
		}
	}
	return js_result
}

//------------------------------------------------------------------------------
func sendUpdator(userid int) string {
	s_json := "{" + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("id", strconv.Itoa(websock_uin_users[userid]), 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("cmnd", "update", 1) + string(0x0D) + string(0x0A)
	s_json = s_json + "}" + string(0x0D) + string(0x0A)
	return s_json
}

//------------------------------------------------------------------------------
func setUnknown(userid, js_name, js_command string) string {
	fmt.Println(getDT(), "Command unknown"+js_command)
	s_json := "{" + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("id", userid, 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("name", js_name, 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("cmnd", js_command, 1) + "," + string(0x0D) + string(0x0A)
	s_json = s_json + getQuatedJSON("param", "STATUS_ERROR", 1) + string(0x0D) + string(0x0A)
	s_json = s_json + "}" + string(0x0D) + string(0x0A)

	return s_json
}

//------------------------------------------------------------------------------

func logGBR(userid, js_name, js_param string, conPosition int) string {
	s_sql := "-2"
	s_sql += "(CONSTVISIB = 0) AND (CONSTKIND = 4) AND (IDCONST = " + userid
	s_sql += ") LIMIT 1"
	//TODO remake valid method
	gbrvalid := bla(userid)

	s_sql = "SELECT IDPERS,FIOPERS,PAROL FROM personality WHERE IDPERS=" + dbQuatedString(js_name)
	s_json := ""

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
			logResult = "LOG_OK"
			fmt.Println("LOG_RESULY: " + logResult)

			if conPosition >= 0 && conPosition < websock_addr_counter {
				s_tocken = websock_send_device[conPosition]
				websock_uin_users[conPosition] = convertIntVal(userid)
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

func getWebsocket(conn *websocket.Conn, useridd []string) {
	getWeb = conn
	userid = useridd
}

//------------------------------------------------------------------------------
func procAlarmList() {
	for i := 0; i < len(alarmList_ID); i++ { //READ ALARMS
		if alarmList_WAS_SEND[i] == false { //CHECK WAS SEND
			for j := 0; j < uptGBRlist; j++ { //SEARCH GBR
				if alarmList_GBR_NAME[i] == gbrListNUM[j] || alarmList_GBR_RESERVE[i] == gbrListNUM[j] { //GBR VALID
					isGBRcon := false
					for k := 0; k < websock_addr_counter; k++ { //SEARCH GBR IN CONNECT LIST
						if websock_uin_gbr[k] == "" {
							websock_uin_gbr[k] = getNameGBR(websock_uin_users[i])
						}
						if websock_uin_gbr[k] == gbrListNUM[j] { //GBR CONNECTED
							sendYeden(i, "")
						}
					} //SEARCH GBR FOR SEND PUSH
					if isGBRcon == false {

					} //END - SEARCH GBR IN CONNECT LIST
				} //END - GBR VALID
			} //END - SEARCH GBR
		} //END - CHECK WAS SEND
	} //END - READ ALARMS
}

//------------------------------------------------------------------------------
func checkUnsendAlarms() {

}
