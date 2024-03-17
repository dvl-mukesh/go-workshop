package dvlutil

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ReadEnvVars reads environment variables based on tags specified in the "env" tag
// of the fields in the provided configuration struct (cfg). It populates the struct
// fields with the corresponding environment variable values.
//
// The "env" tag should contain the name of the environment variable. An optional
// "required" tag can be added to indicate that the variable must be present; otherwise,
// an error is returned.
//
// Supported field types: string, bool and int. If a field is of type bool or int, the environment
// variable is expected to be a string representing a boolean value or an int value.
func ReadEnvVars(cfg interface{}) error {
	cfgType := reflect.TypeOf(cfg).Elem()
	cfgValue := reflect.ValueOf(cfg).Elem()

	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		tags := strings.Split(field.Tag.Get("env"), ",")

		if len(tags) > 0 && tags[0] != "" {
			envName := tags[0]
			envValue, exists := os.LookupEnv(envName)
			required := isRequired(tags)

			if !exists && required {
				return fmt.Errorf("required environment variable %s not found", tags[0])
			}

			err := setFieldValueFromEnv(cfgValue.Field(i), envName, envValue)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isRequired(tags []string) bool {
	required := false

	if len(tags) > 1 && tags[1] == "required" {
		required = true
	}
	return required
}

func setFieldValueFromEnv(field reflect.Value, envName, envValue string) error {
	if envValue == "" {
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(envValue)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(envValue)
		if err != nil {
			return fmt.Errorf("value: %s, in env var %s is not of type bool", envValue, envName)
		}
		field.SetBool(boolVal)
	case reflect.Int:
		intVal, err := strconv.ParseInt(envValue, 10, 64)
		if err != nil {
			return fmt.Errorf("value: %s, in env variable %s is not of type int", envValue, envName)
		}
		field.SetInt(intVal)
	}
	return nil
}
