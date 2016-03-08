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
	metrics "github.com/rcrowley/go-metrics"
)

var (
	mx4j          gowrapmx4j.MX4J
	host          string
	port          string
	loglvl        string
	gprefix       string
	hostnameid    string
	queryInterval int
)

// Query all registered MX4J endpoints and compose their data into the MX4JMetric
// array or return error
func QueryMX4J(mx4j gowrapmx4j.MX4J) (*[]gowrapmx4j.MX4JMetric, error) {
	reg := gowrapmx4j.RegistryGetAll()

	for _, mm := range reg {
		var newData gowrapmx4j.MX4JData
		var err error
		data := mm.Data
		log.Debugf("Metric being queried: %#v", mm)

		// If first time querying endpoint, create data struct
		if data == nil {
			newData = gowrapmx4j.Bean{}
			mx4jData, err := newData.QueryMX4J(mx4j, mm)
			if err != nil {
				retErr := fmt.Errorf("QueryMX4J Error: %v%s", newData, err)
				return nil, retErr
			}
			gowrapmx4j.RegistrySet(mm, mx4jData)
		} else {
			newData, err = data.QueryMX4J(mx4j, mm)
			gowrapmx4j.RegistrySet(mm, newData)
		}

		if mm.MetricFunc != nil && newData != nil {
			log.Debugf("Metric func running: %s", mm.HumanName)
			mm.MetricFunc(mm.Data, mm.HumanName)
		}

		if newData == nil {
			log.Errorf("No data returned from querying; blanking the metric registries")
			metrics.DefaultRegistry.UnregisterAll()
			gowrapmx4j.RegistryFlush()
		}

		if err != nil {
			retErr := fmt.Errorf("QueryMX4J Error: %v %s", newData, err)
			return nil, retErr
		}
	}

	updated := gowrapmx4j.RegistryGetAll()
	return &updated, nil
}

// Cassandra MX4J status returns all data from the gowrapmx4j.registry
// in its raw form marshalled into JSON.
func cassStatus(w http.ResponseWriter, r *http.Request) {
	metrics, err := QueryMX4J(mx4j)
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

// API Endpoint which will execute the optionally specified metric function
// on the data structure for cleanup.
func cleanStatus(w http.ResponseWriter, r *http.Request) {
	metrics := gowrapmx4j.RegistryGetAll()

	mjs := make(map[string]interface{})
	for _, m := range metrics {
		if m.ValFunc != nil {
			log.Infof("%s", m.HumanName)
			mdata, err := m.ValFunc(m.Data)
			if err != nil {
				log.Errorf("Error running value function for %s: %v", m.HumanName, err)
				continue
			}
			mjs[m.HumanName] = mdata
		} else {
			mjs[m.HumanName] = m.Data
		}
	}

	js, err := json.Marshal(mjs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, fmt.Sprintf("CleanFuncAPI: Error marshaling JSON from MX4J data: %#v", err), 500)
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
	mx4j = gowrapmx4j.MX4J{Host: host, Port: port}
	mx4j.Init()

	ll, err := log.ParseLevel(loglvl)
	if err != nil {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(ll)
	}

	// Query singlenton values from MX4J
	mm := gowrapmx4j.NewMX4JMetric("compactions.active", "org.apache.cassandra.internal:type=CompactionExecutor", "array", "ActiveCount")
	mm.ValFunc = gowrapmx4j.DistillAttribute
	gowrapmx4j.RegistrySet(mm, nil)

	mm = gowrapmx4j.NewMX4JMetric("compactions.pending", "org.apache.cassandra.internal:type=CompactionExecutor", "array", "PendingTasks")
	mm.ValFunc = gowrapmx4j.DistillAttribute
	gowrapmx4j.RegistrySet(mm, nil)

	//OR
	// Query MBean attribute maps
	mname := "NodeStatusBinary"
	mm = gowrapmx4j.NewMX4JMetric(mname, "org.apache.cassandra.net:type=FailureDetector", "", "")
	mm.ValFunc = gowrapmx4j.DistillAttributeTypes
	gowrapmx4j.RegistrySet(mm, nil)

	mname = "CompactionExecutor"
	mm = gowrapmx4j.NewMX4JMetric(mname, "org.apache.cassandra.internal:type=CompactionExecutor", "", "")
	mm.ValFunc = gowrapmx4j.DistillAttributeTypes
	gowrapmx4j.RegistrySet(mm, nil)

	// Query Cluster information
	mname = "StorageService"
	mm = gowrapmx4j.NewMX4JMetric(mname, "org.apache.cassandra.db:type=StorageService", "", "")
	mm.ValFunc = gowrapmx4j.DistillAttributeTypes
	gowrapmx4j.RegistrySet(mm, nil)

	// Simple run loop to query MX4J
	go func() {
		for {
			log.Debug("Querying MX4J")
			_, err := QueryMX4J(mx4j)
			if err != nil {
				log.Errorf("Error Querying MX4J: %v", err)
			}
			time.Sleep(time.Second * time.Duration(queryInterval))
		}
	}()

	http.HandleFunc("/", cassStatus)
	http.HandleFunc("/clean", cleanStatus)
	http.HandleFunc("/status", nodeStatus)
	http.ListenAndServe(":8082", nil)
}
