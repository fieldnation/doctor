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
	//
	// this example declares the appointment and variadic options
	// within the schedule parameters
	doc.Schedule(doctor.Appointment{
		Name:        "ping",
		HealthCheck: ping,
	}, doctor.Regularity(1*time.Second), doctor.TTL(5*time.Second))

	// schedule an appointment that occurs every second for 20 seconds
	//
	// this example declares the appointment and schedule ahead of time
	pongAppt := doctor.Appointment{
		Name:        "pong",
		HealthCheck: pong,
	}
	pongOpts := []doctor.Options{doctor.Regularity(1 * time.Second), doctor.TTL(20 * time.Second)}
	doc.Schedule(pongAppt, pongOpts...)

	// schedule an appointment that only occurs once
	//
	// this example does not require any variadic options
	doc.Schedule(doctor.Appointment{Name: "only once", HealthCheck: onlyOnce})

	// start the examination and save the recieving channel
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
