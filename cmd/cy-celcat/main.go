package main

import (
	"fmt"
	"time"

	"github.com/Obito1903/CY-celcat/internal/cyCelcat"

	config "github.com/Obito1903/CY-celcat/pkg"
)

func main() {
	config := config.Configure()
	fmt.Println(config)
	cyCelcat.Query(config, cyCelcat.Period{Start: time.Now(), End: time.Now().Add(time.Hour * 24 * 7 * 3)})
}
