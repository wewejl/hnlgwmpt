package models

import (
	"github.com/astaxie/beego/orm"
	"errors"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"encoding/json"

)

//插入用户信息
func IsertUser(username string,PassWord string) error  {
	//获取结构体
	var user User

	//获取orm数据库
	o:=orm.NewOrm()
	//赋值
	user.Name=username
	user.PassWord=PassWord
	//插入
	_,err:=o.Insert(&user)
	return err
}

//处理邮箱激活
func UpdataUserEmail(username string,email string) error {
	//获取结构体
	var user User
	//获取orm数据库对象
	o:=orm.NewOrm()
	//赋值
	user.Name=username
	//查询
	err:=o.Read(&user,"Name")
	if err!=nil {
		return err
	}
	//赋新值
	user.Active=true
	user.Email=email
	//更新操作
	_,err=o.Update(&user,"Active","Email")
	return err
}

//处理登录   查询数据库
func SeclectLogin(userName string,password string) error {
	//获取结构体
	var user User
	//获取orm数据库对象
	o:=orm.NewOrm()
	//赋值
	user.Name=userName
	//查询
	err:=o.Read(&user,"Name")
	if err!=nil {
		return  errors.New("用户密码错误")
	}
	//校验
	if user.PassWord!=password {
		return errors.New("用户密码错误")
	}
	if !user.Active {
		return errors.New("用户未激活")
	}
	return nil
}

//根据用户名 查询用户信息
func SeclectUserinfo(userName string) ([]Address,error) {
	//获取结构体
	var address []Address
	//获取数据库orm对象
	//var user User

	o:=orm.NewOrm()
	//赋值

	_,err:=o.QueryTable("Address").Filter("User__Name",userName).All(&address)
	return address,err
}

//插入默认地址
func InsertPersonalcentersite(receiptname,detailedaddress,zipCode,phone string,userName string) error {
	//获取结构体
	var address Address
	//获取orm数据库对象
	o:=orm.NewOrm()
	//获取user对象
	var user User
	user.Name=userName
	//查询
	err:=o.Read(&user,"Name")
	if err!=nil {
		return errors.New("没有查到对应的user")
	}
	//根据user删除所有对应的Isdefault成0
	o.QueryTable("Address").Filter("User",&user).Update(orm.Params{"Isdefault":0})
	//在赋值
	address.Receiver=receiptname
	address.Addr=detailedaddress
	address.Zipcode=zipCode
	address.Phone=phone
	address.Isdefault=1
	address.User=&user
	_,err=o.Insert(&address)
	if err!=nil {
		return errors.New("插入失败")
	}
	return nil

}

//查询用户默认地址
func SelectPersonalcentersite(userName string) (Address,error) {
	//获取结构体
	var address Address
	//获取数据库orm对象
	o:=orm.NewOrm()
	//高级查询
	err:=o.QueryTable("Address").RelatedSel("User").Filter("User__Name",userName).Filter("Isdefault",1).One(&address)

	return address,err
}

//查询redis用户最近的浏览记录
func SelectRedisDetailinfo(userName string) []GoodsSKU {
	conn,_:=redis.Dial("tcp","192.168.73.128:6379")
	var goodssku GoodsSKU
	var goodsskus []GoodsSKU
	goods,err:=redis.Values(conn.Do("lrange",userName+"_hostory",0,4))
	if err!=nil {
		fmt.Println(err)
		return nil
	}
	for _,v:=range goods {
		json.Unmarshal(v.([]byte),&goodssku)
		goodsskus=append(goodsskus,goodssku)
	}
	return goodsskus


}