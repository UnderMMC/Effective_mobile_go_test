package main

import (
	"EffectiveMobile_Go/internal/app"
	_ "database/sql"
	_ "github.com/lib/pq"
)

func main() {
	app.New().Run()
}
