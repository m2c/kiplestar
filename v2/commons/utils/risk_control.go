package utils

import (
	"encoding/json"
	slog "github.com/m2c/kiplestar/commons/log"
	"net/http"
)

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
	RiskTransfer            RiskPath = "?type=transfer"
	RiskCompleteTransaction RiskPath = "?type=completeTransaction"
	RiskBlackList           RiskPath = "?type=blacklist"

	RiskUserCreation RiskPath = "/api/user-creation"
	RiskUserUpdate   RiskPath = "/api/user-update"
	RiskUserStatus   RiskPath = "/api/user-status"
)

type RiskControl struct {
	host       string
	xApiKey    string
	mock       bool
	portalHost string
}

//
func RiskInstance(host, portalHost, xApiKey string, mock bool) *RiskControl {
	r := new(RiskControl)
	r.host = host
	r.xApiKey = xApiKey
	r.mock = mock
	r.portalHost = portalHost
	return r
}

func (r *RiskControl) ExecAsync(url RiskPath, req interface{}) {
	go func() {
		if r := recover(); r != nil {
			slog.Errorf("error to invoke ExecAsync , %v", r)
		}
		r.Exec(url, req)
	}()
}

func (r *RiskControl) Exec(url RiskPath, req interface{}) (*RiskResp, error) {
	if r.mock {
		slog.Info("======= mock Risk Control ======")
		return new(RiskResp), nil
	}
	bts, err := RequestBaseForm(r.host+string(url), req, http.Header{"x-api-key": []string{r.xApiKey}})
	if err != nil {
		//network error ,will Through risk control
		return new(RiskResp), nil
	}
	resp := new(RiskResp)
	err = json.Unmarshal(bts, resp)
	if err != nil {
		//parse error ,will Through risk control
		slog.Errorf("error to Unmarshal:%s", err.Error())
		return resp, nil
	}
	return resp, nil
}

func (r *RiskControl) PortalExecAsync(url RiskPath, req interface{}) {
	go func() {
		if r := recover(); r != nil {
			slog.Errorf("error to invoke PortalExecAsync , %v", r)
		}
		r.PortalExec(url, req)
	}()
}

func (r *RiskControl) PortalExec(url RiskPath, req interface{}) (*PortalResp, error) {
	if r.mock {
		slog.Info("======= mock Risk Control ======")
		return &PortalResp{Success: true}, nil
	}
	bts, err := RequestBaseForm(r.portalHost+string(url), req, http.Header{"x-api-key": []string{r.xApiKey}})
	if err != nil {
		//network error ,will Through risk control
		return &PortalResp{Success: true}, nil
	}
	resp := new(PortalResp)
	err = json.Unmarshal(bts, resp)
	if err != nil {
		//parse error ,will Through risk control
		slog.Errorf("error to Unmarshal:%s", err.Error())
		return &PortalResp{Success: true}, nil
	}
	return resp, nil
}
