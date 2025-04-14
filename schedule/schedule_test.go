package schedule

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewCronScheduler(t *testing.T) {
	sc, err := NewCronScheduler(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	message := "waiting"

	done := func() {
		message = "done"
	}

	id, err := sc.AddFunc("@minute", done)
	if err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Second * 61)

	fmt.Println(id, message)
}
