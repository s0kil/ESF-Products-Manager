package fault

import "log"

type Fault struct {
	err    error
	reason string
}

func Report(err error, reason string) {
	if err != nil {
		log.Fatal(reason, err)
	}
}
