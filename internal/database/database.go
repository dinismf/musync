package database

import (
	"context"
	"log"
	"time"

	"github.com/dinis/musync/internal/config"
	"github.com/dinis/musync/internal/errors"
	"github.com/dinis/musync/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a wrapper around gorm.DB that provides context support
type DB struct {
	*gorm.DB
	ctx context.Context
}

var GlobalDB *DB

// GetDB returns the global DB instance
func GetDB() *DB {
	return GlobalDB
}

// InitDB initializes the database connection and returns a DB instance
func InitDB(cfg config.DatabaseConfig) *DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database connection:", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Lifetime) * time.Second)

	// Auto migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.SocialLink{},
		&models.Artist{},
		&models.Label{},
		&models.Release{},
		&models.FeedItem{},
		&models.MusicLibrary{},
		&models.Track{},
		&models.Tempo{},
		&models.Playlist{},
		&models.PlaylistTrack{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	dbWrapper := &DB{DB: db, ctx: context.Background()}
	GlobalDB = dbWrapper // Set global DB instance for backward compatibility
	log.Println("Database connection established and migrations completed")
	return dbWrapper
}

// WithContext returns a new DB instance with the given context
func (db *DB) WithContext(ctx context.Context) *DB {
	return &DB{DB: db.DB.WithContext(ctx), ctx: ctx}
}

// Create creates a new record with context support
func (db *DB) Create(ctx context.Context, value interface{}) error {
	result := db.WithContext(ctx).DB.Create(value)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to create record")
	}
	return nil
}

// CreateWithoutContext creates a new record using the stored context
func (db *DB) CreateWithoutContext(value interface{}) error {
	return db.Create(db.ctx, value)
}

// First finds the first record that matches the given conditions with context support
func (db *DB) First(ctx context.Context, dest interface{}, conds ...interface{}) error {
	result := db.WithContext(ctx).DB.First(dest, conds...)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("record not found", nil)
		}
		return errors.Wrap(result.Error, "failed to find record")
	}
	return nil
}

// FirstWithoutContext finds the first record that matches the given conditions using the stored context
func (db *DB) FirstWithoutContext(dest interface{}, conds ...interface{}) error {
	return db.First(db.ctx, dest, conds...)
}

// Find finds all records that match the given conditions with context support
func (db *DB) Find(ctx context.Context, dest interface{}, conds ...interface{}) error {
	result := db.WithContext(ctx).DB.Find(dest, conds...)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to find records")
	}
	return nil
}

// FindWithoutContext finds all records that match the given conditions using the stored context
func (db *DB) FindWithoutContext(dest interface{}, conds ...interface{}) error {
	return db.Find(db.ctx, dest, conds...)
}

// Save saves a record with context support
func (db *DB) Save(ctx context.Context, value interface{}) error {
	result := db.WithContext(ctx).DB.Save(value)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to save record")
	}
	return nil
}

// SaveWithoutContext saves a record using the stored context
func (db *DB) SaveWithoutContext(value interface{}) error {
	return db.Save(db.ctx, value)
}

// Delete deletes a record with context support
func (db *DB) Delete(ctx context.Context, value interface{}, conds ...interface{}) error {
	result := db.WithContext(ctx).DB.Delete(value, conds...)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete record")
	}
	return nil
}

// DeleteWithoutContext deletes a record using the stored context
func (db *DB) DeleteWithoutContext(value interface{}, conds ...interface{}) error {
	return db.Delete(db.ctx, value, conds...)
}

// Where adds a where condition to the query with context support
func (db *DB) Where(ctx context.Context, query interface{}, args ...interface{}) *DB {
	return &DB{DB: db.WithContext(ctx).DB.Where(query, args...), ctx: ctx}
}

// WhereWithoutContext adds a where condition to the query using the stored context
func (db *DB) WhereWithoutContext(query interface{}, args ...interface{}) *DB {
	return db.Where(db.ctx, query, args...)
}

// Transaction starts a transaction with context support
func (db *DB) Transaction(ctx context.Context, fn func(tx *DB) error) error {
	return db.WithContext(ctx).DB.Transaction(func(tx *gorm.DB) error {
		return fn(&DB{DB: tx})
	})
}
