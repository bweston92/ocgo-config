# OpenCensus Configuration for Golang

Parses URI type value for OpenCensus exporters.

If no value is passed to `Setup` then it will check `OCGO_EXPORTER`.

# Example

```golang
package main

import (
	"flag"
	"time"

	"github.com/bweston92/ocgoconfig"
	"go.opencensus.io/trace"
)

var (
	ocgov string
)

func init() {
	flag.StringVar(&ocgov, "oc-exporter", "", "Configuration for OpenCensus exporter")
}

func main() {
	flag.Parse()

	ocgoconfig.Setup(ocgov)

	s := trace.NewSpan("myheavywork")
	for i := 0; i < 100 {
		time.Sleep(time.Microsecond)
	}
	s.End()
}

```


