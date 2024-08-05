package database

import (
	"log"

	"github.com/giovane-aG/video-encoder/encoder/domain"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

type Database struct {
	DB            *gorm.DB
	Dsn           string
	DsnTest       string
	DbType        string
	DbTypeTest    string
	Debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDb() *Database {
	return &Database{}
}

func NewDbTest() *gorm.DB {
	db := NewDb()

	db.Env = "test"
	db.DbTypeTest = "sqlite3"
	db.DsnTest = ":memory:"
	db.AutoMigrateDb = true
	db.Debug = true

	conn, err := db.Connect()

	if err != nil {
		log.Fatalf("Test db error: %v", err)
	}
	return conn
}

func (db *Database) Connect() (*gorm.DB, error) {
	var err error

	if db.Env == "test" {
		db.DB, err = gorm.Open(db.DbTypeTest, db.DsnTest)
	} else {
		db.DB, err = gorm.Open(db.DbType, db.Dsn)
	}

	if err != nil {
		return nil, err
	}

	if db.Debug {
		db.DB.LogMode(true)
	}

	if db.AutoMigrateDb {
		db.DB.AutoMigrate(&domain.Video{}, &domain.Job{})
		db.DB.Model(domain.Job{}).AddForeignKey("video_id", "videos(id)", "cascade", "cascade")
	}

	return db.DB, nil
}
