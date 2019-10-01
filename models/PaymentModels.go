package models

import (
	"strconv"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"time"
	"strings"
	"errors"
	"fmt"
)

func PaymentCloseAccount(skuids []string, userName string) ([]map[string]interface{}, int, int, error) {
	//循环出skuids
	var goodsSKUscart []map[string]interface{}
	//链接redis
	conn, err := redis.Dial("tcp", "192.168.73.128:6379")
	if err != nil {
		return nil, 0, 0, err
	}
	//创建一个总价
	var llroparic int
	//创建一个购买商品总数
	var llrocontent int

	var goodsSKU GoodsSKU
	for _, id := range skuids {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return nil, 0, 0, err
		}
		temp := make(map[string]interface{})
		goodsSKU.Id = idInt
		//获取orm数据库对象
		o := orm.NewOrm()
		err = o.Read(&goodsSKU, "Id")
		if err != nil {
			return nil, 0, 0, err
		}
		count, err := redis.Int(conn.Do("hget", userName+"_host", idInt))
		if err != nil {
			return nil, 0, 0, err
		}
		temp["llpairc"] = count * goodsSKU.Price
		temp["goodsSKU"] = goodsSKU
		temp["count"] = count
		llroparic += count * goodsSKU.Price
		llrocontent += count
		goodsSKUscart = append(goodsSKUscart, temp)
	}
	return goodsSKUscart, llroparic, llrocontent, nil
}

//插入订单表
func InsertOrderInfoList(userName string, addrId, payId int, skuids string, totalCount, totalPrice, transit int) (error, string) {
	//链接redis数据库
	conn, err := redis.Dial("tcp", "192.168.73.128:6379")
	if err != nil {
		return err, "redis数据库链接错误"
	}
	//获取orm数据库里面对象
	o := orm.NewOrm()
	//获取结构体
	var orderinfo OrderInfo
	//赋值
	//--付款方式
	orderinfo.PayMethod = payId
	//--商品数量
	orderinfo.TotalCount = totalCount
	//--商品总价
	orderinfo.TotalPrice = totalPrice
	//--快递费用
	orderinfo.TransitPrice = transit
	//用户的对象指针
	var user User
	user.Name = userName
	//查询user用户
	err = o.Read(&user, "Name")
	if err != nil {
		return err, "查询user用户错误"
	}
	var address Address
	address.Id = addrId
	err = o.Read(&address)
	if err != nil {
		return err, "查询address地址表错误"
	}
	//orderinfo赋值
	orderinfo.Address = &address
	orderinfo.User = &user
	//把订单号添加上去
	orderId := time.Now().Format("200601021504051234") + strconv.Itoa(user.Id)
	orderinfo.OrderId = orderId
	//开启回滚
	o.Begin()
	_, err = o.Insert(&orderinfo)
	if err != nil {
		o.Rollback()
		return err, "插入订单表"
	}
	//处理skuids数据
	//掐头去尾 skuids[1:len(skuids)-1]
	skuqieb := strings.Split(skuids[1:len(skuids)-1], " ")
	for _, cartId := range skuqieb {
		var ordergoods OrderGoods
		ordergoods.OrderInfo = &orderinfo
		//循环数据

		//根据id查到商品
		var goodssku GoodsSKU
		var count int
		for i := 10; ; {
			//把商品id进行转化
			Idint, err := strconv.Atoi(cartId)
			if err != nil {
				o.Rollback()
				return err, "id字符转成int类型"
			}
			goodssku.Id = Idint
			err = o.Read(&goodssku, "Id")
			if err != nil {
				o.Rollback()
				return err, "循环读goodssku"
			}
			historyStock := goodssku.Stock
			time.Sleep(10 * time.Second)
			o.Read(&goodssku)
			count, err = redis.Int(conn.Do("hget", userName+"_host", Idint))
			if err != nil {
				o.Rollback()
				return err, "redis读哈希表 商品数量错误"
			}
			fmt.Println("历史库存为,", historyStock, "现在库存为,", goodssku.Stock)
			ordergoods.GoodsSKU = &goodssku
			ordergoods.Count = count
			ordergoods.Price = count * goodssku.Price
			if goodssku.Stock < count {
				o.Rollback()
				return errors.New("店家库存不够"), "店家的" + goodssku.Name + "库存不够"
			}
			if goodssku.Stock != historyStock {
				if i >= 0 {
					i -= 1
				} else {
					o.Rollback()
					return errors.New("店家库存不够"), "店家的" + goodssku.Name + "库存不够"
				}
			} else {
				conn.Do("hdel", userName+"_host", Idint)
				break
			}
		}
		_, err = o.Insert(&ordergoods)

		if err != nil {
			o.Rollback()
			return err, "插入订单商品表错误"
		}
		goodssku.Stock -= count
		goodssku.Sales += count
		o.Update(&goodssku, "Stock", "Sales")
	}
	//插入订单商品表
	o.Commit()
	return nil, ""
}
