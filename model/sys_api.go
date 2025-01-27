package model

import "time"

// Api 接口
type Api struct {
	Id     int       `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`
	Name   string    `xorm:"VARCHAR(255) comment('名称')" json:"name"`
	Path   string    `xorm:"VARCHAR(255) comment('路径')" json:"path"`
	Group  string    `xorm:"VARCHAR(255) comment('组')" json:"group"`
	Method string    `xorm:"VARCHAR(255) comment('方法')" json:"method"`
	Ctime  time.Time `xorm:"DATETIME comment('时间')" json:"ctime"`
}

func (Api) TableName() string {
	return "sys_api"
}

// ApiGet 单条接口信息
func ApiGet(id int) (*Api, bool) {
	mod := &Api{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// ApiAll 所有接口信息
func ApiAll() ([]Api, error) {
	mods := make([]Api, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// ApiPage 接口分页信息
func ApiPage(pi int, ps int, cols ...string) ([]Api, error) {
	mods := make([]Api, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// ApiCount 接口分页信息总数
func ApiCount() int {
	mod := &Api{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// ApiAdd 添加接口信息
func ApiAdd(mod *Api) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// ApiEdit 编辑接口信息
func ApiEdit(mod *Api, cols ...string) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// ApiMapIds 通过id集合返回接口信息
func ApiMapIds(ids []int) map[int]*Api {
	mods := make([]Api, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	mapSet := make(map[int]*Api, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}

// ApiDrop 删除单条接口信息
func ApiDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&Api{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// ------------------------------------------------------ 角色接口 ------------------------------------------------------

// RoleApiAll  通过RoleId查询所有接口信息
func RoleApiAll(roleId int) ([]Api, error) {
	mods := make([]Api, 0, 8)
	err := Db.SQL("SELECT sys_api.* FROM sys_api LEFT JOIN sys_role_api ON sys_api.id = sys_role_api.api_id WHERE sys_role_api.role_id = ?", roleId).Find(&mods)
	return mods, err
}
