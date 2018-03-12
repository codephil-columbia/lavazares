package models

import (
	"testing"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const conn = "user=codephil dbname=lavazaresDB password=password port=5432 host=localhost sslmode=disable"

type DBSuite struct {
	suite.Suite
	DB *gorm.DB
}

var testUser = User{
	Username: "Test",
	Password: "$2a$14$zc5HJi/TedE7EVfkHU1YDuUv0wqhG7XMSCcD2DJPb.SD1eLWkZQqa",
	Email:    "testEmail@test.com",
}

func (suite *DBSuite) SetupSuite() {
	db, err := initTestDB(conn)
	assert.Nil(suite.T(), err)

	suite.DB = db
}

func (suite *DBSuite) TestCanInsertUser() {
	var t = suite.T()
	assert.NotNil(t, suite.DB)
	err := CreateNewUser(&testUser)
	suite.Nil(t, err)
}

func (suite *DBSuite) TestRetrieveUser(withEmail string) {
	var t = suite.T()

	u, err := RetrieveUser(testUser.Email)
	assert.Nil(t, err)
	assert.EqualValues(t, u, testUser)
}

func TestDB(t *testing.T) {
	suite.Run(t, new(DBSuite))
}
