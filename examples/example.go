package main

import (
	"fmt"
	"time"

	"github.com/fieldnation/doctor"
)

func main() {

	// create new doctor
	doc := doctor.New()

	// schedule an appointment that occurs every second for 5 seconds
	doc.Schedule(doctor.Appointment{
		Name:        "ping",
		HealthCheck: ping,
	}, doctor.Regularity(1*time.Second), doctor.TTL(5*time.Second))

	// schedule an appointment that occurs every second for 20 seconds
	doc.Schedule(doctor.Appointment{
		Name:        "pong",
		HealthCheck: pong,
	}, doctor.Regularity(1*time.Second), doctor.TTL(20*time.Second))

	// schedule an appointment that only occurs once
	doc.Schedule(doctor.Appointment{
		Name:        "only once",
		HealthCheck: onlyOnce,
	})

	// start the examination
	ch := doc.Examine()

	// slurp on the channel to recieve bills of health resulting from each health check
	for boh := range ch {
		// print out info on the bill of health
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
