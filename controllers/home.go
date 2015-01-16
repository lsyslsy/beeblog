package controllers

import (
	"github.com/astaxie/beego"
	"github.com/lsyslsy/beeblog/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["IsHome"] = true
	c.TplNames = "home.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	topics, err := models.GetAllTopics(c.Input().Get("cate"), c.Input().Get("label"), true)
	if err != nil {
		beego.Error(err)
		return
	} else {
		c.Data["Topics"] = topics
	}
	c.Data["Categories"], err = models.GetAllCategories()
	if err != nil {
		beego.Error(err)
		return
	}

}
