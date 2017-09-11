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

	// schedule an appointment that occurs every 5 seconds and runs forever with no TTL
	doc.Schedule(doctor.Appointment{
		Name:        "forever",
		HealthCheck: forever,
	}, doctor.Regularity(5*time.Second))

	// schedule an appointment that only occurs once
	//
	// this example does not require any variadic options
	doc.Schedule(doctor.Appointment{Name: "only once", HealthCheck: onlyOnce})

	// start the examination and save the recieving channel
	ch := doc.Examine()

	go func() {
		time.Sleep(5 * time.Second)
		doc.Close()
	}()

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

func forever(b doctor.BillOfHealth) doctor.BillOfHealth {
	return b
}

func halfsecond(b doctor.BillOfHealth) doctor.BillOfHealth {
	return b
}
