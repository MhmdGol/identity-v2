package main

import (
	"context"
	"identity-v2/internal/model"
	"identity-v2/internal/service"

	"github.com/bwmarrin/snowflake"
	"github.com/casbin/casbin/v2"
)

func InsertAdmin(userSvc service.UserService, sf *snowflake.Node, e *casbin.Enforcer) error {
	b, err := userSvc.Exists(context.Background(), "su@gmail.com")
	if err != nil {
		return err
	}
	if !b {
		e.LoadPolicy()
		e.AddGroupingPolicy("su@gmail.com", "admin")
		e.SavePolicy()

		return userSvc.Create(context.Background(), model.RawUser{
			ID:       model.ID(sf.Generate().Int64()),
			Username: "su",
			Password: "Admin@123",
			Email:    "su@gmail.com",
			Role:     "admin",
			Status:   "active",
		})
	}

	return nil
}
