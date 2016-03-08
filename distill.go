package gowrapmx4j

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var AttributeError = errors.New("gowrapmx4j: Attribute parsing error")

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

func parseArray(s string) []string {
	return separateValues(removeBrackets(s))
}

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

// Cleanly extracts the name and value from a singleton MX4J Bean struct
func AttributeClean(mb MX4JData) (map[string]interface{}, error) {
	dataMap := make(map[string]interface{})
	switch mb.(type) {
	case *Bean:
		x := mb.(*Bean)
		dataMap[x.Attributes[0].Name] = x.Attributes[0].Value
		return dataMap, nil
	default:
		return nil, errors.New("gowrapmx4j.PercentileClean() type error")
	}
}

// ExtractAttributes parses the queried MX4JMetric endpoints and yields
// a map of metric fields to their original string values.
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

// ExtractAtributeTypes parses Bean struct []Attributes data and returns
// map parsed from attribute information which can be marsahlled into JSON.
func ExtractAttributeTypes(mb MX4JData) (map[string]interface{}, error) {
	attributes := make(map[string]interface{})

	switch mb.(type) {
	case *Bean:
		b := mb.(*Bean)
		for _, attr := range b.Attributes {
			log.Info(attr)

			strippedVal := removeBrackets(removeBraces(attr.Value))
			if strippedVal == "" {
				log.Debugf("Attribute %s is empty", attr.Name)
				continue
			}
			switch attr.Aggregation {
			case "":
				attributes[attr.Name] = attr.Value
			case "collection":
				attributes[attr.Name] = parseArray(attr.Value)
			case "map":
				attributes[attr.Name] = parseMap(attr.Value)
			}
		}
		return attributes, nil
	default:
		return nil, errors.New("gowrapmx4j.ExtractAttributeTypes() Error: attribute type not handled")
	}
}
