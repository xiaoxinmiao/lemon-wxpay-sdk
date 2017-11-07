package wxpay

import (
	"fmt"
	"strings"

	"github.com/relax-space/go-kit/random"

	"github.com/relax-space/go-kit/httpreq"

	"github.com/relax-space/go-kit/base"
	"github.com/relax-space/go-kit/sign"
)

func Pay(reqDto reqPayDto, custDto reqCustomerDto) (result map[string]interface{}, err error) {
	wxPayData := BuildCommonparam(reqDto.reqBaseDto)

	SetValue(wxPayData, "body", reqDto.Body)
	SetValue(wxPayData, "out_trade_no", random.Uuid(OUTTRADENO))
	SetValue(wxPayData, "total_fee", reqDto.TotalFee)
	SetValue(wxPayData, "auth_code", reqDto.AuthCode)
	SetValue(wxPayData, "device_info", reqDto.DeviceInfo)

	SetValue(wxPayData, "detail", reqDto.Detail)
	SetValue(wxPayData, "attach", reqDto.Attach)
	SetValue(wxPayData, "fee_type", reqDto.FeeType)
	SetValue(wxPayData, "goods_tag", reqDto.GoodsTag)
	SetValue(wxPayData, "limit_pay", reqDto.LimitPay)
	if len(strings.TrimSpace(reqDto.SpbillCreateIp)) == 0 {
		SetValue(wxPayData, "spbill_create_ip", "8.8.8.8")
	} else {
		SetValue(wxPayData, "spbill_create_ip", reqDto.SpbillCreateIp)
	}
	signStr := base.JoinMapObject(wxPayData.DataAttr)
	SetValue(wxPayData, "sign", sign.MakeMd5Sign(signStr, custDto.Key))
	_, body, err := httpreq.NewPost(URLPAY, []byte(wxPayData.ToXml()),
		&httpreq.Header{ContentType: httpreq.MIMEApplicationXMLCharsetUTF8}, nil)
	if err != nil {
		err = fmt.Errorf("%v:%v", MESSAGE_WECHAT, err)
		return
	}
	result, err = RespParse(body, custDto.Key)
	if err != nil {
		return
	}
	return
}

// func LoopQuery(rawMap map[string]string, limit, interval int) (queryResult map[string]interface{}, err error) {
// 	count := limit / interval
// 	waitTime := time.Duration(interval) * time.Second //2s
// 	for index := 0; index < count; index++ {
// 		queryResult, err = Query(rawMap)
// 		if err == nil { // 1. query success
// 			tradeStatusObj, ok := queryResult["trade_state"]
// 			if !ok {
// 				continue
// 			}
// 			tradeStatus := tradeStatusObj.(string)

// 			switch {
// 			case tradeStatus == "TRADE_SUCCESS":
// 				return
// 			case tradeStatus == "TRADE_CLOSED" || tradeStatus == "TRADE_FINISHED":
// 				err = errors.New("wechat pay failure")
// 				return
// 			case tradeStatus == "WAIT_BUYER_PAY":
// 				time.Sleep(waitTime)
// 				continue
// 			default:
// 				err = errors.New("wechat pay failure")
// 				return
// 			}
// 		}
// 	}
// 	err = fmt.Errorf("wechat pay overtime:%v(s)", limit)
// 	return
// }
