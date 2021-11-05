package utils

// response start

type RiskResp struct {
	FraudScore int64              `json:"fraud_score"`
	Action     string             `json:"action"`
	Results    []ScoreResultArray `json:"results"`
}

type PortalResp struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (risk *RiskResp) IsBlocked() bool {
	return risk.Action == "blocked"
}

type ScoreResultArray struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

//response end

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
	TransactionId    string `json:"transaction_id"`
	GatewayEntryId   string `json:"gateway_entry_id"`
	OrdDate          string `json:"ord_date"` //2020-07-08 16:30:00
	OrdShipName      string `json:"ord_shipname"`
	OrdMercID        string `json:"ord_mercID"`
	OrdMercName      string `json:"ord_mercName"`
	MID              string `json:"merchant_id"`
	MercName         string `json:"merchant_name"`
	OrdMercref       string `json:"ord_mercref"`
	OrdTotalamt      string `json:"ord_totalamt"`
	Amount           string `json:"amount"`
	OrdEmail         string `json:"ord_email"`
	OrdTelephone     string `json:"ord_telephone"`
	IpAddress        string `json:"ip_address"`
	MasterMerchantId string `json:"master_merchant_id"`
	ForeignAmount    string `json:"foreign_amount"`
	Currency         string `json:"currency"`
	ServiceCharges   string `json:"service_charges"`
	DeliveryCharges  string `json:"delivery_charges"`
	PaymentMethod    string `json:"payment_method"`
	CardBin          string `json:"card_bin"`
	TransactionType  string `json:"transaction_type"`
	AccountNo        string `json:"account_no"`
}

type RiskAfterPaymentReq struct {
	OrdMercref      string `json:"ord_mercref"`
	OrdDate         string `json:"ord_date"` //2020-07-08 16:30:00
	OrdTotalamt     string `json:"ord_totalamt"`
	Amount          string `json:"amount"`
	OrdEmail        string `json:"ord_email"`
	OrdShipName     string `json:"ord_shipname"`
	OrdDelcharges   string `json:"ord_delcharges"`
	OrdShipcountry  string `json:"ord_shipcountry"`
	OrdTelephone    string `json:"ord_telephone"`
	OrdReturnURL    string `json:"ord_return_url"`
	OrdSvccharges   string `json:"ord_svccharges"`
	OrdGstamt       string `json:"ord_gstamt"`
	IpAddress       string `json:"ip_address"`
	CardNumber      string `json:"card_number"`
	PaymentMethod   string `json:"payment_method"`
	TransactionId   string `json:"transaction_id"`
	GatewayEntryId  string `json:"gateway_entry_id"`
	TransactionType string `json:"transaction_type"`
	OrdMercID       string `json:"ord_mercID"`
	OrdMercName     string `json:"ord_mercName"`
	MID             string `json:"merchant_id"`
	MercName        string `json:"merchant_name"`
	AccountNo       string `json:"account_no"`
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

//订单号,金额, email, merchant Id, merchant Name
type RiskTransferReq struct {
	AccountNo     string `json:"account_no"`
	TransactionId string `json:"transaction_id"`
	OrdTotalamt   string `json:"amount"`
	OrdEmail      string `json:"ord_email"`
	OrdMercID     string `json:"ord_mercID"`
	OrdMercName   string `json:"ord_mercName"`
	MID           string `json:"merchant_id"`
	MercName      string `json:"merchant_name"`
}

type PortalUserInfo struct {
	AccountNo        string `json:"account_no"`
	SourceMerchantId string `json:"source_merchant_id"`
	Name             string `json:"name"`
	CreatedAt        string `json:"created_at"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	IdentificationId string `json:"identification_id"`
	UserRace         string `json:"user_race"`
	Address          string `json:"address"`
	City             string `json:"city"`
	State            string `json:"state"`
	Postcode         string `json:"postcode"`
	MobileNumber     string `json:"mobile_number"`
}

type PortalUserStatus struct {
	AccountNo string `json:"account_no"`
	Status    string `json:"status"`
}

type RiskBackListReq struct {
	OrdMercref        string `json:"ord_mercref"`
	OrdDate           string `json:"ord_date"`
	OrdTotalamt       string `json:"ord_totalamt"`
	Amount            string `json:"amount"`
	OrdEmail          string `json:"ord_email"`
	OrdShipname       string `json:"ord_shipname"`
	OrdMercID         string `json:"ord_mercID"`
	OrdMercName       string `json:"ord_mercName"`
	MID               string `json:"merchant_id"`
	MercName          string `json:"merchant_name"`
	OrdDelcharges     string `json:"ord_delcharges"`
	OrdShipcountry    string `json:"ord_shipcountry"`
	OrdTelephone      string `json:"ord_telephone"`
	OrdReturnURL      string `json:"ord_return_url"`
	OrdSvccharges     string `json:"ord_svccharges"`
	OrdGstamt         string `json:"ord_gstamt"`
	MerchantHashvalue string `json:"merchant_hashvalue"`
	IpAddress         string `json:"ip_address"`
	GatewayEntryId    string `json:"gateway_entry_id"`
	PaymentMethod     string `json:"payment_method"`
	IdentificationNo  string `json:"identification_no"`
}
