package main

import (
	"context"
	"fmt"
	"identity-v2/cmd/config"
	"identity-v2/internal/repository/sql/sqlmodel"
	"identity-v2/internal/store"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run() error {
	conf, err := config.Load()
	if err != nil {
		return err
	}

	db, err := store.NewSQLStorage(conf)
	if err != nil {
		return err
	}

	node, _ := strconv.Atoi(conf.SnowflakeNode)
	if err != nil {
		return err
	}

	sf, err := snowflake.NewNode(int64(node))
	if err != nil {
		return err
	}

	// --------------------------------------------------------------------------------

	_, err = db.NewInsert().Model(&sqlmodel.Role{
		Name: "staff",
	}).Exec(context.Background())
	if err != nil {
		return err
	}
	_, err = db.NewInsert().Model(&sqlmodel.Status{
		Name: "active",
	}).Exec(context.Background())
	if err != nil {
		return err
	}

	var role sqlmodel.Role
	err = db.NewSelect().Model(&role).Where("name = ?", "staff").Scan(context.Background())
	if err != nil {
		return err
	}

	fmt.Println(role)

	var status sqlmodel.Status
	err = db.NewSelect().Model(&status).Where("name = ?", "active").Scan(context.Background())
	if err != nil {
		return err
	}

	fmt.Println(status)

	user := &sqlmodel.User{
		ID:             sf.Generate().Int64(),
		UUN:            "A123",
		Username:       "Mhmd",
		HashedPassword: "hashpasshash",
		Email:          "mhmd@gol.com",
		Created_at:     time.Now(),
		TOTPSecret:     "totpKey",
		Role:           role.ID,
		Status:         status.ID,
	}

	_, err = db.NewInsert().Model(user).Exec(context.Background())
	if err != nil {
		return err
	}

	// --------------------------------------------------------------------------------

	return nil
}