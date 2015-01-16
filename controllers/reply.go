package controllers

import (
	"github.com/astaxie/beego"
	"github.com/lsyslsy/beeblog/models"
)

type ReplyController struct {
	beego.Controller
}

func (c *ReplyController) Add() {
	tid := c.Input().Get("tid")
	err := models.AddReply(tid,
		c.Input().Get("nickname"), c.Input().Get("content"))
	if err != nil {
		beego.Error(err)
		c.Redirect("/topic/view/"+tid, 302)
		return

	}

	c.Redirect("/topic/view/"+tid, 302)
}

func (c *ReplyController) Delete() {
	if !checkAccount(c.Ctx) {
		return
	}
	rid := c.Input().Get("rid") // reply id
	err := models.DeleteReply(rid)
	if err != nil {
		beego.Error(err)
	}
	tid := c.Input().Get("tid")
	c.Redirect("/topic/view/"+tid, 302)
}
