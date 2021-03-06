package model

import (
	"51job/cons"
	// "log"
	"51job/log"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	_ "io"
	"os"
	"time"
)

func init() {
	// 设置为 UTC 时间
	orm.DefaultTimeLoc = time.Local
	//打出查询语句
	orm.Debug = cons.OpenDbLog
	w, _ := os.OpenFile(cons.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	orm.DebugLog = orm.NewLog(w)

	maxIdle := 30
	maxConn := 30
	orm.RegisterDataBase("default", "mysql", cons.Db, maxIdle, maxConn)
	// 需要在init中注册定义的model
	orm.RegisterModelWithPrefix("51job_", new(User), new(UserKeyword), new(Userinfo), new(Keyword))
	// create table
	// orm.RunSyncdb("default", false, true)
}

func SaveUser(today string, u *User, data []byte, keepfile string, keyword string, address string, kind string) error {
	t, _ := time.ParseInLocation("20060102", today, time.Local)
	// log.Fatal("%s", t)
	o := orm.NewOrm()
	err := o.Begin()
	uk := new(UserKeyword)
	uk.Id51 = u.Id51
	uk.FileAddress = keepfile
	uk.Date51 = u.Date51

	ui := new(Userinfo)
	ui.Id51 = u.Id51
	ui.Date51 = u.Date51
	ui.Content = string(data)

	k := new(Keyword)

	/**/
	k.Created = t
	k.Keyword = keyword
	k.Address = address
	k.Kind = kind
	ke := o.Read(k, "Keyword", "Address", "Kind")
	if ke == orm.ErrNoRows {
		log.Println("数据库不存在该关键字")
		k.Updated = t
		_, err1 := o.Insert(k)

		if err1 != nil {
			o.Rollback()
			return err1
		}

	} else if ke == orm.ErrMissPK {
		log.Println("找不到主键")
	} else {
		log.Println("数据库存在该关键字记录")
		k.Updated = t
		k.Time51 = k.Time51 + 1
		if _, err := o.Update(k); err != nil {
			o.Rollback()
			return err
		}
	}
	/**/
	r := o.Read(u, "Id51")
	if r == orm.ErrNoRows {
		log.Println("数据库不存在该简历记录")

		_, err1 := o.Insert(u)

		if err1 != nil {
			o.Rollback()
			return err1
		}

	} else if r == orm.ErrMissPK {
		log.Println("找不到主键")
	} else {
		log.Println("数据库存在该简历记录")
		u.Date51 = uk.Date51
		if _, err := o.Update(u); err != nil {
			o.Rollback()
			return err
		}
	}
	/**/
	uke := o.Read(uk, "FileAddress")
	if uke == orm.ErrNoRows {
		log.Println("数据库不存在该简历-关键字记录")
		uk.User = u
		uk.Keyword = k
		_, err2 := o.Insert(uk)
		if err2 != nil {
			o.Rollback()
			return err2
		}
	} else if uke == orm.ErrMissPK {
		log.Println("找不到主键")
	} else {
		log.Println("数据库存在该该简历-关键字记录")
		uk.Date51 = u.Date51
		if _, err := o.Update(uk); err != nil {
			o.Rollback()
			return err
		}
	}

	/**/
	rui := o.Read(ui, "Id51")

	if rui == orm.ErrNoRows {
		log.Println("数据库不存在该简历详情页")
		ui.User = u
		_, err3 := o.Insert(ui)
		if err3 != nil {
			o.Rollback()
			return err3
		}

	} else if r == orm.ErrMissPK {
		log.Println("找不到主键")
	} else {
		log.Println("数据库存在该简历详情页")
		ui.Content = string(data)
		ui.Date51 = u.Date51
		ui.User = u
		if _, err := o.Update(ui); err != nil {
			o.Rollback()
			return err
		}
	}

	err = o.Commit()
	return err
}
