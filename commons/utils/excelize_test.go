package utils

import "testing"

func TestDataToExcelByte(t *testing.T) {
	type sct struct {
		UserID          uint64 `json:"user_id"`
		ClientID        uint64 `json:"client_id"`
		Amount          string `json:"amount"`
		MerchantOrderNo string `json:"merchant_order_no"`
		CardBin         string `json:"card_bin"`
		ClientIP        string `json:"client_ip"`
		RefererUrl      string `json:"referer_url"`
	}
	var data []sct

	for i := 0; i < 10; i++ {
		s := sct{
			UserID:          uint64(i + 1),
			ClientID:        10,
			Amount:          "100",
			MerchantOrderNo: "423f22",
			CardBin:         "1d924234",
			ClientIP:        "127.0.0.1",
			RefererUrl:      "http://111.com",
		}
		data = append(data, s)
	}

	excelByte, err := DataToExcelByte(data)
	if err != nil {
		t.Errorf("TestDataToExcelByte err : %s", err.Error())
		return
	}

	t.Log(excelByte)
}
