package utils

import (
	"fmt"
	slog "github.com/m2c/kiplestar/commons/log"
	"net/http"
)

const emailApi = "/se/email/api/sendmail"

type NotifyService interface {
	SeedEmail(notifyReq *NotifyEntity) error
}

type notifyService struct {
	appKey string
	secret string
	host string

}

func NotifyServiceInstance(appKey, secret, host string) NotifyService {
	return &notifyService{
		appKey:       appKey,
		secret:     secret,
		host: host,
	}
}

//seed email
func (e *notifyService) SeedEmail(notifyReq *NotifyEntity) error {
	notifyReq.AppKey = e.appKey
	notifyReq.Secret = e.secret
	url := e.host + emailApi

	code, err := Request(http.MethodPost, url, notifyReq, nil, nil)

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

type NotifyEntity struct {
	ChannelId int `json:"channelId"`
	TemplateName string `json:"templateName"`
	MailTo string `json:"mailTo"`
	ReplaceWords []string `json:"replaceWords"`
	AttachFile []byte `json:"attachFile"`
	AttachFileName string `json:"attachFileName"`
	//system config
	AppKey string `json:"appKey"`
	Secret string `json:"secret"`
}