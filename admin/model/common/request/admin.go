package request

type Admin struct {
	Id            uint   `json:"id" form:"id" `
	Account       string `json:"account" gorm:"account" comment:"用户账号"`
	Mail          string `json:"mail" comment:"用户邮箱"`
	Name          string `comment:"用户昵称"`
	Mobile        uint64 `json:"mobile" comment:"用户手机号码"`
	Password      string `json:"-" comment:"用户密码"`
	GroupId       string `comment:"用户所属群组" json:"group_id"  gorm:"group_id"`
	Status        uint8  `comment:"用户状态"`
	LoginVisit    uint   `comment:"登录次数"`
	LastLoginIp   string `comment:"最后登录IP"`
	LastLoginedAt uint   `comment:"最后登录时间"`
	Groupname     string `json:"groupname" comment:"群组名" `
	RoleIds       string `json:"-" `
}
