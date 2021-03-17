package utils

import (
	"fmt"
	slog "github.com/m2c/kiplestar/commons/log"
	"net/http"
)

//send  email with file
func SendEmail(url, appKey, secret string, data interface{}, fileName string, address, template string) error {
	file, err := DataToExcelByte(data)
	if err != nil {
		return err
	}

	content := make(map[string]interface{})
	content["templateName"] = template
	content["mailTo"] = address
	content["attachFile"] = file
	content["attachFileName"] = fileName
	content["appKey"] = appKey
	content["secret"] = secret

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
