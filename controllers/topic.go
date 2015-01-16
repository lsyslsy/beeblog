package controllers

import (
	"github.com/astaxie/beego"
	"github.com/lsyslsy/beeblog/models"
	"path"
	"strings"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	c.Data["IsTopic"] = true
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	topics, err := models.GetAllTopics("", "", false)
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["Topics"] = topics
	}
	c.TplNames = "topic.html"
}

func (c *TopicController) Post() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	title := c.Input().Get("title")
	content := c.Input().Get("content")
	tid := c.Input().Get("tid")
	category := c.Input().Get("category")
	label := c.Input().Get("label")

	// 获取附件
	_, fh, err := c.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}
	var attachment string

	// save attachment
	if fh != nil {
		attachment = fh.Filename
		beego.Info(attachment)
		err = c.SaveToFile("attachment", path.Join("attachment", attachment))
		if err != nil {
			beego.Error(err)
		}
	}
	if len(tid) == 0 {
		err = models.AddTopic(title, category, label, content, attachment)
	} else {
		err = models.ModifyTopic(tid, title, category, label, content, attachment)
	}

	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic", 302)
}

func (c *TopicController) Add() {
	c.TplNames = "topic_add.html"
}

func (c *TopicController) View() {
	c.TplNames = "topic_view.html"

	topic, err := models.GetTopic(c.Ctx.Input.Params["0"])
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}

	replies, err := models.GetAllReplies(c.Ctx.Input.Params["0"], false)
	if err != nil {
		beego.Error(err)
		return
	}

	c.Data["Replies"] = replies
	c.Data["Topic"] = topic
	c.Data["Labels"] = strings.Split(topic.Labels, ",")
	c.Data["Tid"] = c.Ctx.Input.Params["0"]
	c.Data["IsLogin"] = checkAccount(c.Ctx)
}

func (c *TopicController) Modify() {
	c.TplNames = "topic_modify.html"
	tid := c.Input().Get("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}

	c.Data["Topic"] = topic
	c.Data["Tid"] = tid
}

func (c *TopicController) Delete() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	err := models.DeleteTopic(c.Ctx.Input.Params["0"])

	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/", 302)

}
