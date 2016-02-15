Go Wrap MX4J
------------

[![GoDoc](https://godoc.org/github.com/lytics/gowrapmx4?status.svg)](https://godoc.org/github.com/lytics/gowrapmx4j)

Golang wrapper package for accesssing MX4J HTTP data. JMX is useful but horrible, MX4J provides an HTTP endpoint to access JMX data but is shrouded in misleading documentation and odd XML data structures.

This library aims ease interfacing with MX4J to extract useful information about your Java process. MX4J/JMX like to reuse XML tag names which can make coming up with a descriptive Go types a bit difficult and is an ongoing process for the library. 

 ...More to come. Hopefully will provide a useful example usage soon.

### External Requirements
[Logrus](https://github.com/Sirupsen/logrus); for nice log handling.

If using the [Go Vendor Experiment](https://medium.com/@freeformz/go-1-5-s-vendor-experiment-fd3e830f52c3#.fq1ap96hb) everything should just work(yay GO 1.6!). Otherwise you might need to `go get -u github.com/Sirupsen/logrus` to make it available in your GOPATH.

