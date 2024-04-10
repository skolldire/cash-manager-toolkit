package viper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/app_profile"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/load_properties"
	"github.com/spf13/viper"
)

var (
	path = setPath()
)

const errorLoadingConfiguration = "error loading configuration - get: %w"

type service struct {
	propertyFiles []string
}

var _ load_properties.LoadProperties = (*service)(nil)

func NewService() *service {
	return &service{propertyFiles: []string{"application.yaml", "application-local.yaml",
		"application-pt.yaml", "application-qa.yaml", "application-prod.yaml"}}
}

// Apply load the properties from the files
func (s *service) Apply() (load_properties.Config, error) {
	var pivot map[string]interface{}
	//Read basic config in application.yaml properties file
	v1 := viper.New()
	v1.AddConfigPath(path)
	v1.SetConfigName("application")
	v1.AutomaticEnv()
	err := v1.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf(errorLoadingConfiguration, err))
	}
	//Read custom property files by scope
	v2 := viper.New()
	v2.AddConfigPath(path)
	v2.SetConfigName(s.getPropertyToLoad())
	err = v2.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf(errorLoadingConfiguration, err))
	}
	//Unify properties
	err = v1.MergeConfigMap(v2.AllSettings())
	if err != nil {
		return load_properties.Config{}, err
	}
	err = v1.Unmarshal(&pivot)
	if err != nil {
		return load_properties.Config{}, err
	}
	pivot, err = validateMapProperties(pivot)
	if err != nil {
		return load_properties.Config{}, err
	}

	var prop load_properties.Config
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &prop,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	err = decoder.Decode(pivot)
	if err != nil {
		panic(fmt.Errorf(errorLoadingConfiguration, err))
	}

	return prop, nil
}

// Method that validates that the base configuration files exist.
func (s *service) validateFiles() error {
	fileLis := listFiles()
	for _, s := range s.propertyFiles {
		if !contains(fileLis, s) {
			return errors.New(fmt.Sprintf("The %s file are required.", s))
		}
	}
	return nil
}

// Method that returns the properties to load.
func (s *service) getPropertyToLoad() string {
	scopeProp := fmt.Sprintf("application-%s", app_profile.GetScopeValue())
	profileProp := fmt.Sprintf("application-%s", app_profile.GetProfileByScope())
	fileList := listFiles()
	for _, s := range fileList {
		if strings.Contains(s, scopeProp) {
			return scopeProp
		}
	}
	return profileProp
}

// Validate if the setting is encrypted or an environment variable
func getEnv(key string) string {
	if strings.HasPrefix(key, "${") && strings.HasSuffix(key, "}") {
		return os.Getenv(strings.TrimSuffix(strings.TrimPrefix(key, "${"), "}"))
	}
	return key
}

// Method that walks through the configuration map and chunks the configurations to validate each key.
func validateMapProperties(config map[string]interface{}) (map[string]interface{}, error) {
	for k, v := range config {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			value, err := validateMapProperties(v.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			config[k] = value
		}

		if reflect.TypeOf(v).Kind() == reflect.Slice {
			value, err := validateSliceProperty(v.([]interface{}))
			if err != nil {
				return nil, err
			}
			config[k] = value
		}

		if reflect.TypeOf(v).Kind() == reflect.String {
			config[k] = getEnv(v.(string))
		} else {
			config[k] = v
		}
	}
	return config, nil
}

// Method that walks through the configuration slice and chunks the configurations to validate each key.
func validateSliceProperty(config []interface{}) ([]interface{}, error) {
	for i, e := range config {
		if reflect.TypeOf(e).Kind() == reflect.Map {
			v, err := validateMapProperties(e.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			config[i] = v
		}

		if reflect.TypeOf(e).Kind() == reflect.Slice {
			v, err := validateSliceProperty(e.([]interface{}))
			if err != nil {
				return nil, err
			}
			config[i] = v
		}

		if reflect.TypeOf(e).Kind() == reflect.String {
			config[i] = getEnv(e.(string))
		}

		config[i] = e

	}
	return config, nil
}

// Check if string is in slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// List files in directory
func listFiles() []string {
	var fileNames []string
	files, err := os.ReadDir(path)
	if err != nil {
		panic(fmt.Errorf(errorLoadingConfiguration, err))
	}
	for _, s := range files {
		fileNames = append(fileNames, s.Name())
	}
	return fileNames
}

func setPath() string {
	if path := os.Getenv("CONF_DIR"); path != "" {
		return path
	}

	return "kit/config"
}

func readSecret(secretName string) (string, error) {
	secretPath := "/run/secrets/" + secretName
	secretData, err := ioutil.ReadFile(secretPath)
	if err != nil {
		return "", err
	}

	return string(secretData), nil
}
