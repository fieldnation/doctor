package main

import (
	"fmt"
	"time"

	"github.com/fieldnation/doctor"
)

func main() {
	doc := doctor.New()

	doc.Schedule("ping", ping, doctor.Regularity(1*time.Second), doctor.TTL(5*time.Second))
	doc.Schedule("pong", pong, doctor.Regularity(1*time.Second), doctor.TTL(10*time.Second))

	ch := doc.Examine()

	for boh := range ch {
		fmt.Printf("%s started at %s\n", boh.Name(), boh.Start())
	}
}

func ping(b doctor.BillOfHealth) doctor.BillOfHealth {
	fmt.Println("ping...")
	return b
}

func pong(b doctor.BillOfHealth) doctor.BillOfHealth {
	fmt.Println("pong...")
	return b
}
