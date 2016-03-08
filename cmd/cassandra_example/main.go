package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"reflect"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lytics/gowrapmx4j"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/ropes/katoptron"
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
		log.Debugf("Metric being queried:\n %#v\n", mm)

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

// Cleanly extracts the name and value from a singleton MX4J Bean struct
func PercentileClean(mb gowrapmx4j.MX4JData) map[string]string {
	switch mb.(type) {
	case *gowrapmx4j.Bean:
		x := mb.(*gowrapmx4j.Bean)
		return map[string]string{x.Attributes[0].Name: x.Attributes[0].Value}
	default:
		return map[string]string{"ERR": "Unable to get type of MX4J Data"}
	}
}

// ExtractAttributes parses the queried MX4JMetric endpoints and yields
// a map of metric fields which can be marshalled cleanly into JSON.
func ExtractAttributes(mb gowrapmx4j.MX4JData) map[string]string {
	data := make(map[string]string)

	katoptron.Display("Bean", reflect.ValueOf(mb))
	//log.Infof("%v", reflect.TypeOf(mb))
	switch mb.(type) {
	case *gowrapmx4j.Bean:
		x := mb.(*gowrapmx4j.Bean)
		for _, attr := range x.Attributes {
			log.Debugf("%s %s", attr.Name, attr.Value)
			if attr.Value != "" {
				data[attr.Name] = attr.Value
			}
		}
		return data
	default:
		return map[string]string{"ERR": "extractAttributes: Unknown type of MX4J Data"}
	}

}

// Cassandra MX4J status endpoint
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
			mjs[m.HumanName] = m.ValFunc(m.Data)
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

func main() {
	flag.StringVar(&host, "host", "localhost", "mx4j host address")
	flag.StringVar(&port, "port", "8081", "mx4j port to query")
	flag.StringVar(&loglvl, "loglvl", "info", "Log level to use")
	flag.StringVar(&gprefix, "gprefix", "mx4jcass.", "Graphite prefix tag")
	flag.StringVar(&hostnameid, "hostid", "", "Hostname to use as graphite classifier for metrics")
	flag.IntVar(&queryInterval, "queryInterval", 10, "Interval seconds between querying MX4J")
	flag.Parse()

	//TODO: Sort out your metrics endpoint
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

	//Pull singlenton values from MX4J
	mm := gowrapmx4j.NewMX4JMetric("compactions.active", "org.apache.cassandra.internal:type=CompactionExecutor", "array", "ActiveCount")
	mm.ValFunc = PercentileClean
	gowrapmx4j.RegistrySet(mm, nil)
	mm = gowrapmx4j.NewMX4JMetric("compactions.pending", "org.apache.cassandra.internal:type=CompactionExecutor", "array", "PendingTasks")
	mm.ValFunc = PercentileClean
	gowrapmx4j.RegistrySet(mm, nil)

	//OR
	// Query MBean attribute maps
	mname := "compactionExecutor"
	mm = gowrapmx4j.NewMX4JMetric(mname, "org.apache.cassandra.internal:type=CompactionExecutor", "", "")
	mm.ValFunc = ExtractAttributes
	gowrapmx4j.RegistrySet(mm, nil)

	// Query Cluster information
	mname = "StorageService"
	mm = gowrapmx4j.NewMX4JMetric(mname, "org.apache.cassandra.db:type=StorageService", "", "")
	mm.ValFunc = ExtractAttributes
	gowrapmx4j.RegistrySet(mm, nil)

	go func() {
		for {
			log.Info("Querying MX4J")
			QueryMX4J(mx4j)
			time.Sleep(time.Second * time.Duration(queryInterval))
		}
	}()

	http.HandleFunc("/", cassStatus)
	http.HandleFunc("/clean", cleanStatus)
	http.ListenAndServe(":8082", nil)
}
