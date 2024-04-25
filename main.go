package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"regexp"
)

// https://ai.youdao.com/DOCSIRMA/html/trans/api/wbfy/index.html
// user: 15803
// pass: .....

const (
	appID    = "4c67797e6790510c"
	appKey   = "OJDxTKImnx15Ld6bOo0Zy0e1rEZ4nYqT"
	baseURL  = "https://openapi.youdao.com/api"
	salt     = "1x2r6y68"
	from     = "auto"
        signType = "v3"
)

// 字段名起始字母必须是大写，否则外部访问不了，即会解析不到值
type ydTransRst struct {
	Basic     ydBasic
	Status    int
	ErrorCode string
	Translation []string
}

type ydBasic struct {
	Phonetic   string
	UkPhonetic string `json:"uk-phonetic"`
	UsPhonetic string `json:"us-phonetic"`
	Explains   []string
	Wfs        []ydWfOut
}

type ydWfOut struct {
	Wf ydWf
}

type ydWf struct {
	Name  string
	Value string
}

func main() {
	// Check translate word
	if len(os.Args) != 2 {
		printErr("Params Error", errors.New("Please input the valid cound of translate word"))
		return
	}
	TransByYd(os.Args[1])
}

//export TransByYd
func TransByYd(q string) string {
        if strings.Compare("", q) == 0 {
		errMsg := "Please input a word for translate"
		printErr("Params Error", errors.New(errMsg))
		return errMsg
	}

	// Generate request params
	toStr := getTo(q)
	curTime := strconv.FormatInt(time.Now().Unix(), 10)
	sign := generateSign(q, curTime)

	// Create http client
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		panic(err)
	}

	// Add query params
	query := req.URL.Query()
	query.Add("appKey", appID)
	query.Add("signType", signType)
	query.Add("from", from)
	query.Add("to", toStr)
	query.Add("salt", salt)
	query.Add("curtime", curTime)
	query.Add("sign", sign)
	query.Add("q", q)
	req.URL.RawQuery = query.Encode()
	// fmt.Printf("reqUrl: %s\n\n", req.URL);

	// Request
	resp, err := client.Do(req)
	if err != nil {
		printErr("Request Error", err)
		return err.Error()
	}

	// Read response
	// defer 总是在方法最后执行，即使它放在代码中间位置
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		printErr("Response Read Error", err)
		return err.Error()
	}

	// Parse
	bodyStr := string(body)
	// fmt.Printf("rspBody: %s\n", bodyStr)
	var transRst ydTransRst
	unmarshalErr := json.Unmarshal(body, &transRst)
	if unmarshalErr != nil {
		printErr("Response Unmarshal Error", unmarshalErr)
		return unmarshalErr.Error()
	}
	status := transRst.Status
	if status != 0 {
		printErr("Response Status Error", errors.New(bodyStr))
		return bodyStr
	}
	errorCode := transRst.ErrorCode
	if strings.Compare("", errorCode) != 0 && strings.Compare("0", errorCode) != 0 {
		printErr("Response Code Error", errors.New(bodyStr))
		return bodyStr
	}

	// Print query word
	fmt.Printf("\n %s \n\n", q)

	// Phonetic
	// printPhonetic(q, transRst.Basic)

	// Ws
	// wfs := transRst.Basic.Wfs
	// if len(wfs) > 0 {
	// 	for _, wfOut := range wfs {
	// 		name := wfOut.Wf.Name
	// 		value := wfOut.Wf.Value
	// 		nameLen := utf8.RuneCountInString(name)
	// 		w := 0
	// 		if nameLen == 2 {
	// 			w = 10
	// 		} else if nameLen == 3 {
	// 			w = 9
	// 		} else if nameLen == 4 {
	// 			w = 8
	// 		} else if nameLen == 5 {
	// 			w = 7
	// 		} else if nameLen == 6 {
	// 			w = 6
	// 		}
	// 		fmt.Printf(" %-"+strconv.Itoa(w)+"s : %s\n", name, value)
	// 	}
	// 	fmt.Println()
	// }

	// Explains
	for _, explain := range transRst.Translation {
		fmt.Println(" - " + explain)
	}
	fmt.Println()
	return bodyStr
}

func printPhonetic(q string, basic ydBasic) {
	phonetic := basic.Phonetic
	ukPhonetic := basic.UkPhonetic
	usPhonetic := basic.UsPhonetic
	if strings.Compare("", phonetic) != 0 {
		fmt.Printf("  音 [ %s ]", phonetic)
	}
	if strings.Compare("", ukPhonetic) != 0 {
		fmt.Printf("  英 [ %s ]", ukPhonetic)
	}
	if strings.Compare("", usPhonetic) != 0 {
		fmt.Printf("  美 [ %s ]", usPhonetic)
	}
	fmt.Printf("\n\n")
}

func printErr(errType string, err error) {
	fmt.Printf("\n%s: \n\n    %v\n\n", errType, err)
}

func getTo(q string) string {
	if regexp.MustCompile(`^[a-z A-Z]+$`).MatchString(q) {
		return "zh-CHS"
	}
	return "en"
}

func generateSign(q string, curTime string) string {
	input := ""
	qLen := utf8.RuneCountInString(q)
	if qLen < 20 {
		input = q
	} else {
		input = q[:10] + strconv.Itoa(qLen) + q[20:qLen]
	}
	oriSign := appID + input + salt + curTime + appKey
	h := sha256.New()
	h.Write([]byte(oriSign))
	x := hex.EncodeToString(h.Sum(nil))
	return x
}
