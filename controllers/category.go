package controllers

import (
	"github.com/astaxie/beego"
	"github.com/lsyslsy/beeblog/models"
)

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get() {

	op := c.Input().Get("op")
	switch op {
	case "add":
		name := c.Input().Get("name")
		if len(name) == 0 {
			break
		}
		err := models.AddCategory(name)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/category", 302)
		return
	case "del":
		id := c.Input().Get("id")
		if len(id) == 0 {
			break
		}

		err := models.DelCategory(id)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/category", 302)
		return
	}
	c.TplNames = "category.html"
	c.Data["IsCategory"] = true

	var err error
	c.Data["Categories"], err = models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
}
