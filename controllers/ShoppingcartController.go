package controllers

import (
	"github.com/astaxie/beego"
	"ttsx/models"
)

type ShoppingcartController struct {
	beego.Controller
}

func (c *ShoppingcartController)AddShoppingcart()  {
	//获取数据
	//商品id
	skuid,err1:=c.GetInt("skuid")
	//商品个数
	count,err:=c.GetInt("count")
	//建立一个容器map
	respErr:=make(map[string]interface{})
	if err!=nil ||err1!=nil{
		//赋值
		respErr["status"]=401
		respErr["msg"]="获取数据错误"
		c.Data["json"]=respErr
		//传出的方式
		c.ServeJSON()
		return
	}
	//获取用户名
	userName:=c.GetSession("userName")
	if userName==nil{
		respErr["status"]=402
		respErr["msg"]="您还没有登录"
		c.Data["json"]=respErr
		//传出的方式
		c.ServeJSON()
		return
	}
	err=models.RedisShoppingcart(userName.(string),skuid,count)
	if err!=nil {
		respErr["status"]=403
		respErr["msg"]="数据库插入错误"
		c.Data["json"]=respErr
		//传出的方式
		c.ServeJSON()
		return
	}
	resp:=make(map[string]interface{})
	resp["status"]=200
	resp["msg"]="成功插入"
	c.Data["json"]=resp
	//传出的方式
	c.ServeJSON()
	return
}


//展示我的购物车
func (c *ShoppingcartController)ShowMyShoppingCart()  {
	//获取数据
	userName:=c.GetSession("userName")
	//根据用户名查询到全部商品
	Shoppinggoods,err:=models.RedisShoppingcartshow(userName.(string))
	if err!=nil {
		c.Data["Shoppinggoods"]=""
	}else{
		c.Data["Shoppinggoods"]=Shoppinggoods
	}
	c.TplName="cart.html"
}

//删除用户的一个行购物车
func (c *ShoppingcartController)DeleteShoppingCart()  {
	//获取数据
	userName:=c.GetSession("userName")
	skuid,err:=c.GetInt("skuid")
	resp:=make(map[string]interface{})
	if err!=nil {
		//赋值
		resp["status"]=400
		resp["msg"]="没有这个商品id"
		//成json格式
		c.Data["json"]=resp
		//返回数据方式
		c.ServeJSON()
		return
	}

	if userName==nil {
		//赋值
		resp["status"]=401
		resp["msg"]="您没有登录状态了"
		//成json格式
		c.Data["json"]=resp
		//返回数据方式
		c.ServeJSON()
		return
	}
	//在数据库
	err=models.RedisDeleteShoppingcartshow(userName.(string),skuid)
	if err!=nil {
		resp["status"]=402
		resp["msg"]="数据库没有这个数据"
		//成json格式
		c.Data["json"]=resp
		//返回数据方式
		c.ServeJSON()
		return
	}
	resp["status"]=200
	resp["msg"]="成功"
	//成json格式
	c.Data["json"]=resp
	//返回数据方式
	c.ServeJSON()
	return
}


