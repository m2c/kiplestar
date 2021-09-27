package utils

import (
	"encoding/json"
	slog "github.com/m2c/kiplestar/commons/log"
	"net/http"
)

type RiskResp struct {
	FraudScore int64              `json:"fraud_score"`
	Action     string             `json:"action"`
	Results    []ScoreResultArray `json:"results"`
}

func (risk *RiskResp) IsBlocked() bool {
	return risk.Action == "blocked"
}

type ScoreResultArray struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

//for invoke kiple risk
type RiskPath string

const (
	//Multiple Identification ID
	RiskRegistration RiskPath = "?type=registration"
	//Surge in Activities in Dormant Account
	// Time Zone Mismatches
	// Sudden Spike in Transaction Value
	//Excessive Number of Transaction or Value Transacted
	RiskTransaction RiskPath = "?type=transaction"
	//Void Transactions (Exceed)
	RiskVoid RiskPath = "?type=void"
	//Exceesive Number of Topup
	RiskTopUp RiskPath = "?type=topup"
	//Widthdrawal Abuse
	RiskWidthDraw RiskPath = "?type=widthdrawal"
	//Profile Change
	RiskUpdateProfile RiskPath = "?type=updateProfile"
	//Login Failure
	//Incorrect Credentials
	RiskLoginFailed RiskPath = "?type=loginFailed"
	//Unusual Login Interval
	RiskLogin RiskPath = "?type=login"
	//Peer to Peer Transfer
	RiskTransfer RiskPath = "?type=transfer"
)

type RiskControl struct {
	host    string
	xApiKey string
	mock    bool
}

//
func RiskInstance(host, xApiKey string, mock bool) *RiskControl {
	r := new(RiskControl)
	r.host = host
	r.xApiKey = xApiKey
	r.mock = mock
	return r
}

func (r *RiskControl) Exec(url RiskPath, req interface{}) (*RiskResp, error) {
	if r.mock {
		slog.Info("======= mock Risk Control ======")
		return new(RiskResp), nil
	}
	bts, err := RequestBaseForm(r.host+string(url), req, http.Header{"x-api-key": []string{r.xApiKey}})
	if err != nil {
		return nil, err
	}
	resp := new(RiskResp)
	err = json.Unmarshal(bts, resp)
	if err != nil {
		slog.Errorf("error to Unmarshal:%s", err.Error())
		return nil, err
	}
	return resp, nil
}

type RiskRegistrationReq struct {
	FullName         string `json:"fullname"`
	Email            string `json:"email"`
	IdentificationId string `json:"identification_id"`
	MobileNumber     string `json:"mobile_number"`
	MemberId         string `json:"member_id"`
}

//Void Transactions
//Transaction
//Topup
//Widthdrawal Abuse
type RiskTransactionReq struct {
	TransactionId     string `json:"transaction_id"`
	GatewayEntryId    string `json:"gateway_entry_id"`
	OrdDate           string `json:"ord_date"` //2020-07-08 16:30:00
	OrdShipName       string `json:"ord_shipname"`
	OrdMercID         string `json:"ord_merc_id"`
	OrdMercName       string `json:"ord_merc_name"`
	OrdMercref        string `json:"ord_mercref"`
	OrdTotalamt       string `json:"ord_totalamt"`
	OrdEmail          string `json:"ord_email"`
	OrdTelephone      string `json:"ord_telephone"`
	OpAddress         string `json:"ip_address"`
	MasterMerchant_id string `json:"master_merchant_id"`
	ForeignAmount     string `json:"foreign_amount"`
	Currency          string `json:"currency"`
	ServiceCharges    string `json:"service_charges"`
	DeliveryCharges   string `json:"delivery_charges"`
	PaymentMethod     string `json:"payment_method"`
	CardBin           string `json:"card_bin"`
	TransactionType   string `json:"transaction_type"`
	AccountNo         string `json:"account_no"`
}

type RiskUpdateProfileReq struct {
	FullName  string `json:"fullname"`
	AccountNo string `json:"account_no"`
}

//Incorrect Credentials
//Login Failure
//Unusual Login Interval
type RiskLoginReq struct {
	AccountNo string `json:"account_no"`
	IpAddress string `json:"ip_address"`
}

type RiskTransferReq struct {
	AccountNo string `json:"account_no"`
}
