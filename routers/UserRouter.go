package routers

import (
	"github.com/astaxie/beego"
	"ttsx/controllers"
)

func init()  {
	//用户注册
	beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:InsertRegister")
	//发送信息
	beego.Router("/Sendmessage",&controllers.UserController{},"post:Sendmessage")
	//邮箱注册
	beego.Router("/registerEmail",&controllers.UserController{},"get:ShowRegisterEmail;post:InsertRegisterEmail")
	//邮箱激活
	beego.Router("/emailseed",&controllers.UserController{},"get:Emailseed")
	//用户登录
	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:SeclectLogin")
	//用户退出
	beego.Router("/layout",&controllers.UserController{},"get:LayoutSseion")
	//个人中心个人信息路由
	beego.Router("/user/Personalcenterinfo",&controllers.UserController{},"get:ShowPersonalcenterinfo")
	//个人中心全部订单路由
	beego.Router("/user/Personalcenterorder",&controllers.UserController{},"get:ShowPersonalcenterorder")
	//个人中心收货地址路由
	beego.Router("/user/Personalcentersite",&controllers.UserController{},"get:ShowPersonalcentersite;post:InsertPersonalcentersite")

}