package utils

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func RandomString(l int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
func Execute(n string) string {
	out, err := exec.Command("sh", "-c", n).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	return string(out)
}
