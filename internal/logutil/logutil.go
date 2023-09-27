package logutil

import (
	"flag"
	"log"
	"os"
	"strconv"
)

const logLevelDefault = 1

var v int

func init() {
	verb := os.Getenv("VERBOSE")
	if verb != "" {
		var err error
		v, err = strconv.Atoi(verb)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		v = logLevelDefault
	}
	flag.IntVar(&v, "v", v, "Уровень «многословности» логов")
}

func V(level int) bool {
	return level <= v
}
