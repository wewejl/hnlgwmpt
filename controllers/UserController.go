package controllers

import (
	"github.com/astaxie/beego"
	"ttsx/models"
	"fmt"
	"math/rand"
	"time"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/astaxie/beego/utils"
)

type UserController struct {
	beego.Controller
}

//展示注册页面
func (c *UserController) ShowRegister() {
	//展示页面
	c.TplName = "register.html"
}

//发送短信业务
func (c *UserController) Sendmessage() {

	//获取电话号码
	phone := c.GetString("phone")
	//种下随机数种子
	fmt.Println("到这里了")
	rand.Seed(time.Now().UnixNano())
	//创建随机数
	ret := rand.Intn(1000000)
	num := fmt.Sprintf("%06v", ret)
	fmt.Println(num)
	c.Ctx.SetCookie(phone+"_user", string(num), 60*60*24)

	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAI4FvKwHjXNkrNuU7jpAKP", "tKAY2OIYOvgCvtC8eHGR3cad4sTutY")

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = phone
	request.SignName = "学院水果超市注册"
	request.TemplateCode = "SMS_174275505"
	request.TemplateParam = "{\"code\":\"" + num + "\"}"

	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	stal := make(map[string]interface{})

	if response.IsSuccess() {
		//返回的是成功的
		fmt.Println(response.IsSuccess())
		//传给前段数据
		stal["start"] = 200
		stal["msg"] = "ok"

	} else {
		stal["start"] = 500
		stal["msg"] = "验证码发送错误"
	}
	//指定传输方式
	c.Data["json"] = stal
	c.ServeJSON();
}

//注册用户
func (c *UserController) InsertRegister() {
	//获取数据
	phone := c.GetString("phone")
	code := c.GetString("code")
	password := c.GetString("password")
	repassword := c.GetString("repassword")
	//校验数据
	if phone == "" || code == "" || password == "" || repassword == "" {
		c.Data["errer"] = "输入数据不完整"
		c.TplName = "register.html"
		return
	}
	codeatl := c.Ctx.GetCookie(phone + "_user")
	fmt.Println("codeatl =",codeatl)
	if code != codeatl {
		c.Data["phone"] = phone
		c.Data["code"] = code
		c.Data["password"] = password
		c.Data["repassword"] = repassword
		c.Data["errer"] = "验证码错误"
		c.TplName = "register.html"
		return
	}
	if password != repassword {
		c.Data["phone"] = phone
		c.Data["code"] = code
		c.Data["password"] = password
		c.Data["repassword"] = repassword
		c.Data["errer"] = "两个密码不相同"
		c.TplName = "register.html"
		return
	}
	err := models.IsertUser(phone, password)
	if err != nil {
		c.Data["errer"] = "用户插入失败"
		c.TplName = "register.html"
		return
	}
	c.Redirect("/registerEmail?userName="+phone, 302)
}

//展示邮箱激活页面
func (c *UserController) ShowRegisterEmail() {
	//获取数据
	userName := c.GetString("userName")
	//校验数据
	if userName == "" {
		c.Ctx.WriteString("userName 没有传过来")
	}
	//处理数据
	c.Data["userName"] = userName
	//返回视图
	c.TplName = "register-email.html"
}

//处理邮箱激活
func (c *UserController) InsertRegisterEmail() {
	//获取数据
	username := c.GetString("userName")
	email := c.GetString("email")
	//校验数据
	if username == "" || email == "" {
		c.Data["errer"] = "input的username没有获取到"
		c.TplName = "register-email.html"
		return
	}
	fmt.Println("username =", username)
	fmt.Println("email =", email)
	//处理数据
	fmt.Println("1")

	// 创建一个字符串变量，存放相应配置信息
		config := `{"username":"zhu1024344053@163.com","password":"zhuxinyebj5q","host":"smtp.163.com","port":25}`
	//通过字符串创建一个email对象
	temail := utils.NewEMail(config)
	//邮箱名称
	temail.Subject = "学校天天生鲜邮箱激活"
	//邮箱的发送着
	temail.From = "zhu1024344053@163.com"
	//邮箱接受者
	temail.To = []string{email}
	//邮箱内容
	temail.HTML = `<html>
	<head>
	</head>
	<body>
	<div>点击下面链接即可完成激活 </div>	<a href="http://192.168.73.128:8080/emailseed?username=` + username + `&email=` + email + `">点击激活</a>
	</body>
	</html>`
	err := temail.Send()
	fmt.Println("4")
	if err != nil {
		c.Ctx.WriteString("发生错误,邮箱未发送")
		return
	}
	c.Ctx.WriteString("邮箱已发,请前往邮箱进行验证")
}

//处理邮箱激活页面
func (c *UserController) Emailseed() {
	//获取数据
	username := c.GetString("username")
	email := c.GetString("email")
	//校验数据
	if username == "" || email == "" {
		c.Data["errer"] = "邮箱激活链接失败"
		c.TplName = "reigster-email.html"
		return
	}
	//处理数据
	err := models.UpdataUserEmail(username, email)
	if err != nil {
		c.Data["errer"] = "邮箱激活链接失败"
		c.TplName = "reigster-email.html"
		return
	}
	c.Data["errer"] = "激活成功"
	c.Data["userName"] = c.GetString("userName")
	c.Redirect("/login?userName="+username, 302)
}

//展示登录界面
func (c *UserController) ShowLogin() {
	//获取数据
	userName := c.GetString("userName")
	//返回视图
	c.Data["userName"] = userName
	c.TplName = "login.html"
}

//处理登录页面
func (c *UserController) SeclectLogin() {
	//获取数据
	userName := c.GetString("userName")
	PassWord := c.GetString("PassWord")
	//校验数据
	fmt.Println("到这里了")
	if userName == "" || PassWord == "" {
		c.Data["errer"] = "输入数据不完整"
		c.Data["userName"] = userName
		c.Data["PassWord"] = PassWord
		c.TplName = "login.html"
		return
	}
	//处理数据
	err := models.SeclectLogin(userName, PassWord)
	if err != nil {
		c.Data["errer"] = err
		c.Data["userName"] = userName
		c.Data["PassWord"] = PassWord
		c.TplName = "login.html"
		return
	}
	//登录成功后进行登录状态更新
	c.SetSession("userName", userName)
	//返回视图
	c.Redirect("/index", 302)

}

//用户注销
func (c *UserController) LayoutSseion() {
	//从新设置Session
	c.DelSession("userName")
	//返回视图
	c.Redirect("/login", 302)
}

//展示个人中心info页面
func (c *UserController) ShowPersonalcenterinfo() {
	//获取数据
	userName := c.GetSession("userName")
	//校验数据
	if userName == nil {
		userName = ""
	}
	//处理数据
	goodssku:=models.SelectRedisDetailinfo(userName.(string))
	if goodssku!=nil {
		c.Data["goodssku"]=goodssku
	}else {
		c.Data["goodssku"]=""
	}
	//返回视图

	c.Data["userName"] = userName.(string)
	c.Layout = "user_conterlayout.html"
	c.TplName = "user_center_info.html"
}

//展示个人中心order页面
func (c *UserController) ShowPersonalcenterorder() {
	c.Layout = "user_conterlayout.html"
	c.TplName = "user_center_order.html"
}

//展示个人中site页面
func (c *UserController) ShowPersonalcentersite() {
	//获取数据
	userName:=c.GetSession("userName")
	//校验数据
	//处理数据
	address,err:=models.SelectPersonalcentersite(userName.(string))
	if err!=nil {
		c.Data["address"]=""
	}else {
		phonstart:=address.Phone[:3]
		phonend:=address.Phone[7:]
		address.Phone=phonstart+"****"+phonend
		c.Data["address"]=address
	}
	c.Layout = "user_conterlayout.html"
	c.TplName = "user_center_site.html"
}

//插入默认地址
func (c *UserController) InsertPersonalcentersite() {
	//获取数据
	userName := c.GetSession("userName")
	receiptname := c.GetString("receiptname")
	detailedaddress := c.GetString("detailedaddress")
	zipCode := c.GetString("zipCode")
	phone := c.GetString("phone")
	//校验数据
	if receiptname == "" || detailedaddress == "" || zipCode == "" || phone == "" {
		c.Data["errer"]="输入数据不完整"
		c.Data["userName"]=userName
		c.Layout="user_conterlayout.html"
		c.TplName = "user_center_site.html"
		return
	}
	//处理数据
	err:=models.InsertPersonalcentersite(receiptname,detailedaddress,zipCode,phone,userName.(string))
	if err != nil {
		c.Data["errer"]=err
		c.Data["userName"]=userName
		c.Layout="user_conterlayout.html"
		c.TplName = "user_center_site.html"
		return
	}
	//返回视图
	c.Redirect("/user/Personalcentersite",302)
}
