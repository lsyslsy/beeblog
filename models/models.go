package models

import (
	"errors"
	"fmt"
	//"github.com/Unknwon/com"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	_DB_NAME        = "data/beeblog.db"
	_SQLITE3_DRIVER = "sqlite3"
	_MYSQL_DRIVER   = "mysql"
)

type Category struct {
	Id              int64
	Title           string    `orm:"null"`
	Created         time.Time `orm:"null;index"`
	Views           int64     `orm:"null;index"` // viewed times
	TopicTime       time.Time `orm:"null;index"`
	TopicCount      int64     `orm:"null"`
	TopicLastUserId int64     `orm:"null"`
}

type Topic struct {
	Id              int64
	Uid             int64  `orm:"null"`
	Title           string `orm:"null"`
	Labels          string `orm:"null"`
	Category        string
	Content         string    `orm:"size(5000)"`
	Attachment      string    `orm:"null"`
	Created         time.Time `orm:"null;index"`
	Updated         time.Time `orm:"null;index"`
	Views           int64     `orm:"null;index"` // viewed times
	Author          string    `orm:"null"`
	ReplyTime       time.Time `orm:"null;index"`
	ReplyCount      int64     `orm:"null"`
	ReplyLastUserId int64     `orm:"null"`
}

// comment
type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

func RegisterDB() {
	// if !com.IsExist(_DB_NAME) {
	// 	os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
	// 	os.Create(_DB_NAME)
	// }

	orm.RegisterModel(new(Category), new(Topic), new(Comment))
	orm.RegisterDriver(_MYSQL_DRIVER, orm.DR_MySQL)
	orm.RegisterDataBase("default", _MYSQL_DRIVER, "beego:bee@/beeblog?charset=utf8", 10)
}

func AddCategory(name string) error {
	o := orm.NewOrm()

	cate := &Category{Title: name}

	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return errors.New("category exists")
	}

	_, err = o.Insert(cate)
	if err != nil {
		return err
	}
	return nil

}

func DelCategory(cid string) error {
	cidNum, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id: cidNum}
	_, err = o.Delete(cate)
	return err
}

func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()

	cates := make([]*Category, 0)

	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

func GetCategory(cid string) (*Category, error) {
	o := orm.NewOrm()
	cidNum, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return nil, err
	}
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("id", cidNum).One(cate)

	return cate, err
}

func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	reply := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
	}
	o := orm.NewOrm()
	if _, err = o.Insert(reply); err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyCount++
		topic.ReplyTime = time.Now()
		_, err = o.Update(topic)
	}
	return err
}

func AddTopic(title, category, label, content, attachment string) error {
	label = "," + strings.Replace(label, " ", "", -1) + ","
	o := orm.NewOrm()
	topic := &Topic{
		Title:      title,
		Content:    content,
		Labels:     label,
		Attachment: attachment,
		Category:   category,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	_, err := o.Insert(topic)
	if err != nil {
		return err
	}

	// 更新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		cate.TopicCount++
		if _, err = o.Update(cate); err != nil {
			return err
		}
	}
	return err

}

func GetAllTopics(cate, label string, isDesc bool) ([]*Topic, error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)

	qs := o.QueryTable("topic")
	if len(cate) > 0 {
		qs = qs.Filter("category", cate)
	}

	if len(label) > 0 {
		qs = qs.Filter("labels__contains", ","+label+",")
	}
	var err error
	if isDesc {
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		_, err = qs.All(&topics)
	}
	return topics, err
}

// GetAllReplies
func GetAllReplies(tid string, isDesc bool) ([]*Comment, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	replies := make([]*Comment, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	if isDesc {
		_, err = qs.Filter("tid", tidNum).OrderBy("-created").All(&replies)
	} else {
		_, err = qs.Filter("tid", tidNum).All(&replies)
	}
	return replies, err
}

func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()

	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}
	topic.Views++
	_, err = o.Update(topic)
	topic.Labels = topic.Labels[1 : len(topic.Labels)-1]
	return topic, err
}

func ModifyTopic(tid, title, category, label, content, attachment string) error {
	label = "," + strings.Replace(label, " ", "", -1) + ","
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	var oldCate, oldAttach string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		oldAttach = topic.Attachment
		topic.Title = title
		topic.Labels = label
		topic.Content = content
		topic.Category = category
		topic.Attachment = attachment
		topic.Updated = time.Now()
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}

	// 更新old分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err := qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			if _, err = o.Update(cate); err != nil {
				return err
			}
		}

	}
	// delete old attachment
	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment", oldAttach))
	}
	// 更新新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		cate.TopicCount++
		if _, err = o.Update(cate); err != nil {
			return err
		}
	}

	return nil

}

func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		if _, err = o.Delete(topic); err != nil {
			return err
		}
		cate := new(Category)
		qs := o.QueryTable("category")
		if err = qs.Filter("title", topic.Category).One(cate); err != nil {
			return err
		}
		_, err = o.Delete(cate)

	}
	return err
}

func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}
	var tidNum int64
	o := orm.NewOrm()
	reply := &Comment{Id: ridNum}
	if o.Read(reply) == nil {
		tidNum = reply.Tid
		if _, err = o.Delete(reply); err != nil {
			return err
		}
	}

	replies := make([]*Comment, 0)
	replies, err = GetAllReplies(fmt.Sprint(ridNum), true)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyCount--
		if len(replies) != 0 {
			topic.ReplyTime = replies[0].Created
		}
		_, err = o.Update(topic)
	}

	return err
}
