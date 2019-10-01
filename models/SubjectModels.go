package models

import (
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"

)

//获取全部的联机菜单
func SeclectAllfrom() []map[string]interface{} {
	//获取orm对象
	o:=orm.NewOrm()
	var munes []TpshopCategory
	var goodsTypes []map[string]interface{}
	//高级查询
	//o.QueryTable("TpshopCategory").Filter("Pid",0).All(&menus)
	o.QueryTable("TpshopCategory").Filter("Pid",0).All(&munes)
	//循环munes

	for _,v1:=range munes{
		temp:=make(map[string]interface{})
		var erji []TpshopCategory
		o.QueryTable("TpshopCategory").Filter("Pid",v1.Id).All(&erji)
		temp["yijidangge"]=v1
		temp["erji"]=erji
		goodsTypes=append(goodsTypes,temp)
	}
	//循环二级
	for _,v1:=range goodsTypes{
		//for _,v2 := range v1["erji"].([]TpshopCategory){
		var erjiTemp  []map[string]interface{}
		for _,v:=range  v1["erji"].([]TpshopCategory){
			temp:=make(map[string]interface{})
			var shanji []TpshopCategory
			o.QueryTable("TpshopCategory").Filter("Pid",v.Id).All(&shanji)
			temp["sanji"]=shanji
			temp["erjidangge"]=v
			erjiTemp=append(erjiTemp,temp)
		}
		v1["二级的值和三级的切片"]=erjiTemp
	}
	return goodsTypes
}


//获取所有的商品类型
func SeclectAllCommodityType() ([]GoodsType,error) {

	//获取orm数据库对象
	o:=orm.NewOrm()
	//获取结构体
	var goodstype []GoodsType
	//高级查询
	_,err:=o.QueryTable("GoodsType").All(&goodstype)
	return goodstype,err
}

//获取所有轮播图
func SeclectAllBanner() ([]IndexGoodsBanner,error) {
	//获取orm数据库对象
	o:=orm.NewOrm()
	//获取结构体
	var indexgoodsbanner []IndexGoodsBanner
	//高级查询
	_,err:=o.QueryTable("IndexGoodsBanner").RelatedSel("GoodsSKU").All(&indexgoodsbanner)
	return indexgoodsbanner,err
}

//获取促销商品类型
func SeclectAllPromotion() ([]IndexPromotionBanner,error) {
	//获取orm数据库对象
	o:=orm.NewOrm()
	//获取结构提
	var indexpromotionBanner []IndexPromotionBanner
	//高级查询
	_,err:=o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&indexpromotionBanner)
	return indexpromotionBanner,err
}

//获取首页分类商品展示表
func SelectIndexTypeGoodsBanner(GoodsType []GoodsType) ([]map[string]interface{}) {
	//获取orm数据库对象
	o:=orm.NewOrm()

	//获取goodsTyps
	var goodsTypeGoodsBanners []map[string]interface{}

	//获取结构体切片


	//高级查询
	for _,v:=range GoodsType{
		//获取图片
		temp:=make(map[string]interface{})

		var TextindexTypeGoodsBanner []IndexTypeGoodsBanner
		var ImgeindexTypeGoodsBanner []IndexTypeGoodsBanner
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType").RelatedSel("GoodsSKU").
			Filter("GoodsType__Id",v.Id).
			Filter("DisplayType",0).OrderBy("Index").All(&TextindexTypeGoodsBanner)

		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType").RelatedSel("GoodsSKU").
			Filter("GoodsType__Id",v.Id).
			Filter("DisplayType",1).OrderBy("Index").All(&ImgeindexTypeGoodsBanner)
			temp["GoodsTypes"]=v
			temp["TextGoods"]=TextindexTypeGoodsBanner
			temp["ImgeGoods"]=ImgeindexTypeGoodsBanner
		goodsTypeGoodsBanners=append(goodsTypeGoodsBanners,temp)
	}
	fmt.Println(goodsTypeGoodsBanners[1])
	return goodsTypeGoodsBanners
}

//根据id获取商品的全部数据
func SelectGoodsAllDetails(id int) (GoodsSKU,error) {
	//获取结构体
	var goodssku GoodsSKU
	//获取orm数据库对象
	o:=orm.NewOrm()
	//赋值
	goodssku.Id=id
	//查询
	err:=o.QueryTable("GoodsSKU").RelatedSel("GoodsType","Goods").Filter("Id",id).One(&goodssku)
	return goodssku,err
}

//根据类型获取类型商品的时间最新的前两个数据
func SelectGoodsAllrecommen(GoodsType string) ([]GoodsSKU,error) {
	//获取orm数据库对象
	o:=orm.NewOrm()
	//获取结构体
	var goodssku []GoodsSKU
	//高级查询
	_,err:=o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name",GoodsType).
		OrderBy("-Time").Limit(2,0).All(&goodssku)
	return goodssku,err
}

//根据类型获取类型商品的全部数据
func SelectGoodsTypeAlllist(GoodsTypeName string,Pagesize,pageindex int) (int64,[]GoodsSKU,error) {
	//获取orm数据库对象
	o:=orm.NewOrm()
	//获取结构体
	var goodssku []GoodsSKU

	//高级查询
	DataNumber,err:=o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name",GoodsTypeName).OrderBy("Id").Limit(Pagesize,(pageindex-1)*Pagesize).All(&goodssku)
	return DataNumber,goodssku,err
}

//根据价格排序查询类型商品的全部数据
func SelectGoodsTypeAlllistprice(GoodsTypeName string,price string,Pagesize,pageindex int) (int64,[]GoodsSKU,error, string)  {
	//获取orm数据库对象
	o:=orm.NewOrm()
	//获取结构体
	var goodssku []GoodsSKU
	if price=="true" {
		DataNumber,err:=o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name",GoodsTypeName).OrderBy("Price").Limit(Pagesize,(pageindex-1)*Pagesize).All(&goodssku)
		price="false"
		return DataNumber,goodssku,err,price
	}else {
		DataNumber,err:=o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name",GoodsTypeName).OrderBy("-Price").Limit(Pagesize,(pageindex-1)*Pagesize).All(&goodssku)
		price="true"
		return DataNumber,goodssku,err,price
	}
	//高级查询

}

//根据人气排序查询类型商品的全部数据
func SelectGoodsTypeAlllistsentiment(GoodsTypeName string,sentiment string,Pagesize,pageindex int) (int64,[]GoodsSKU,error, string)  {
	//获取orm数据库对象
	o:=orm.NewOrm()
	//获取结构体
	var goodssku []GoodsSKU
	if sentiment=="true" {
		DataNumber,err:=o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name",GoodsTypeName).Limit(Pagesize,(pageindex-1)*Pagesize).OrderBy("Sales").All(&goodssku)
		sentiment="false"
		return DataNumber,goodssku,err,sentiment
	}else {
		DataNumber,err:=o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name",GoodsTypeName).Limit(Pagesize,(pageindex-1)*Pagesize).OrderBy("-Sales").All(&goodssku)
		sentiment="true"
		return DataNumber,goodssku,err,sentiment
	}
	//高级查询
}

//根据类型查询到所有类型数据多少
func SelectgoodsTypePageDataNumber(GoodsTypeName string) (PageDataNumber int64) {
	//获取数据库对象
	o:=orm.NewOrm()
	//获取结构体
	PageDataNumber,_=o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name",GoodsTypeName).Count()
	return
}

//最近浏览的时候进行redis插入数据
func SelectUserGoodsinfoDetail(userName string,id int) error {
	//现在是在浏览这个商品
	 o:=orm.NewOrm()
	 var goodssku GoodsSKU
	 //获取数据库
	 err:=o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("Id",id).One(&goodssku)
	if err!=nil {
		return err
	}
	conn,_ := redis.Dial("tcp","192.168.73.128:6379")
	goodsskubyus,err:=json.Marshal(&goodssku)
	if err!=nil {
		return err
	}
	//在添加的时候删除
	conn.Do("lrem",userName+"_hostory",0,goodsskubyus)
	_,err=conn.Do("lpush",userName+"_hostory",goodsskubyus)
	return err

}

//搜索的是全部商品信息
func SeclectAllGoosSrcre() ([]GoodsSKU,error) {
	//获取orm数据库对象
	o:=orm.NewOrm()
	//获取结构体
	var goodssku []GoodsSKU
	//高级查询
	_,err:=o.QueryTable("GoodsSKU").All(&goodssku)
	return goodssku,err
}

//搜索的是搜索内容的全部商品信息
func SeclectGoodsNameGoosSrcre(GoodsName string) ([]GoodsSKU,error) {
	//获取orm的数据库向
	o:=orm.NewOrm()
	//获取结构体
	var goodssku []GoodsSKU
	//高级操作
	_,err:=o.QueryTable("GoodsSKU").Filter("Name__contains",GoodsName).All(&goodssku)
	return goodssku,err
}