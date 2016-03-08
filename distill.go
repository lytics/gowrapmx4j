package gowrapmx4j

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var AttributeError = errors.New("gowrapmx4j: Attribute parsing error")

// ExtractAttributes parses the queried MX4JMetric endpoints and yields
// a map of metric fields which can be marshalled cleanly into JSON.
func ExtractAttributes(mb MX4JData) map[string]string {
	data := make(map[string]string)

	switch mb.(type) {
	case *Bean:
		x := mb.(*Bean)
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

func removeBrackets(s string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, "]"), "[")
}

func removeBraces(s string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, "}"), "{")
}

func separateValues(s string) []string {
	r := strings.NewReplacer(" ", "")
	csl := r.Replace(s)
	return strings.Split(csl, ",")
}

func parseArray(s string) {}

func parseMap(s string) map[string]interface{} {
	list := separateValues(removeBraces(s))
	strMap := make(map[string]interface{})

	for _, v := range list {
		kv := strings.Split(v, "=")
		if len(kv) != 2 {
			log.Errorf("Error in parseMap with value: %s", v)
			continue
		}
		strMap[kv[0]] = kv[1]
	}
	return strMap
}

func ExtractAttributeTypes(mb MX4JData) interface{} {
	attributes := make(map[string]interface{})

	switch mb.(type) {
	case *Bean:
		b := mb.(*Bean)
		for _, attr := range b.Attributes {
			log.Info(attr)
		}
		return attributes
	default:
		return errors.New("gowrapmx4j: Attribute extraction type error")
	}
}
