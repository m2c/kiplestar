package utils

import (
	"encoding/base64"
	"fmt"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/m2c/kiplestar/config"
	"net/http"
)

//send  email with file
func SendEmail(appkey, secret string, data interface{}, fileName string, address, template string) error {
	file, err := DataToExcelByte(data)
	if err != nil {
		return err
	}
	bufStore := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(bufStore, file)

	content := make(map[string]interface{})
	content["channelId"] = 1
	content["templateName"] = template
	content["mailTo"] = address
	content["attachFile"] = bufStore
	content["attachFileName"] = fileName
	content["apiKey"] = appkey
	content["secret"] = secret
	url := config.Configs.Notify.Url + config.EmailSendUrl

	code, err := Request(http.MethodPost, url, content, nil, nil)

	if err != nil {
		slog.Errorf("send email err[%v]", err)
		return err
	}
	if code != 0 {
		err := fmt.Errorf("emil send error")
		return err
	}
	return nil
}
