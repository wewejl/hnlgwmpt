package routers

import (
	"github.com/astaxie/beego"
	"ttsx/controllers"
)

func init()  {
	//去结算用户的账单
	beego.Router("/CloseAccount",&controllers.PaymentController{},"get:ShowCloseAccount")
	//提交订单
	beego.Router("/insertOrder",&controllers.PaymentController{},"post:InsertOrder")
	//到支付页面
	beego.Router("/user/userCenterOrder",&controllers.PaymentController{},"get:ShowUserCenterOrder")
}