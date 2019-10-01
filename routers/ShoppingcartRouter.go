package routers

import (
	"github.com/astaxie/beego"
	"ttsx/controllers"
)

func init()  {
	//添加addCart
	beego.Router("/addCart",&controllers.ShoppingcartController{},"post:AddShoppingcart")
	//显示我的购物车
	beego.Router("/user/MyShoppingCart",&controllers.ShoppingcartController{},"get:ShowMyShoppingCart")
	//删除一个用户的行购物车
	beego.Router("/deleteCart",&controllers.ShoppingcartController{},"post:DeleteShoppingCart")

}