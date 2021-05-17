package utils

import (
	"fmt"
	slog "github.com/m2c/kiplestar/commons/log"
	"net/http"
)

type SendEmailForm struct {
	TemplateId         uint64   `json:"templateId"`
	TemplateName       string   `json:"templateName"`
	ApiKey             string   `json:"appKey" validate:"required"`
	Secret             string   `json:"secret" validate:"required"`
	EmailCustomerTitle string   `json:"email_customer_title" validate:"max=150"`
	MailTo             string   `json:"mailTo"`
	Mails              []string `json:"mails"`
	ReplaceWords       []string `json:"replaceWords"`
	AttachFileName     string   `json:"attachFileName"`
	AttachFile         []byte   `json:"attachFile"`
}

//send  email with file
func SendEmailWithFile(url, appKey, secret string, data interface{}, fileName string, address, template,title string,replaceWords []string) error {
	file, err := DataToExcelByte(data)
	if err != nil {
		return err
	}
	content :=SendEmailForm{
		TemplateName:       template,
		ApiKey:             appKey,
		Secret:             secret,
		EmailCustomerTitle: title,
		MailTo:             address,
		ReplaceWords:       replaceWords,
		AttachFileName:     fileName,
		AttachFile:         file,
	}

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

//send email with no file
func SendEmail(url, appKey, secret , address, template,title string,replaceWords []string) error {
	content :=SendEmailForm{
		TemplateName:       template,
		ApiKey:             appKey,
		Secret:             secret,
		EmailCustomerTitle: title,
		MailTo:             address,
		ReplaceWords:       replaceWords,
	}

	code, err := Request(http.MethodPost, url, content, nil, nil)

	if err != nil {
		slog.Errorf("send email failed content[%v ]err[%v]",content, err)
		return err
	}
	if code != 0 {
		err := fmt.Errorf("emil send error content[%v]",content)
		return err
	}
	return nil
}