package gowrapmx4j

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestBasicUnmarshal(t *testing.T) {
	//<?xml version="1.0" encoding="UTF-8"?>
	input := `
<?xml version="1.0" encoding="UTF-8"?>
<MBean classname="com.yammer.metrics.reporting.JmxReporter$Timer" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.metrics:type=ColumnFamily,keyspace=lio4,scope=node,name=ReadLatency">
	<Attribute classname="double" isnull="false" name="Max" value="100.0"/>
</MBean>
`

	x := MBean{ObjectName: "neh"}

	err := xml.Unmarshal([]byte(input), &x)
	if err != nil {
		t.Errorf("Error unmarshalling xml: %#v", err)
	}

	if x.ObjectName == "neh" {
		fmt.Printf("%#v\n", x)
		t.Errorf("Incorrect default value not overwritten.")
	}

	if x.Attribute.Name != "Max" {
		t.Errorf("Attribute 'Name' was not unmarshalled correctly")
	}

}

func TestUnmarshalFunction(t *testing.T) {
	input := `
<?xml version="1.0" encoding="UTF-8"?>
<MBean classname="com.yammer.metrics.reporting.JmxReporter$Timer" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.metrics:type=ColumnFamily,keyspace=lio4,scope=node,name=ReadLatency">
	<Attribute classname="double" isnull="false" name="Max" value="100.0"/>
</MBean>
`

	x, err := GetAttrUnmarshal([]byte(input))
	if err != nil {
		t.Errorf("Error running GetAttrUnmarshal: %v\n", err)
	}

	if x.ObjectName == "neh" {
		fmt.Printf("%#v\n", x)
		t.Errorf("Incorrect default value not overwritten.")
	}

	if x.Attribute.Name != "Max" {
		t.Errorf("Attribute 'Name' was not unmarshalled correctly")
	}
}
