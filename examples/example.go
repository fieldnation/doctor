package main

import (
	"fmt"
	"time"

	"github.com/fieldnation/doctor"
)

func main() {
	doc := doctor.New()

	doc.Schedule(doctor.Appointment{
		Name:        "ping",
		HealthCheck: ping,
	}, doctor.Regularity(1*time.Second), doctor.TTL(5*time.Second))

	doc.Schedule(doctor.Appointment{
		Name:        "pong",
		HealthCheck: pong,
	}, doctor.Regularity(1*time.Second), doctor.TTL(20*time.Second))

	doc.Schedule(doctor.Appointment{
		Name:        "only once",
		HealthCheck: onlyOnce,
	})

	ch := doc.Examine()

	for boh := range ch {
		fmt.Printf("%s started at %s\n", boh.Name(), boh.Start())
	}
}

func ping(b doctor.BillOfHealth) doctor.BillOfHealth {
	return b
}

func pong(b doctor.BillOfHealth) doctor.BillOfHealth {
	return b
}

func onlyOnce(b doctor.BillOfHealth) doctor.BillOfHealth {
	return b
}
