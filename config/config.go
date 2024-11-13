package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Configuration interface {
	GetConfiguration() (config *AppConfiguration)
}

type AppConfiguration struct {
	Appname    *string `json:"appname"`
	AppIP      *string `json:"appip"`
	AppPort    *string `json:"appport"`
	Production *string `json:"production"`
	Version    *string `json:"version"`

	// path
	PathTemp       *string `json:"pathtemp"`
	PathLog        *string `json:"pathlog"`
	PathFile       *string `json:"pathfile"`
	PathUpload     *string `json:"pathupload"`
	Certificate    *string `json:"certificate"`
	CertificateKey *string `json:"certificatekey"`

	// misc
	FileSep *string `json:"filesep"`
}

func NewConfiguration(filename string) (config *AppConfiguration, err error) {

	jsonFile, err := os.Open(filename)
	if err != nil {
		err = fmt.Errorf("Unable to load config.json :" + err.Error())
		return config, err
	}
	defer jsonFile.Close()

	byteData, err := io.ReadAll(jsonFile)
	if err != nil {
		err = fmt.Errorf("Unable to read jsonFile :" + err.Error())
		return config, err
	}

	err = json.Unmarshal(byteData, &config)
	if err != nil {
		err = fmt.Errorf("Unable to unmarshall jsonFile -> byteData :" + err.Error())
		return config, err
	}

	return config, err
}

func (c *AppConfiguration) GetConfiguration() (config *AppConfiguration) {
	return c
}
