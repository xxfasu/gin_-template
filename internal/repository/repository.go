package repository

import (
	"context"
	"gin_template/pkg/logs"
	"gin_template/pkg/zapgorm2"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"time"
)

var ProviderSet = wire.NewSet(
	NewDB,
	NewRepository,
	NewTransaction,
	NewUserRepository,
)

const ctxTxKey = "TxKey"

type Repository struct {
	db *gorm.DB
	//rdb    *redis.Client
	logger *logs.Logger
}

func NewRepository(
	logger *logs.Logger,
	db *gorm.DB,
	// rdb *redis.Client,
) *Repository {
	return &Repository{
		db: db,
		//rdb:    rdb,
		logger: logger,
	}
}

type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTransaction(r *Repository) Transaction {
	return r
}

// DB return tx
// If you need to create a Transaction, you must call DB(ctx) and Transaction(ctx,fn)
func (r *Repository) DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ctxTxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return r.db.WithContext(ctx)
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, ctxTxKey, tx)
		var err error
		go func() {
			err = fn(ctx)
		}()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second * 10):
			return err
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
