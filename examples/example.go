package main

import (
	"fmt"
	"time"

	"github.com/fieldnation/doctor"
)

func main() {
	doc := doctor.New()

	doc.Schedule("ping test", ping, doctor.Regularity(1*time.Second), doctor.TTL(5*time.Second))

	ch, err := doc.Examine()
	if err != nil {
		panic(err)
	}

	for boh := range ch {
		fmt.Printf("ping started at %s", boh.Start())
	}
}

func ping(b doctor.BillOfHealth) doctor.BillOfHealth {
	fmt.Println("yay")
	b.Healthy()
	return b
}
