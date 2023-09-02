package test

import (
	"context"
	"fmt"
	"identity-v2/cmd/config"
	"identity-v2/internal/model"
	"identity-v2/internal/repository/cache"
	"identity-v2/internal/repository/sql"
	"identity-v2/internal/repository/sql/sqlmodel"
	service "identity-v2/internal/service/impl"
	"identity-v2/internal/store"
	"identity-v2/pkg/bcrypthash"
	"identity-v2/pkg/jwt"
	"identity-v2/pkg/redis"
	"testing"
	"time"
)

func TestHello(t *testing.T) {
	if "Hello" == "hello" {
		t.Errorf("it should be case sensitive")
	}
}

func TestLogin_WithoutTOTP(t *testing.T) {
	conf, err := config.Load()
	if err != nil {
		t.Errorf("config load error")
	}

	db, err := store.NewSQLStorage(conf)
	if err != nil {
		fmt.Println(1)
		t.Errorf("database connection error")
	}

	j, err := jwt.NewJwtHandler(conf.RSAPair)
	if err != nil {
		fmt.Println(2)
		t.Errorf("jwt package config error")
	}

	rds := redis.NewRedisClient(conf)

	userRepo := sql.NewUserRepo(db)
	userCache := cache.NewUserCache(userRepo, rds)

	sessionRepo := sql.NewSessionRepo(db)
	sessionCache := cache.NewSessionCache(sessionRepo, rds)

	trackRepo := sql.NewTrackRepo(db)
	loginAttemptRepo := sql.NewLoginAttemptRepo(db)

	loginAttemptSvc := service.NewLoginAttempService(loginAttemptRepo)
	// userSvc := service.NewUserService(userCache, sf, e)
	authSvc := service.NewAuthService(userCache, sessionCache, trackRepo, loginAttemptSvc, j)

	// ------------------------------ test ------------------------------

	tID := int64(123456789)
	tEmail := "Mhmd@Gol.com"
	tPassword := "1234"

	hpass, err := bcrypthash.HashPassword(tPassword)
	if err != nil {
		t.Errorf("password hash error")
	}

	_, err = db.NewInsert().Model(sqlmodel.User{
		ID:             tID,
		UUN:            "A123",
		Username:       "Mhmd",
		HashedPassword: hpass,
		Email:          tEmail,
		Created_at:     time.Now(),
		TOTPIsActive:   false,
		TOTPSecret:     "",
		RoleID:         2,
		StatusID:       1,
	}).Exec(context.Background())
	if err != nil {
		t.Errorf("user insertion error")
	}

	defer func() {
		db.NewDelete().Model((*sqlmodel.User)(nil)).Where("email = ?", tEmail).Exec(context.Background())
	}()

	token, err := authSvc.Login(context.Background(), model.LoginInfo{
		Email:    tEmail,
		Password: tPassword,
		TOTPCode: "",
	})
	if err != nil {
		t.Errorf("auth error")
	}

	tc, err := j.ExtractClaims(token)
	if err != nil {
		t.Errorf("token claims extraction error")
	}

	if tEmail != tc.Email {
		t.Errorf("claims are not correct")
	}

	_, err = sessionCache.ByID(context.Background(), model.ID(tID))
	if err != nil {
		t.Errorf("cached session error")
	}

	_, err = sessionRepo.ByID(context.Background(), model.ID(tID))
	if err != nil {
		t.Errorf("session not in db")
	}

	err = sessionCache.Remove(context.Background(), model.ID(tID))
	if err != nil {
		t.Errorf("session removing from cache error")
	}

	err = sessionRepo.Remove(context.Background(), model.ID(tID))
	if err != nil {
		t.Errorf("session removing error")
	}
}
