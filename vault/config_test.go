package vault

import (
	"github.com/tag-service/test"
	"os"
	"strings"
	"testing"
)

func TestLoadConfig_With_Default(t *testing.T) {
	t.Logf("Given I load default vault config")
	{
		os.Setenv(CONFIG, "../config")
		os.Setenv(ENVIRONMENT, "default")
		config := LoadConfig()
		if !config.Enabled {
			t.Logf("\t\t vault login should have been disabled: %v", test.CheckMark)
		} else {
			t.Errorf("\t\t vault login should have been disabled: %v", test.BallotX)
		}

	}
}

// Loading config
func TestLoadConfig_Dev(t *testing.T) {

	t.Logf("Given I load the dev vault config")
	{
		os.Setenv(ENVIRONMENT, "dev")
		config := LoadConfig()
		if config.Enabled {
			t.Logf("\t\t vault login should have been enabled: %v", test.CheckMark)
		} else {
			t.Errorf("\t\t vault login should have been enabled: %v", test.BallotX)
		}

		if config.Address == "https://vault.dev.astoapp.co.uk" {
			t.Logf("\t\t The vault address should have been: %s,  %v", config.Address, test.CheckMark)
		} else {
			t.Errorf("\t\t The vault address should have been: %s,  %v", config.Address, test.BallotX)
		}

	}
}

func TestLoadConfig_Prod(t *testing.T) {
	t.Logf("Given I load the prod vault config")
	{
		os.Setenv(ENVIRONMENT, "prod")
		config := LoadConfig()
		if config.Enabled {
			t.Logf("\t\t vault login should have been enabled: %v", test.CheckMark)
		} else {
			t.Errorf("\t\t vault login should have been enabled: %v", test.BallotX)
		}

		if config.Address == "https://vault.prod.astoapp.co.uk" {
			t.Logf("\t\t The vault address should have been: %s,  %v", config.Address, test.CheckMark)
		} else {
			t.Errorf("\t\t The vault address should have been: %s,  %v", config.Address, test.BallotX)
		}

	}
}

func TestLoadConfig_Dummy_Env(t *testing.T) {
	t.Logf("Given the vault config is not present")
	{
		os.Setenv(CONFIG, "config")
		os.Setenv(ENVIRONMENT, "dummy")
		got := panicValue(func() { LoadConfig() })
		a, ok := got.(error)
		messageContainedInError := "config/vault-config-dummy.yml: no such file or directory"
		if strings.Contains(a.Error(), messageContainedInError) || !ok {
			t.Logf("No such file or dir error is thrown: %s, %v ", a.Error(), test.CheckMark)
		} else {
			t.Errorf("No such file or dir error is thrown:: %s, %v ", a.Error(), test.BallotX)
		}
	}
}

func TestLoadConfig_Invalid_Config_File(t *testing.T) {
	t.Logf("Given the vault config is invalid")
	{
		os.Setenv(ENVIRONMENT, "default")
		os.Setenv(CONFIG, "../test/data")
		got := panicValue(func() { LoadConfig() })
		a, ok := got.(error)
		expectedErrorMessage := "cannot unmarshal !!str `dajkldf...` into vault.Config"
		errorGot := a.Error()
		if strings.Contains(errorGot, expectedErrorMessage) || !ok {
			t.Logf("Unmarshalling error should have been thrown: %s, %v ", a.Error(), test.CheckMark)
		} else {
			t.Errorf("Unmarshalling error should have been thrown: %s, %v ", a.Error(), test.BallotX)
		}
	}
}

func TestGetEnv(t *testing.T) {
	defaultValue := "blah"
	key := "key"
	got := GetEnv(key, defaultValue)
	if got == defaultValue {
		t.Logf("It should get the default value %s, %v ", got, test.CheckMark)
	} else {
		t.Errorf("It should get the default value %s, %v ", got, test.BallotX)
	}
}

// catching a panic
func panicValue(fn func()) (recovered interface{}) {
	defer func() {
		recovered = recover()
	}()
	fn()
	return
}
