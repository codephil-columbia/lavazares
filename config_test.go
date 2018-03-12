package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	yaml "gopkg.in/yaml.v2"
)

const (
	filePath = "./env.yml"
)

var requiredFields = []string{"user", "dbname", "password", "port", "host", "sslmode"}

type ConfigSuite struct {
	suite.Suite
	FilePath string
}

func (suite *ConfigSuite) SetupTest() {
	suite.FilePath = filePath
}

func (suite *ConfigSuite) TestConfigFields() {
	var t = suite.T()

	file, err := os.Open(filePath)
	assert.Nil(t, err)

	data, err := ioutil.ReadAll(file)
	assert.Nil(t, err)

	configMap := make(map[string]interface{})
	err = yaml.Unmarshal(data, &configMap)
	assert.Nil(t, err)

	foundFields := make([]string, 0, len(configMap["db"].(map[interface{}]interface{})))

	for key := range configMap["db"].(map[interface{}]interface{}) {
		foundFields = append(foundFields, key.(string))
	}

	assert.EqualValues(t, requiredFields, foundFields)
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
