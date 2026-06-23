package mycasbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"gorm.io/gorm"
	"sync"
	"wallet_chain.com/global"

	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// Initialize the model from a string.
var text = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`

var (
	enforcer *casbin.SyncedEnforcer
	once     sync.Once
)

func Setup(db *gorm.DB, _ string) *casbin.SyncedEnforcer {
	once.Do(func() {
		Apter, err := gormadapter.NewAdapterByDBUseTableName(db, "sys", "casbin_rule")
		if err != nil && err.Error() != "invalid DDL" {
			panic(err)
		}

		m, err := model.NewModelFromString(text)
		if err != nil {
			panic(err)
		}
		enforcer, err = casbin.NewSyncedEnforcer(m, Apter)
		if err != nil {
			panic(err)
		}
		err = enforcer.LoadPolicy()
		if err != nil {
			panic(err)
		}
		// set redis watcher if redis config is not nil

		enforcer.EnableLog(true)
	})

	return enforcer
}

func updateCallback(msg string) {

	global.SHOP_LOG.Info(msg)
	err := enforcer.LoadPolicy()
	if err != nil {

		global.SHOP_LOG.Info(fmt.Sprintf("casbin LoadPolicy err: %v", err))
	}
}
func SetCasbin(key string, enforcer *casbin.SyncedEnforcer) {
	/*e.mux.Lock()
	defer e.mux.Unlock()
	e.casbins[key] = enforcer*/
}
