package main

import "github.com/tj/docopt"
import "encoding/json"
import "strconv"
import "time"
import "log"
import "os"
import "io"

var Version = "0.0.1"

const Usage = `
  Usage:
    throttle <ops> [<duration>]
    throttle -h | --help
    throttle --version

  Options:
    -h, --help        output help information
    -v, --version     output version

`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	if err != nil {
		log.Fatalf("error parsing arguments: %s", err)
	}

	dur := time.Second
	if s, ok := args["<duration>"].(string); ok {
		dur, err = time.ParseDuration(s)
		if err != nil {
			log.Fatalf("error parsing duration: %s", err)
		}
	}

	ops, err := strconv.Atoi(args["<ops>"].(string))
	if err != nil {
		log.Fatalf("error parsing ops: %s", err)
	}

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	throttle := time.Tick(dur / time.Duration(ops))

	for {
		<-throttle

		var v interface{}
		err := dec.Decode(&v)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error decoding: %s", err)
		}

		err = enc.Encode(v)

		if err != nil {
			log.Fatalf("error encoding: %s", err)
		}
	}
}
