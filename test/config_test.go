package config_test

import (
	"testing"

	config "github.com/Obito1903/CY-celcat/pkg"
)

func TestConfig(t *testing.T) {
	config := config.ReadConfig("../config.json")
	t.Log("userName: ", config.UserName, " userPassword: ", config.UserPassword, " celcatHost: ", config.CelcatHost)
	for _, groupe := range config.Groups {
		t.Log("name: ", groupe, " id: ", groupe)
	}
}
