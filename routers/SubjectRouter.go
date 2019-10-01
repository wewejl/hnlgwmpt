package routers

import (
	"github.com/astaxie/beego"
	"ttsx/controllers"
	"github.com/astaxie/beego/context"
)
var UserFilterFunc = func(ctr *context.Context) {
	userName:=ctr.Input.Session("userName")
	if userName==nil {
		ctr.Redirect(302,"/login")
	}
}
func init()  {
	//路由过滤器
	beego.InsertFilter("/user/*",beego.BeforeExec,UserFilterFunc)
	//index页面展示
	beego.Router("/",&controllers.SubjectCOntroller{},"get:ShowIndex")
	beego.Router("/index",&controllers.SubjectCOntroller{},"get:ShowIndex")
	//天天生鲜
	beego.Router("/indexSx",&controllers.SubjectCOntroller{},"get:ShowIndexSx")

	//商品点击到商品详情页
	beego.Router("/Goodsdetails",&controllers.SubjectCOntroller{},"get:ShowDetails")

	//商品列表单
	beego.Router("/indexSxTypelist",&controllers.SubjectCOntroller{},"get:ShowTypelist")
	//商品列表价格
	//商品查询方法
	beego.Router("/querygoodsSrcre",&controllers.SubjectCOntroller{},"get:ShowgoodsSrcre")

}