package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lytics/gowrapmx4j"
	"github.com/lytics/gowrapmx4j/cassandra"
)

var (
	mx4j          gowrapmx4j.MX4JService
	host          string
	port          string
	loglvl        string
	gprefix       string
	hostnameid    string
	queryInterval int
)

// Cassandra MX4J status returns all data from the gowrapmx4j.registry
// in its raw form marshalled into JSON.
func cassStatus(w http.ResponseWriter, r *http.Request) {
	metrics, err := gowrapmx4j.QueryMX4J(mx4j)
	if err != nil {
		errString := fmt.Sprintf("cassStatus Error: %#v", err)
		http.Error(w, errString, http.StatusServiceUnavailable)
	}
	log.Debugf("metrics: %#v", metrics)

	mbeans := gowrapmx4j.RegistryBeans() //metrics := mm.ValFunc(mm.MBean)
	js, err := json.Marshal(mbeans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "cassStatus: Error marshaling JSON from MX4J data: %v", err)
	}
	fmt.Fprintf(w, "%s", js)
}

func nodeStatus(w http.ResponseWriter, r *http.Request) {
	mjs := make(map[string]interface{})
	nsb := gowrapmx4j.RegistryGet("NodeStatusBinary")
	metricMap, err := gowrapmx4j.DistillAttributeTypes(nsb.Data)
	if err != nil {
		mjs["ERR"] = "Error extracting node status data"
		mjs["error"] = fmt.Sprintf("%v", err)
	}

	states, ok := metricMap["SimpleStates"]
	log.Debugf("%#v", states)
	if !ok {
		mjs["ERR"] = "Error extracting node status data"
		mjs["error"] = "Key: SimpleStates not in data map"
	}
	ss := states.(map[string]interface{})

	var hostKey string
	for k, v := range ss {
		log.Debug("Nodestatus: %s %#v", k, v)
		hostMatch := regexp.MustCompile(fmt.Sprintf(".*%s.*", hostnameid))
		if hostMatch.MatchString(k) {
			hostKey = k
		}
	}

	mjs[hostKey] = ss[hostKey]
	js, err := json.Marshal(mjs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, fmt.Sprintf("nodeStatus: Error marshaling JSON from MX4J data: %#v", err), 500)
	}
	fmt.Fprintf(w, "%s", js)
}

func main() {
	flag.StringVar(&host, "host", "localhost", "mx4j host address")
	flag.StringVar(&port, "port", "8081", "mx4j port to query")
	flag.StringVar(&loglvl, "loglvl", "info", "Log level to use")
	flag.StringVar(&gprefix, "gprefix", "mx4jcass.", "Graphite prefix tag")
	flag.StringVar(&hostnameid, "hostid", "", "Hostname to use as graphite classifier for metrics")
	flag.IntVar(&queryInterval, "queryInterval", 10, "Interval seconds between querying MX4J")
	flag.Parse()

	//TODO: Sort out metrics writer
	//mlog := golog.New(os.Stderr, "", golog.LstdFlags)
	//go LogToWriter(metrics.DefaultRegistry, time.Second*time.Duration(queryInterval), mlog, gprefix, stop)

	log.Infof("Initializing Cassandra mx4j status endpoint")
	mx4j = gowrapmx4j.MX4JService{Host: host, Port: port}
	mx4j.Init()

	ll, err := log.ParseLevel(loglvl)
	if err != nil {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(ll)
	}

	// TODO: Improve the initialization of the Registry
	// TODO: Multiple definitions of the NewMX4JMetric function is confusing and un-go like

	// Query singlenton values from MX4J
	mm := gowrapmx4j.MX4JMetric{HumanName: "compactions.active", ObjectName: "org.apache.cassandra.internal:type=CompactionExecutor",
		Format: "array", Attribute: "ActiveCount", ValFunc: gowrapmx4j.DistillAttribute}
	gowrapmx4j.RegistrySet(mm, nil)

	mm = gowrapmx4j.MX4JMetric{HumanName: "compactions.pending", ObjectName: "org.apache.cassandra.internal:type=CompactionExecutor",
		Format: "array", Attribute: "PendingTasks", ValFunc: gowrapmx4j.DistillAttribute}
	gowrapmx4j.RegistrySet(mm, nil)

	//OR
	// Query MBean attribute maps
	mm = gowrapmx4j.MX4JMetric{HumanName: "NodeStatusBinary", ObjectName: "org.apache.cassandra.net:type=FailureDetector",
		ValFunc: gowrapmx4j.DistillAttributeTypes}
	gowrapmx4j.RegistrySet(mm, nil)

	mm = gowrapmx4j.MX4JMetric{HumanName: "CompactionExecutor", ObjectName: "org.apache.cassandra.internal:type=CompactionExecutor",
		ValFunc: gowrapmx4j.DistillAttributeTypes}
	gowrapmx4j.RegistrySet(mm, nil)

	// Query Cluster information
	mm = gowrapmx4j.MX4JMetric{HumanName: "StorageService", ObjectName: "org.apache.cassandra.db:type=StorageService",
		ValFunc: gowrapmx4j.DistillAttributeTypes}
	gowrapmx4j.RegistrySet(mm, nil)

	// Simple run loop to query MX4J
	go func() {
		for {
			log.Debug("Querying MX4J")
			_, err := gowrapmx4j.QueryMX4J(mx4j)
			if err != nil {
				log.Errorf("Error Querying MX4J: %v", err)
			}
			time.Sleep(time.Second * time.Duration(queryInterval))
		}
	}()

	http.HandleFunc("/", cassStatus)
	http.HandleFunc("/clean", cassandra.HttpCleanStatus)
	//http.HandleFunc("/status", nodeStatus)
	http.ListenAndServe(":8082", nil)
}
