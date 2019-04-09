package builder

import "log"

func noerr(name string, err error) {
	if err == nil {
		return
	}
	log.Fatalf("%s: %s", name, err)
}
