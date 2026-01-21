package fixture

import (
	"testing"

	"github.com/vincent-tien/bookmark-management/pkg/sqldb"
	"gorm.io/gorm"
)

type Fixture interface {
	SetupDB(db *gorm.DB)
	DB() *gorm.DB
	Migrate() error
	GenerateData() error
}

func NewFixture(t *testing.T, f Fixture) *gorm.DB {
	// create test db
	f.SetupDB(sqldb.InitMockDb(t))

	// Migrate schema
	if err := f.Migrate(); err != nil {
		t.Fatal("Failed to migrate schema:", err)
	}

	// Generate test data
	err := f.GenerateData()
	if err != nil {
		t.Fatal("Failed to generate test data:", err)
	}

	return f.DB()
}
