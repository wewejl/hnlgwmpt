package models

import (
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego/orm"
	"strconv"
	"fmt"
)

//redis传值到购物车上面
func RedisShoppingcart(userName string,skuid,count int) error {
	//链接redis数据库
	conn,err:=redis.Dial("tcp","192.168.73.128:6379")
	if err!=nil {
		return err
	}
	rediscount,err:=redis.Int(conn.Do("hget",userName+"_host",skuid))
	if err!=nil {
		_,err=conn.Do("hset",userName+"_host",skuid,count)
		return err
	}else{
		count=count+rediscount
		_,err=conn.Do("hset",userName+"_host",skuid,count)
		return err
	}
}

//查询这个用户的全部redis购物车的数据
func RedisShoppingcartshow(userName string) ([]map[string]interface{},error)  {

	//创建一个大容器
	var  goods []map[string]interface{}
	//链接redis数据库
	conn,err:=redis.Dial("tcp","192.168.73.128:6379")
	if err!=nil {
		return nil,err
	}
	fmt.Println(userName+"_host")
	newmap,err:=redis.IntMap(conn.Do("hgetall",userName+"_host"))
	fmt.Println("newmap =",newmap)
	if err!=nil {
		return nil,err
	}
	for id,count:=range newmap {
		temp:=make(map[string]interface{})
		//把id进行转化成int类型
		newId,err:=strconv.Atoi(id)
		if err!=nil {
			return nil,err
		}
		//获取orm数据库对象
		o:=orm.NewOrm()
		//获取结构体
		var goodsSKU GoodsSKU
		//赋值
		goodsSKU.Id=newId
		//查询
		err=o.Read(&goodsSKU,"Id")
		if err!=nil {
			return nil,err
		}
		llCount:=goodsSKU.Price*count
		temp["llCount"]=llCount
		temp["GoodsSKU"]=goodsSKU
		temp["Count"]=count
		goods=append(goods,temp)
	}
	return goods,nil
}

//删除用户redis购物车的行数据
func RedisDeleteShoppingcartshow(userName string,skuid int) error {
	//用户的数据
	//链接redis数据库
	conn,_:=redis.Dial("tcp","192.168.73.128:6379")
	_,err:=conn.Do("hdel",userName+"_host",skuid)
	return err
}
