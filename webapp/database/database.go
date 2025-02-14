package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseHandler interface {
	InitDatabase() error
	GetProduct(id uint) (Product, error)
	GetProducts() ([]Product, error)
	GetCheckouts(userID uint) ([]Checkout, error)
	CreateCheckout(userID uint, productID uint, productQuantity uint) (uint, error)
	GetCheckout(checkoutID uint) (Checkout, error)
}

const InitDataJSONFileName = "initdata.json"

type Database struct {
	Conn *gorm.DB
}

type Product struct {
	gorm.Model
	Name  string
	Price uint
	Image string
}

type User struct {
	gorm.Model
	Name     string
	Password string
}

type Checkout struct {
	gorm.Model
	UserID          uint
	User            User
	ProductID       uint
	Product         Product
	ProductQuantity uint
}

type Blob struct {
	Products []Product `json:"products"`
	Users    []User    `json:"users"`
}

func NewDatabaseHandler(environment string, conn *gorm.DB) (DatabaseHandler, error) {
	switch environment {
	case "production":
		return NewProdDatabaseHandler(conn), nil
	case "development":
		return NewDevDatabaseHandler(conn), nil
	default:
		return nil, fmt.Errorf("environment: %s is not supported", environment)
	}
}

func InitializeProdDBConn() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		getEnvDBHostname(),
		getEnvDBPort(),
		getEnvDBUsername(),
		getEnvDBName(),
		getEnvDBPassword(),
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func InitializeDevDBConn() (*gorm.DB, error) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	return gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
}

func getEnvDBHostname() string {
	return getEnv("DB_HOSTNAME", "scstore-database")
}

func getEnvDBPort() int {
	val := getEnv("DB_PORT", "5432")
	port, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("DB_PORT should be integer, but %s: %v", val, err)
	}

	return port
}

func getEnvDBUsername() string {
	return getEnv("DB_USERNAME", "scstore")
}

func getEnvDBPassword() string {
	return getEnv("DB_PASSWORD", "scstore")
}

func getEnvDBName() string {
	return getEnv("DB_NAME", "scstore")
}

func getEnv(key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}
