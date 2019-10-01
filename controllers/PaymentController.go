package controllers

import (
	"github.com/astaxie/beego"
	"ttsx/models"
	"fmt"
	"github.com/smartwalle/alipay"
)

type PaymentController struct {
	beego.Controller
}

//展现提交订单页面
func (this *PaymentController)ShowCloseAccount() {
	//获取数据
	skuids := this.GetStrings("skuid")
	if len(skuids) == 0 {
		this.Redirect("/user/MyShoppingCart", 302)
		return
	}
	//获取用户名
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("login", 302)
		return
	}
	//在数据调取用户信息的地址
	address, err := models.SeclectUserinfo(userName.(string))
	if err != nil {
		//fmt.Println(err)
		this.Redirect("/user/Personalcentersite", 302)
		return
	}
	//判断用户有没有注册地址
	//if address.Addr ==""{
	//	this.Redirect("/user/Personalcentersite",302)
	//	return
	//}
	//在数据库调取
	goodsmap, llroparic, llrocontent, err := models.PaymentCloseAccount(skuids, userName.(string))
	if err != nil {
		fmt.Println(err)
		this.Ctx.WriteString("数据库错误")
		return
	}
	this.Data["llroparic"] = llroparic
	this.Data["llrocontent"] = llrocontent
	this.Data["llroparicyf"] = llroparic + 1
	this.Data["skuids"] = skuids
	this.Data["address"] = address
	this.Data["goodsmap"] = goodsmap
	this.TplName = "place_order.html"
}

//递交一个完整的订单
func (this *PaymentController) InsertOrder() {
	resp := make(map[string]interface{})
	defer Jsonserver(this, resp)
	userName := this.GetSession("userName")
	if userName == nil {
		resp["status"] = 401
		resp["msg"] = "用户没有登录"
		return
	}
	//地址Id
	addrId, err := this.GetInt("addrId")
	//支付方式id
	payId, err1 := this.GetInt("payId")
	//skuids 所有商品id
	skuids := this.GetString("skuids")
	//获取总件数，总价格，快递费
	totalCount, err2 := this.GetInt("totalCount")
	totalPrice, err3 := this.GetInt("totalPrice")
	transit, err4 := this.GetInt("transit")
	//数据校验
	if err != nil || err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		resp["status"] = 402
		resp["msg"] = "信息获取失败"
		return
	}
	//把数据插入订单表里面
	err,msg:=models.InsertOrderInfoList(userName.(string),addrId, payId,skuids, totalCount, totalPrice, transit)
	if err!=nil {
		resp["status"] = 403
		resp["msg"] =msg
		return
	}

	resp["status"]=200
	resp["msg"]="成功"
}

func Jsonserver(this *PaymentController, resp map[string]interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()

}

//跳转到支付页面
func (this *PaymentController)ShowUserCenterOrder()  {
	var privateKey =`MIIEpAIBAAKCAQEAvaLI0mc/eDS0ysmR1+pOTqLNVgWafDBKDWOWXYjU5prslw8X
KTs+zvC8SXMaqER4VaEdaFpTrk3QefiClc6MY7L08TxyG0Y+PXFaSJup43cgM8FB
BYf+lmhoqPK61o1aTv0JrFD0KPiEV/QZDbkIsmTKrdKLv2EQq1XumS+vxI95qi4B
vmzAKWRri7n8oFU13gH+Ics2y7NVD6TDtRVEV6OfOvDgpLnmLXq2tj38lnSLBRgP
4GS/w+Pfh64D9lcTxn69RudD+Sxz6vYoW7uaziSjq15Ni3LURev/IZ34Kzk12GjP
PDW2pAQIpDXzaPkVQYXuQ4Oz5ZaKsGF0Yv5ipQIDAQABAoIBAAPGrwsJhUkGe6ci
FmZfQwnr0fzphab9aywTFJZuOBcTdKyZX1Ox21FRl946jYhWPLMvzx8Z1Vq+L+2N
1kPXZhJCKQB4vKjwYCLnE+4oM1zVLW36ZioPCDHEiHj8xF2rWOYDweKNhh8eu7vo
n2sXiSXMPgDyFVhNPYC76FFikrRuQiOfQDQtZiUr6fyeYpKjEL0xI6cJym1hYqmD
gtvctzggDCiQj7JMBUH44+GmiNCnGLzNG6RmbS0TA8gEuQu8z73qzTItEMTJjCaJ
nQPcdKa3YlhWQZoTSfl910l1W4mhyuubB540ipVHRJTulrgb7I1x2UyJ5Ff+nc+4
6SQmuo0CgYEA8zJ0FPfC822hUA+4dA7ZMMCiF61TZdxj+daAKEUMB7hU1a9GN1FY
c6bf2dfUIH6vlKKvpK/8SA7mPd3OXMMo9XupLodcvfHd/yT4a+jmx4Qv96j1wfNe
ioJh1CP7VLS0OMSwsZMbNo6UMZhIs+6Au9YKyrYBduvePo7xFuhMLgMCgYEAx559
vD0H5bspee7jSTOXNrHl6fN/klLN4HEEXMlgUdb3of7fYSnZMpOF5sRwC7mgM84i
fucHDhHv7ZME6ZxawC9DlgLGj92oHgzata+NH5CRynhiR6eZFaBD+6fIf/VgMltU
nVH5MXSbnPEEqRONNaUTan1irZk2eMplA581gDcCgYEA6DDUiaxfsiCKcjEAL7Z5
gMV6PNbcGBWKUl+MbmY17Sz9uiKlDG2a4JiDgq5AtmGd63BD+B2Z5YZsJscdno0q
Du5pAaZ1UliZVl+K2yQ7KmQ3k+H5+ZoNOnrvQia0cBQzOTv5YyELS1Rngs5dI4Vj
3XKnTRDmZw8dWmcJIZDaItcCgYEAnUmFuye/rEWAFeKkVk5/TIp6JZBGqc3zCHEk
xdOqwHGIp61C57VovZA+BqoruyFlWMyIo8N37J83lNOuIEChxSK4t1+ygzNdP2hT
gKs1oHRyW73lep5VYhPo3UbEFgcK6ELMdjVcC5rc7pl+WZbdQjKzDMqFUVIS+LRJ
ScRODJsCgYB5PLWqzyTelIgE/ocSnpF94CzV4WQTi/9k4xjikFNUJPJVOKztHM3C
Dr7ayz24CSyIHHXH777sU86RAIAY75nJgoSmlpFR7nKdDo7I0q3ETCSi3sTS7L0f
r3Ql55OeUcIs5CmAZzvMFrGuWLPzNDutXa6XwLxbPmviNwNmv6ph5g==`
	var appId = "2016101200668306"
	var aliPublicKey =`MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAq5pdKVjpZzvKFaXFMOUYg/7YggJiZ9NZQlShikRiI2Q847JG0f1pZqv4aCqn7CbOTYeM9QXkwLJ4Q1GxU9tc05j+mpDu3zZyaFccSyD5rQHP7S3pdKqu+NnbZ0M2mHS3aepwjyx3hUBW0XeZH08MGaH2yiW+UmsEpJSVgCTqCUEgaV7+tGGDT+f7b501cYKdWQgB02GagzfwsznGHTP/nUWMLC6/Gp9orv+67j7aMLZqwdqpopGW4Os44rPClJyMRrDH/zncNWDwxFazXqASKnDYC5J8lzi3ltL0QN4LC611FPxiw6/5G48hQfLgxf50cuLwXEngLhB9iKEYpo0aHQIDAQAB`
	//func New(appId, aliPublicKey, privateKey string, isProduction bool) (client *Client, err error) {
	client ,err:= alipay.New(appId,aliPublicKey,privateKey,false)
	if err!=nil {
		this.Data["errer"]="支付失败"
		this.TplName="place_order.html"
		return
	}
	//支付配置
	p :=alipay.TradePagePay{}
	//设置支付主题
	p.Subject="学校生鲜平台"
	//设置订单号
	p.OutTradeNo="201909270150266"
	//设置支付金额
	p.TotalAmount="1000.00"
	//设置异步返回方式
	p.NotifyURL="http://192.168.73.128:8080/user/userCenterInfo"
	//设置同步返回方式
	p.ReturnURL="http://192.168.73.128:8080/user/userCenterInfo"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	url,err:=client.TradePagePay(p)
	if err!=nil {
		this.Data["errer"]="支付失败2"
		this.TplName="place_order.html"
		return
	}

	this.Redirect(url.String(),302)
	}
