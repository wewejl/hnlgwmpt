package controllers

import (
	"github.com/astaxie/beego"
	"ttsx/models"
	"fmt"
	"math"
)

type SubjectCOntroller struct {
	beego.Controller
}

//展示index页面
func (c *SubjectCOntroller) ShowIndex() {
	//获取数据
	userName := c.GetSession("userName")
	//校验数据
	if userName != nil {
		c.Data["userName"] = userName.(string)
	} else {
		c.Data["userName"] = ""
	}
	//处理数据
	goodsTypes := models.SeclectAllfrom()

	//返回视图
	c.Data["goodsTypes"] = goodsTypes

	c.TplName = "index.html"
}

//展示indexSx页面
func (c *SubjectCOntroller) ShowIndexSx() {
	fmt.Println("奥斯卡的减肥啦开始的缴费卡收到付款了")
	//获取所有商品类型
	GoodsType, err := models.SeclectAllCommodityType()
	if err != nil {
		c.Data["errer"] = "数据库数据获取错误"
		c.TplName = "index_sx.html"
		return
	}
	//获取所有轮播图
	IndexGoodsBanner, err1 := models.SeclectAllBanner()
	//获取所有促销类型
	IndexPromotionBanner, err2 := models.SeclectAllPromotion()
	if err != nil || err1 != nil || err2 != nil {
		c.Data["errer"] = "数据库数据获取错误"
		c.TplName = "index_sx.html"
		return
	}

	//获取所有商品IndexTypeGoodsBanner
	indexTypeGoodsBanner := models.SelectIndexTypeGoodsBanner(GoodsType)

	c.Data["IndexTypeGoodsBanner"] = indexTypeGoodsBanner
	c.Data["GoodsType"] = GoodsType
	c.Data["IndexGoodsBanner"] = IndexGoodsBanner
	c.Data["IndexPromotionBanner"] = IndexPromotionBanner
	c.TplName = "index_sx.html"
}

//展示Details页面
func (c *SubjectCOntroller) ShowDetails() {
	//获取数据
	id, err := c.GetInt("Id")
	if err != nil {
		c.Redirect("/indexSx", 302)
		return
	}
	//根据id获取这个商品的全部数据
	goodssku, err := models.SelectGoodsAllDetails(id)
	if err != nil {
		c.Ctx.WriteString("数据库查询错误")
		return
	}
	userName := c.GetSession("userName")
	if userName != nil {
		err := models.SelectUserGoodsinfoDetail(userName.(string), id)
		if err != nil {
			fmt.Println(err)
			c.Ctx.WriteString("redis查入错误")
			return
		}
	}
	//根据类型查询到新品商品
	newgoodssku, err := models.SelectGoodsAllrecommen(goodssku.GoodsType.Name)
	if err != nil {
		c.Ctx.WriteString("新品数据错误")
		return
	}
	c.Data["newgoodssku"] = newgoodssku
	c.Data["goodssku"] = goodssku
	c.TplName = "detail.html"
}

//展示list页面
func (c *SubjectCOntroller) ShowTypelist() {
	//获取数据
	GoodsTypeName := c.GetString("GoodsTypeName")
	price := c.GetString("price")
	sentiment := c.GetString("sentiment")
	pageindex, err := c.GetInt("pageindex")
	Pagesize := 2
	if err != nil {
		pageindex = 1
	}
	if price == "" {
		price = "true"
	}
	//处理数据
	//根据类型查询到全部商品
	var DataNumber int64
	if price == "" && sentiment == "" {
		_, goodssku, err := models.SelectGoodsTypeAlllist(GoodsTypeName, Pagesize, pageindex)
		if err != nil {
			c.Ctx.WriteString("对不起查询失败")
			return
		}

		//根据类型查询到新品商品
		newgoodssku, err := models.SelectGoodsAllrecommen(GoodsTypeName)
		if err != nil {
			c.Ctx.WriteString("新品数据错误")
			return
		}
		c.Data["newprice"] = price
		c.Data["newsentiment"] = "true"
		c.Data["newgoodssku"] = newgoodssku
		c.Data["goodsskus"] = goodssku

	} else if price != "" && sentiment == "" {
		_, goodssku, err, _ := models.SelectGoodsTypeAlllistprice(GoodsTypeName, price, Pagesize, pageindex)
		if err != nil {
			c.Ctx.WriteString("对不起查询失败")
			return
		}

		//根据类型查询到新品商品
		newgoodssku, err := models.SelectGoodsAllrecommen(GoodsTypeName)
		if err != nil {
			c.Ctx.WriteString("新品数据错误")
			return
		}
		c.Data["newprice"] = price
		c.Data["newsentiment"] = "true"
		c.Data["goodsskus"] = goodssku
		c.Data["newgoodssku"] = newgoodssku

	} else if price == "" && sentiment != "" {
		_, goodssku, err, newsentiment := models.SelectGoodsTypeAlllistsentiment(GoodsTypeName, sentiment, Pagesize, pageindex)
		if err != nil {
			c.Ctx.WriteString("对不起查询失败")
			return
		}

		//根据类型查询到新品商品
		newgoodssku, err := models.SelectGoodsAllrecommen(GoodsTypeName)
		if err != nil {
			c.Ctx.WriteString("新品数据错误")
			return
		}

		c.Data["newprice"] = "true"
		c.Data["newsentiment"] = newsentiment
		c.Data["goodsskus"] = goodssku
		c.Data["newgoodssku"] = newgoodssku
	} else {
		_, goodssku, err, newsentiment := models.SelectGoodsTypeAlllistsentiment(GoodsTypeName, sentiment, Pagesize, pageindex)
		if err != nil {
			c.Ctx.WriteString("对不起查询失败")
			return
		}

		//根据类型查询到新品商品
		newgoodssku, err := models.SelectGoodsAllrecommen(GoodsTypeName)
		if err != nil {
			c.Ctx.WriteString("新品数据错误")
			return
		}

		c.Data["newprice"] = "true"
		c.Data["newsentiment"] = newsentiment
		c.Data["goodsskus"] = goodssku
		c.Data["newgoodssku"] = newgoodssku
	}

	//页码处理
	DataNumber = models.SelectgoodsTypePageDataNumber(GoodsTypeName)

	pageCount := math.Ceil(float64((DataNumber)) / float64(Pagesize))
	pageqie := ShowPageIndex(int(pageCount), pageindex)
	//返回视图
	var pageStart int
	var pageEnd int
	if pageindex == 1 {
		pageStart = 1
		pageEnd = pageindex + 1
	} else if pageindex == int(pageCount) {
		pageStart = pageindex - 1
		pageEnd = pageindex
	} else {
		pageStart = pageindex - 1
		pageEnd = pageindex + 1
	}
	c.Data["pageStart"] = pageStart
	c.Data["pageEnd"] = pageEnd
	c.Data["pageindex"] = pageindex
	c.Data["pageqie"] = pageqie
	c.Data["GoodsTypeName"] = GoodsTypeName
	c.TplName = "list.html"
}

func ShowPageIndex(pageCount int, pageindex int) (pageqie []int) {
	//页码数据＜5
	if pageCount < 5 {
		for i := 1; i <= pageCount; i++ {
			pageqie = append(pageqie, i)
		}
		return pageqie
	}
	if pageindex <= 3 {
		for i := 1; i <= 5; i++ {
			pageqie = append(pageqie, i)
		}
		return pageqie
	}
	if pageindex < pageCount-2 {
		for i := pageindex - 2; i <= pageindex+2; i++ {
			pageqie = append(pageqie, i)
		}
		return pageqie
	}
	for i := pageCount - 4; i <= pageCount; i++ {
		pageqie = append(pageqie, i)
	}
	return pageqie
}

//搜索商品名称
func (c *SubjectCOntroller) ShowgoodsSrcre() {
	//获取数据
	GoodsName := c.GetString("GoodsName")
	//校验数据
	if GoodsName == "" {
		//获取全部
		goodssku, err := models.SeclectAllGoosSrcre()
		if err != nil {
			c.Data["goodssku"] = ""
		} else {
			c.Data["goodssku"] = goodssku
		}
	} else {
		//获取名称
		goodssku, err := models.SeclectGoodsNameGoosSrcre(GoodsName)
		if err != nil {
			c.Data["goodssku"] = ""
		} else {
			c.Data["goodssku"] = goodssku
		}
	}
	c.TplName = "srcre.html"
}
