package fixture

import (
	"testing"

	"github.com/vincent-tien/bookmark-management/pkg/sqldb"
	"gorm.io/gorm"
)

// Fixture is an interface that defines the methods for setting up the database
// and generating test data.
//
// A fixture is an object that provides a known state for testing. It is used
// to set up the database with specific data for testing purposes.
type Fixture interface {

	// Constraint is a function that returns a string that describes the
	Constraint() string

	// SetupDB returns the database connection that is used by the fixture.
	SetupDB(db *gorm.DB)

	// DB returns the database connection that is used by the fixture.
	//
	// It is used to get the database connection for testing purposes.
	DB() *gorm.DB

	// Migrate performs database schema migration for the fixture's data models and returns an error if migration fails.
	Migrate() error

	// GenerateData generates test data for the fixture and returns an error if generation fails.
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
