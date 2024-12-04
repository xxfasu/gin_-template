package repository

import (
	"context"
	"gin_template/internal/repository/gen"
	"gin_template/internal/repository/user_repository"
	"gin_template/pkg/logs"
	"gin_template/pkg/safe"
	"gin_template/pkg/zapgorm2"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var ProviderSet = wire.NewSet(
	NewDB,
	NewTransaction,
	user_repository.NewUserRepository,
)

type transaction struct {
	DB *gorm.DB
	// rdb    *redis.Client
}

func NewTransaction(
	DB *gorm.DB,
) Transaction {
	return &transaction{
		DB: DB,
	}
}

type Transaction interface {
	Transaction(ctx context.Context, fn func(query *gen.Query) error) error
}

func (r *transaction) Transaction(ctx context.Context, fn func(query *gen.Query) error) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var flag error
		var wg sync.WaitGroup
		done := make(chan struct{})
		wg.Add(1)
		safe.Go(func() {
			defer wg.Done()
			query := gen.Use(tx)
			flag = fn(query)
		})
		safe.Go(func() {
			defer close(done)
			wg.Wait()
		})
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
			return flag
		}
	})
}

func NewDB(l *logs.Logger) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	logger := zapgorm2.New(l.Logger)
	dsn := ""

	// GORM doc: https://gorm.io/docs/connecting_to_the_database.html

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})

	if err != nil {
		panic(err)
	}
	db = db.Debug()
	// Connection Pool config
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
