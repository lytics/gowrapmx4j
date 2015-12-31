package gowrapmx4j

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestBasicUnmarshal(t *testing.T) {
	//<?xml version="1.0" encoding="UTF-8"?>
	input := `
<MBean classname="com.yammer.metrics.reporting.JmxReporter$Timer" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.metrics:type=ColumnFamily,keyspace=lio4,scope=node,name=ReadLatency">
	<Attribute classname="double" isnull="false" name="Max" value="100.0"/>
</MBean>
`

	x := MBean{ObjectName: "neh"}
	fmt.Printf("input len %d\n", len(input))

	err := xml.Unmarshal([]byte(input), &x)
	if err != nil {
		t.Errorf("Error unmarshalling xml: %#v", err)
	}

	fmt.Printf("%#v\n", x)
	if x.ObjectName == "neh" {
		t.Errorf("Incorrect default value not overwritten.")
	}

}
