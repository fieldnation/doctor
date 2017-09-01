package doctor

import (
	"fmt"
	"testing"
	"time"
)

func TestDocktor(t *testing.T) {
	doc := New()

	doc.Schedule("ping test", ping, Regularity(1*time.Second))
	ch, err := doc.Examine()
	if err != nil {
		t.Fatal(err)
	}

	for boh := range ch {
		fmt.Printf("ping started at %s", boh.start)
	}

}

func ping(b BillOfHealth) BillOfHealth {
	return b
}
