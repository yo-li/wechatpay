// wechatpay project wechatpay.go
package wechatpay

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

//报文结果
type Items struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	Appid      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	PrepayId   string `xml:"prepay_id"`
	TradeType  string `xml:"trade_type"`
	CodeUrl    string `xml:"code_url"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

type Result struct {
	XMLName xml.Name `xml:"xml"`
	Items
}

//h5返回对象
type PayModel struct {
	ReturnResult bool
	ReturnMsg    string
	AppId        string
	TimeStamp    string
	NonceStr     string
	Package      string
	SignType     string
	PaySign      string
}

//返回支付二维码内容
func PayNATIVE(pram map[string]string, key string) (bool, string) {

	//给参数key排序
	var pramKey []string
	pramIndex := 0
	pramKey = make([]string, len(pram))
	for key, _ := range pram {
		pramKey[pramIndex] = key
		pramIndex++
	}
	sort.Strings(pramKey)

	//拼接签名内容
	signContent := ""
	for i := 0; i < len(pramKey); i++ {
		if pram[pramKey[i]] == "" {
			continue
		}
		if signContent == "" {
			signContent = pramKey[i] + "=" + pram[pramKey[i]]
		} else {
			signContent += "&" + pramKey[i] + "=" + pram[pramKey[i]]
		}
	}
	signContent += "&key=" + key

	//生成签名
	sign := strings.ToUpper(Get_MD5(signContent))
	pram["sign"] = sign

	//生成请求报文
	xmlf := Get_MapToXML(pram, "xml")

	//fmt.Println(xmlf)

	pbody := strings.NewReader(xmlf)

	//发送https请求
	resp, err := http.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", "text/xml;charset=utf-8", pbody)
	if err != nil {
		// handle error
		return false, "http请求出错！"
	}
	//处理返回结果
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println(string(body))

	v := Result{}
	//data := string(body)
	//err = xml.Unmarshal([]byte(data), &v)
	err = xml.Unmarshal(body, &v)

	if err != nil {
		//fmt.Printf("error: %v", err)
		return false, "解析报文出错！"
	}

	if v.Items.ReturnCode != "SUCCESS" {
		return false, v.Items.ReturnMsg
	}

	if v.Items.ResultCode != "SUCCESS" {
		return false, v.Items.ErrCodeDes
	}
	//返回结果
	return true, v.Items.CodeUrl
}

//返回公众号H5支付对象
func PayH5(pram map[string]string, key string) PayModel {
	//定义返回对象
	p := PayModel{}

	//拼接签名内容
	signContent := Get_MapToString(pram)

	signContent += "&key=" + key

	//生成签名
	sign := strings.ToUpper(Get_MD5(signContent))
	pram["sign"] = sign

	//生成请求报文
	xmlf := Get_MapToXML(pram, "xml")

	//fmt.Println(xmlf)

	pbody := strings.NewReader(xmlf)

	//发送https请求
	resp, err := http.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", "text/xml;charset=utf-8", pbody)
	if err != nil {

		p.ReturnResult = false
		p.ReturnMsg = "http请求出错！"
		return p
	}
	//处理返回结果
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println(string(body))

	v := Result{}
	//data := string(body)
	//err = xml.Unmarshal([]byte(data), &v)
	err = xml.Unmarshal(body, &v)

	if err != nil {

		p.ReturnResult = false
		p.ReturnMsg = "解析报文出错！"
		return p
	}

	if v.Items.ReturnCode != "SUCCESS" {

		p.ReturnResult = false
		p.ReturnMsg = v.Items.ReturnMsg
		return p
	}

	if v.Items.ResultCode != "SUCCESS" {
		p.ReturnResult = false
		p.ReturnMsg = v.Items.ErrCodeDes
		return p
	}

	//准备h5支付参数

	p.AppId = v.Items.Appid

	t := time.Now()

	p.TimeStamp = strconv.FormatInt(t.Unix(), 10)
	p.NonceStr = Get_Nonce_Str()
	p.Package = "prepay_id=" + v.Items.PrepayId
	p.SignType = "MD5"

	var h5pram map[string]string
	h5pram = make(map[string]string)
	h5pram["appId"] = p.AppId
	h5pram["timeStamp"] = p.TimeStamp
	h5pram["nonceStr"] = p.NonceStr
	h5pram["package"] = p.Package
	h5pram["signType"] = p.SignType

	signContent = Get_MapToString(h5pram)
	signContent += "&key=" + key

	//生成签名
	sign = strings.ToUpper(Get_MD5(signContent))
	p.PaySign = sign

	p.ReturnResult = true
	p.ReturnMsg = "SUCCESS"
	return p

}

//其他方法
func Get_Nonce_Str() string {
	//获随数码，32位以内
	rand.Seed(time.Now().Unix())                   //使用时间做为种子包
	rNumber := rand.Intn(999999999999999)          //设置随时数0-999999999999999
	data := []byte(strconv.Itoa(rNumber))          //将随机数转为字串符，并转为[]byte
	str := base64.StdEncoding.EncodeToString(data) //将随机数的[]byte转为base64
	var posend int = 32
	if len(str) < 32 {
		posend = len(str)
	}

	return str[0:posend]
}

func Get_MD5(str string) string {

	h := md5.New()
	h.Write([]byte(str)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	md5str1 := fmt.Sprintf("%x", cipherStr) //将[]byte转成16进制

	return md5str1
}

func Get_MapToXML(pram map[string]string, root string) string {
	xmlContent := "<" + root + ">"
	for key, value := range pram {
		if value == "" {
			continue
		}
		xmlContent += "<" + key + ">"
		xmlContent += value
		xmlContent += "</" + key + ">"
	}
	xmlContent += "</" + root + ">"
	return xmlContent
}

func Get_MapToString(pram map[string]string) string {
	var pramKey []string
	pramIndex := 0
	pramKey = make([]string, len(pram))
	for key, _ := range pram {
		pramKey[pramIndex] = key
		pramIndex++
	}
	sort.Strings(pramKey)

	//拼接签名内容
	signContent := ""
	for i := 0; i < len(pramKey); i++ {
		if pram[pramKey[i]] == "" {
			continue
		}
		if signContent == "" {
			signContent = pramKey[i] + "=" + pram[pramKey[i]]
		} else {
			signContent += "&" + pramKey[i] + "=" + pram[pramKey[i]]
		}
	}
	return signContent
}
