package main

import (
	_ "ttsx/routers"
	_ "ttsx/models"
	"github.com/astaxie/beego"


)

func main() {
	beego.Run()
}

