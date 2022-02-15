package main

import (
	"fmt"
	"time"

	"github.com/Obito1903/CY-celcat/internal/cyCelcat"
	"github.com/Obito1903/CY-celcat/pkg/http"

	config "github.com/Obito1903/CY-celcat/pkg"
)

func main() {
	config := config.Configure()
	fmt.Println(config)
	if config.Web {
		go http.StartServer(config)
	}
	if config.Continuous {
		for {
			cyCelcat.Query(config, cyCelcat.Period{Start: time.Now(), End: time.Now().Add(time.Hour * 24 * 7 * 3)})
			time.Sleep(time.Duration(config.QueryDelay) * time.Second)
		}
	} else {
		cyCelcat.Query(config, cyCelcat.Period{Start: time.Now(), End: time.Now().Add(time.Hour * 24 * 7 * 3)})

	}

}
