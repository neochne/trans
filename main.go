package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"unicode/utf8"
	"crypto/sha256"
	"os"
	"time"
	"strconv"
	"strings"
	"encoding/hex"
	"encoding/json"
)

const (
	APP_ID = "4c67797e6790510c"
	APP_KEY = "OJDxTKImnx15Ld6bOo0Zy0e1rEZ4nYqT"
	BASE_URL = "https://openapi.youdao.com/api"
	SALT = "1x2r6y68"
	FROM = "auto"
	SIGN_TYPE = "v3"
)

// 字段名起始字母必须是大写，否则外部访问不了，即会解析不到值
type ydTransRst struct {
	Basic ydBasic
}

type ydBasic struct {
	Phonetic string
	UkPhonetic string `json:"uk-phonetic"`
	UsPhonetic string `json:"us-phonetic"`
	Explains []string
	Wfs []ydWfOut
}

type ydWfOut struct {
	Wf ydWf
}

type ydWf struct {
	Name string
	Value string
}

func main() {
    // Create http client
    client := &http.Client{}
    req, err := http.NewRequest("GET",BASE_URL,nil)
    if err != nil {
    	panic(err)
    }

    // Get and calc param
    q := os.Args[1]
    toStr := getTo(q)
    curTime := strconv.FormatInt(time.Now().Unix(),10)
    sign := generateSign(q,curTime)

    // Add query params
    query := req.URL.Query()
    query.Add("appKey",APP_ID)
    query.Add("signType",SIGN_TYPE)
    query.Add("from",FROM)
    query.Add("to",toStr)
    query.Add("salt",SALT)
    query.Add("curtime",curTime)
    query.Add("sign",sign)
    query.Add("q",q)
    req.URL.RawQuery = query.Encode()

    // Request
    resp, err := client.Do(req)
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
    }

    // Parse response
    // bodyStr := string(body)
    // fmt.Println(bodyStr)
    var transRst ydTransRst
    errUnmarshal := json.Unmarshal(body, &transRst)
    if errUnmarshal != nil {
         fmt.Println(errUnmarshal)
    }

    // Phonetic
    fmt.Println()
    fmt.Printf(" %s  音 [ %s ] 英 [ %s ]  美 [ %s ]",q,transRst.Basic.Phonetic,transRst.Basic.UkPhonetic,transRst.Basic.UsPhonetic)
    fmt.Println()
    fmt.Println()

	// Ws
    wfs := transRst.Basic.Wfs
    if len(wfs) > 0 {
        for _,wfOut := range wfs {
            name := wfOut.Wf.Name
            value := wfOut.Wf.Value
            nameLen := utf8.RuneCountInString(name)
            w := 0
            if nameLen == 2 {
            	w = 10
            } else if nameLen == 3{
            	w = 9
            } else if nameLen == 4 {
            	w = 8
            }else if nameLen == 5 {
            	w = 7
            }else if nameLen == 6 {
                w = 6
            }
            fmt.Printf(" %-" + strconv.Itoa(w) + "s : %s\n", name, value)
        }
    	fmt.Println()
    	fmt.Println()
    }

	// Explains
    for _,explain := range transRst.Basic.Explains {
        fmt.Println(" - " + explain)
    }
    fmt.Println()
}

func getTo(q string) string{
	if strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",q) {
		return "zh-CHS"
	} else {
		return "en"
	}
}

func generateSign(q string,curTime string) string{
	input := ""
	qLen := utf8.RuneCountInString(q)
	if qLen < 20 {
		input = q
	} else {
		input = q[:10] + string(qLen) + q[20:qLen]
	}	
	oriSign := APP_ID + input + SALT + curTime + APP_KEY
	h := sha256.New()
	h.Write([]byte(oriSign))
	x := hex.EncodeToString(h.Sum(nil))
	return x
}