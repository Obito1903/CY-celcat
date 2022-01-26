package config_test

import (
	"testing"

	config "github.com/Obito1903/CY-celcat/pkg"
)

func TestConfig(t *testing.T) {
	config := config.ReadConfig("../config.json")
	t.Log("userName: ", config.UserName, " userPassword: ", config.UserPassword, " celcatHost: ", config.CelcatHost)
	for _, groupe := range config.Groupes {
		t.Log("name: ", groupe.Name, " id: ", groupe.Id)
	}
}
