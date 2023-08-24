package casbin

import (
	"fmt"
	"identity-v2/cmd/config"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v2"
)

func NewEnforcer(conf config.Config) (*casbin.Enforcer, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
	)
	a, _ := xormadapter.NewAdapter("mssql", dsn, true)

	m, _ := model.NewModelFromString(`
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[role_definition]
		g = _, _

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
	`)

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	InitializeEnforcer(e)

	return e, nil
}

func InitializeEnforcer(e *casbin.Enforcer) {
	e.LoadPolicy()

	e.AddPolicy("admin", "users", "create")
	e.AddPolicy("admin", "users", "update")
	e.AddPolicy("admin", "users", "read")
	e.AddPolicy("admin", "permissions", "create")
	e.AddPolicy("admin", "permissions", "update")
	e.AddPolicy("admin", "permissions", "read")
	e.AddPolicy("admin", "status", "create")
	e.AddPolicy("admin", "status", "update")
	e.AddPolicy("admin", "status", "read")

	e.AddPolicy("staff", "users", "update")
	e.AddPolicy("staff", "users", "read")
	e.AddPolicy("staff", "permissions", "update")
	e.AddPolicy("staff", "permissions", "read")
	e.AddPolicy("staff", "status", "update")
	e.AddPolicy("staff", "status", "read")

	e.AddPolicy("user", "users", "read")

	e.SavePolicy()
}
