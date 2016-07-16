package gowrapmx4j

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
)

// HTTPRegistryRaw returns raw Cassandra MX4J status endpoint data types
func HTTPRegistryRaw(w http.ResponseWriter, r *http.Request) {
	mbeans := RegistryBeans()
	js, err := json.Marshal(mbeans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "HttpRegistryRaw: Error marshaling JSON from MX4J data: %v", err)
	}
	fmt.Fprintf(w, "%s", js)
}

// HTTPRegistryGetAll makes call to get registry metrics non-blocking
// and can be closed by context timeout.
func HTTPRegistryGetAll(ctx context.Context) *[]MX4JMetric {
	metrics := make(chan *[]MX4JMetric)
	defer close(metrics)

	go func() {
		ret := RegistryGetAll()
		metrics <- &ret
	}()

	select {
	case m := <-metrics:
		return m

	case <-ctx.Done():
		return nil
	}
}

// HTTPRegistryProcessed API Endpoint which will execute the optionally specified ValFunc function
// on the data structure to process the metric's data.
func HTTPRegistryProcessed(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	//metrics := RegistryGetAll(context.WithTimeout(ctx, to))
	metrics := HTTPRegistryGetAll(ctx)
	if metrics == nil {
		http.Error(w, fmt.Sprintf("HttpRegistryProcessed: error retreiving metrics"), 500)
		return
	}

	mjs := make(map[string]interface{})
	for _, m := range *metrics {
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
		http.Error(w, fmt.Sprintf("HttpRegistryProcessed: Error marshaling JSON from MX4J data: %#v", err), 500)
	}
	fmt.Fprintf(w, "%s", js)
}
