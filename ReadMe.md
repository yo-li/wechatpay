# 微信支付
微信最常用两种支付模式: **公众号支付**、**扫码支付**   

获取包
```
go get -u github.com/yo-li/wechatpay
```
#### 公众号支付调用示例

```go
import(
	"github.com/yo-li/wechatpay"
	"fmt"
)

func main() {
	var pram map[string]string
	pram = make(map[string]string)
	//公众账号ID
	pram["appid"] = "微信appid"
	//商户号
	pram["mch_id"] = "商户号"
	//设备号
	pram["device_info"] = "WEB"
	//随机字符串
	pram["nonce_str"] = wechatpay.Get_Nonce_Str()
	//签名
	//pram["sign"] = "wx075f07aad3dbbac6"
	//签名类型
	pram["sign_type"] = "MD5"
	//商品描述
	pram["body"] = "golang 语言测试"
	//商品详情
	pram["detail"] = ""
	//附加数据。
	pram["attach"] = ""
	//商户订单号
	pram["out_trade_no"] = "test201710111956"
	//标价币种
	pram["fee_type"] = "CNY"
	//标价金额
	pram["total_fee"] = "2"
	//终端IP
	pram["spbill_create_ip"] = "113.87.129.111"
	//交易起始时间
	pram["time_start"] = ""
	//交易结束时间
	pram["time_expire"] = ""
	//订单优惠标记
	pram["goods_tag"] = ""
	//通知地址
	pram["notify_url"] = "回调地址"
	//交易类型
	pram["trade_type"] = "JSAPI"
	//商品ID
	pram["product_id"] = ""
	//指定支付方式
	pram["limit_pay"] = ""
	//用户标识
	pram["openid"] = "oYjOVjg0EpQS7pp5P90AZTtNI23M"
	//场景信息
	pram["scene_info"] = ""
	//门店id
	pram["id"] = ""
	//门店名称
	pram["name"] = ""
	//门店行政区划码
	pram["area_code"] = ""
	//门店详细地址
	pram["address"] = ""

	//返回支付参数对象
	//
	p := wechatpay.PayH5(pram, "支付密钥")

	//返回结果ReturnResult=true,表示执行成功;ReturnResult=false,表示失败，具体查看p.ReturnMsg具体错误
	//如果成功时返回对象属性是对应js调起支付参数:appId,timeStamp,nonceStr,package,signType,paySign
	if p.ReturnResult {
		//参考微信支付官网:https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=7_7&index=6
		//p.AppId
		//p.TimeStamp
		//p.NonceStr
		//p.Package
		//p.SignType
		//p.PaySign
		fmt.Printf("%+v\n", p)
	} else {
		fmt.Println(p.ReturnMsg)
	}

}
```


**返回成功示例,公众号页面js调起支付所需参数**  
```
{
  ReturnResult:true
  ReturnMsg:SUCCESS
  AppId:wx075f07aad3dbbac6
  TimeStamp:1511532594
  NonceStr:NjgwOTQ3NDgwMjk5NDU5
  Package:prepay_id=wx2017112422095418ba7728590601419874
  SignType:MD5
  PaySign:A61FDFA1376B057021AD6EB7D0FC126F
}
```
**返回失败示例**  
```
appid参数长度有误
```

#### 扫码支付调用示例
```go
import (
	"fmt"
	"wechatpay"
)

func main() {
	var pram map[string]string
	pram = make(map[string]string)
	//公众账号ID
	pram["appid"] = "微信appid"
	//商户号
	pram["mch_id"] = "商户号"
	//设备号
	pram["device_info"] = "WEB"
	//随机字符串
	pram["nonce_str"] = wechatpay.Get_Nonce_Str()
	//签名
	//pram["sign"] = "wx075f07aad3dbbac6"
	//签名类型
	pram["sign_type"] = "MD5"
	//商品描述
	pram["body"] = "golang 语言测试"
	//商品详情
	pram["detail"] = ""
	//附加数据。
	pram["attach"] = ""
	//商户订单号
	pram["out_trade_no"] = "test201710111957"
	//标价币种
	pram["fee_type"] = "CNY"
	//标价金额
	pram["total_fee"] = "2"
	//终端IP
	pram["spbill_create_ip"] = "113.87.129.111"
	//交易起始时间
	pram["time_start"] = ""
	//交易结束时间
	pram["time_expire"] = ""
	//订单优惠标记
	pram["goods_tag"] = ""
	//通知地址
	pram["notify_url"] = "http://www.baidu.com"
	//交易类型
	pram["trade_type"] = "NATIVE"
	//商品ID
	pram["product_id"] = ""
	//指定支付方式
	pram["limit_pay"] = ""
	//用户标识
	pram["openid"] = ""
	//场景信息
	pram["scene_info"] = ""
	//门店id
	pram["id"] = ""
	//门店名称
	pram["name"] = ""
	//门店行政区划码
	pram["area_code"] = ""
	//门店详细地址
	pram["address"] = ""

	//返回支付参数对象
	//
	p, qrurl := wechatpay.PayNATIVE(pram, "支付密钥")
  //p=true,qrurl返回支付二维码链接
	//p=false,qrurl返回错误内容，请根据提醒修改正大确
	if p {
		fmt.Println(qrurl)
	} else {
		fmt.Println(qrurl)
	}

}
```
**返回成功示例，此处不会生成二维码**  
```
weixin://wxpay/bizpayurl?pr=RYTtuge
```
**返回失败示例**  
```
201 商户订单号重复
```
