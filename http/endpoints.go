package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/lytics/gowrapmx4j"
)

func NodeStatus(w http.ResponseWriter, r *http.Request) {
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
		//TODO: Define hostnameid
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
