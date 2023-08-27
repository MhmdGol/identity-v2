package main

import (
	"fmt"
	"identity-v2/cmd/config"
	"identity-v2/internal/casbin"
	"identity-v2/internal/controller"
	authapiv1 "identity-v2/internal/proto/authapi/v1"
	userapiv1 "identity-v2/internal/proto/userapi/v1"
	"identity-v2/internal/repository/cache"
	"identity-v2/internal/repository/sql"
	service "identity-v2/internal/service/impl"
	"identity-v2/internal/store"
	"identity-v2/pkg/jwt"
	"identity-v2/pkg/redis"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"google.golang.org/grpc"
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

	j, err := jwt.NewJwtHandler(conf.RSAPair)
	if err != nil {
		return err
	}

	rds := redis.NewRedisClient(conf)

	e, err := casbin.NewEnforcer(conf)
	if err != nil {
		return err
	}

	// --------------------------------------------------------------------------------

	// _, err = db.NewInsert().Model(&sqlmodel.Role{
	// 	Name: "staff",
	// }).Exec(context.Background())
	// if err != nil {
	// 	return err
	// }
	// _, err = db.NewInsert().Model(&sqlmodel.Status{
	// 	Name: "active",
	// }).Exec(context.Background())
	// if err != nil {
	// 	return err
	// }

	// var role sqlmodel.Role
	// err = db.NewSelect().Model(&role).Where("name = ?", "staff").Scan(context.Background())
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(role)

	// var status sqlmodel.Status
	// err = db.NewSelect().Model(&status).Where("name = ?", "active").Scan(context.Background())
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(status)

	// user := &sqlmodel.User{
	// 	ID:             sf.Generate().Int64(),
	// 	UUN:            "A123",
	// 	Username:       "Mhmd",
	// 	HashedPassword: "hashpasshash",
	// 	Email:          "mhmd@gol.com",
	// 	Created_at:     time.Now(),
	// 	TOTPSecret:     "totpKey",
	// 	Role:           role.ID,
	// 	Status:         status.ID,
	// }

	// _, err = db.NewInsert().Model(user).Exec(context.Background())
	// if err != nil {
	// 	return err
	// }

	// userSecret, _ := totp.Generate(totp.GenerateOpts{
	// 	Issuer:      "IdentityServer",
	// 	AccountName: "mhmd",
	// })

	// fmt.Println(userSecret.Secret())

	// isValid := totp.Validate("751994", "XVVKYV5ARYAMEAKE465JXX25AQTKLO73")
	// if !isValid {
	// 	fmt.Println("Failed")
	// }

	// --------------------------------------------------------------------------------

	userRepo := sql.NewUserRepo(db)
	userCache := cache.NewUserCache(userRepo, rds)

	sessionRepo := sql.NewSessionRepo(db)
	sessionCache := cache.NewSessionCache(sessionRepo, rds)

	trackRepo := sql.NewTrackRepo(db)
	loginAttemptRepo := sql.NewLoginAttemptRepo(db)

	loginAttemptSvc := service.NewLoginAttempService(loginAttemptRepo)
	userSvc := service.NewUserService(userCache, sf, e)
	authSvc := service.NewAuthService(userCache, sessionCache, trackRepo, loginAttemptSvc, j)

	if err := InsertAdmin(userSvc, sf, e); err != nil {
		fmt.Println(err)
	}

	userCtrl := controller.NewUserController(userSvc, authSvc, j, e)
	authCtrl := controller.NewAuthController(authSvc, j)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.HttpPort))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	userapiv1.RegisterUserServiceServer(server, userCtrl)
	authapiv1.RegisterAuthServiceServer(server, authCtrl)

	return server.Serve(lis)
}
