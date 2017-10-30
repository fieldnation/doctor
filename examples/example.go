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
		Name: "ping",
		HealthCheck: func(b doctor.BillOfHealth) doctor.BillOfHealth {
			return b
		},
	}, doctor.Regularity(1*time.Second), doctor.TTL(5*time.Second))

	// schedule an appointment that occurs every second for 20 seconds
	doc.Schedule(doctor.Appointment{
		Name: "pong",
		HealthCheck: func(b doctor.BillOfHealth) doctor.BillOfHealth {
			return b
		},
	}, doctor.Regularity(1*time.Second), doctor.TTL(20*time.Second))

	// schedule an appointment that occurs every 5 seconds and runs forever with no TTL
	doc.Schedule(doctor.Appointment{
		Name: "forever",
		HealthCheck: func(b doctor.BillOfHealth) doctor.BillOfHealth {
			return b
		},
	}, doctor.Regularity(100*time.Millisecond))

	// schedule an appointment that only occurs once
	// this example does not require any variadic options
	doc.Schedule(doctor.Appointment{
		Name: "only once",
		HealthCheck: func(b doctor.BillOfHealth) doctor.BillOfHealth {
			return b
		},
	})

	// schedule a bad appointment that will fail
	doc.Schedule(doctor.Appointment{
		Name: "unhealth",
		HealthCheck: func(b doctor.BillOfHealth) doctor.BillOfHealth {
			b.SetHealth(!b.Healthy())
			return b
		},
	})

	// start the examination and record the recieving channel
	ch := doc.Examine()

	// lets run doctor for 5 seconds and then gracefully stop execution
	go func() {
		time.Sleep(5 * time.Second)
		doc.Close()
	}()

	// slurp on the channel to recieve bills of health
	for boh := range ch {
		// print out info on the bill of health
		fmt.Printf("%s started at %s\n", boh.Name(), boh.Start())
	}
}
