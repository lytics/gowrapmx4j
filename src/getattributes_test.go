package gowrapmx4j

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestAttributesFromioReadCloser(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
<MBean classname="com.yammer.metrics.reporting.JmxReporter$Timer" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.metrics:type=ColumnFamily,keyspace=yourkeyspace,scope=node,name=ReadLatency">
	<Attribute classname="double" isnull="false" name="Max" value="100.0"/>
</MBean>
`

	f, err := ioutil.TempFile("/tmp", "gowrapmx4jtest-")
	defer f.Close()
	if err != nil {
		t.Errorf("Failed to open temp file in /tmp/: %#v\n", err)
	}

	f.WriteString(input)
	_, err = f.Seek(0, 0)
	if err != nil {
		t.Errorf("Error setting seek on tmp file: %#v", err)
	}
	defer os.Remove(f.Name())

	mbean, err := getAttributes(f, getAttrUnmarshal)
	if err != nil {
		t.Errorf("Error reading tmp file in getAttributes: %#v\n", err)
	}

	if mbean.Attribute.Name != "Max" {
		t.Errorf("Attribute 'Name' was not unmarshalled correctly")
	}
}

func TestBasicUnmarshal(t *testing.T) {
	//<?xml version="1.0" encoding="UTF-8"?>
	input := `
<?xml version="1.0" encoding="UTF-8"?>
<MBean classname="com.yammer.metrics.reporting.JmxReporter$Timer" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.metrics:type=ColumnFamily,keyspace=yourkeyspace,scope=node,name=ReadLatency">
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
<MBean classname="com.yammer.metrics.reporting.JmxReporter$Timer" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.metrics:type=ColumnFamily,keyspace=yourkeyspace,scope=node,name=ReadLatency">
	<Attribute classname="double" isnull="false" name="Max" value="100.0"/>
</MBean>
`

	x, err := getAttrUnmarshal([]byte(input))
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

func TestUnmarshalFunctionMap(t *testing.T) {
	input := `<MBean classname="org.apache.cassandra.gms.FailureDetector" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.net:type=FailureDetector">
	<Attribute classname="java.util.Map" isnull="false" name="SimpleStates">
		<Map length="1">
			<Element element="UP" elementclass="java.lang.String" index="0" key="/127.0.0.1" keyclass="java.lang.String"/>
		</Map>
	</Attribute>
</MBean>`

	x, err := getAttrUnmarshal([]byte(input))
	if err != nil {
		t.Errorf("Error running GetAttrUnmarshal: %v\n", err)
	}

	if x.ObjectName != "org.apache.cassandra.net:type=FailureDetector" {
		fmt.Printf("%#v\n", x)
		t.Errorf("Parsing failure of objectname")
	}

	if x.Attribute.Name != "SimpleStates" {
		t.Errorf("Attribute 'Name' was not unmarshalled correctly")
	}

	if x.Attribute.Map.Length != "1" {
		t.Errorf("Map Lenght incorrect: %s\n", x.Attribute.Map.Length)
	}

	e0 := x.Attribute.Map.Elements[0]
	if e0.Element != "UP" {
		t.Errorf("Value 'element' not 'UP': %#v", e0)
	}
	if e0.Key != "/127.0.0.1" {
		t.Errorf("Key of map incorrect")
	}
}
