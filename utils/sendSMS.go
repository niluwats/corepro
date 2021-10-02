package utils

import (
	"bytes"
	"core/errs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func SendSms(to string, strCode string) *errs.AppError {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}

	urlStr := fmt.Sprintf("https://richcommunication.dialog.lk/api/sms/send/%s", config.SmsAuthKey)
	formatted_no := formatMobileNO(to)
	formatted_no = "94" + formatted_no

	msgs := make(map[string][]interface{})
	msArr := make([]interface{}, 1)
	msg1 := make(map[string]string)

	msg1["clientRef"] = ""
	msg1["mask"] = config.Mask
	msg1["campaignName"] = config.CampaignName
	msg1["number"] = formatted_no
	msg1["text"] = "Dear user, your verification code is " + strCode
	msArr[0] = msg1
	msgs["messages"] = msArr

	jStr, _ := json.Marshal(msgs)
	b := bytes.NewReader(jStr)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, b)

	t := time.Now().Format(time.RFC3339)
	ti := t[:len(t)-6]

	req.Header.Add("USER", config.SmsUser)
	req.Header.Add("DIGEST", config.SmsDigest)
	req.Header.Add("CREATED", ti)

	req.Header.Add("Content-Type", "application/json")

	resp, err0 := client.Do(req)
	if err0 != nil {
		log.Fatal(err0)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			return errs.NewUnexpectedError("failed decoding response body ! ")
		}
	} else {
		return errs.NewUnexpectedError("failed sending verification SMS ! " + resp.Status)
	}
	return nil
}

func formatMobileNO(mobile string) string {
	mob_no := mobile
	mb := mobile[1:len(mob_no)]
	fmt.Println(mb)
	return mb
}
