package gowrapmx4j

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestBraceRemoval(t *testing.T) {
	teststr := "{10.100.190.140=115.8 GB, 10.100.37.24=112.99 GB, 10.100.168.68=125.47 GB, 10.100.100.117=125.49 GB, 10.100.18.46=127.44 GB}"
	expected := "10.100.190.140=115.8 GB, 10.100.37.24=112.99 GB, 10.100.168.68=125.47 GB, 10.100.100.117=125.49 GB, 10.100.18.46=127.44 GB"
	out := removeBraces(teststr)
	if out != expected {
		t.Errorf("Braces not removed: \"%s\"", out)
	}
}

func TestBracketRemoval(t *testing.T) {
	teststr := "[10.100.168.68, 10.100.100.117, 10.100.37.24, 10.100.190.140, 10.100.18.46]"
	expected := "10.100.168.68, 10.100.100.117, 10.100.37.24, 10.100.190.140, 10.100.18.46"
	out := removeBrackets(teststr)
	if out != expected {
		t.Errorf("Brackets not removed")
	}
}

func TestValueSplitArray(t *testing.T) {
	teststr := "[10.100.168.68, 10.100.100.117, 10.100.37.24, 10.100.190.140, 10.100.18.46]"
	expected := []string{"10.100.168.68", "10.100.100.117", "10.100.37.24", "10.100.190.140", "10.100.18.46"}
	vals := separateValues(removeBrackets(teststr))

	if !reflect.DeepEqual(expected, vals) {
		t.Errorf("Separated values unequal\n%v != %v", expected, vals)
	}
}

func TestMapParse(t *testing.T) {
	teststr := "{10.100.190.140=115.8 GB, 10.100.37.24=112.99 GB, 10.100.168.68=125.47 GB, 10.100.100.117=125.49 GB, 10.100.18.46=127.44 GB}"
	parsed := parseMap(teststr)
	v, ok := parsed["10.100.190.140"]
	vstr := v.(string)
	if len(parsed) != 5 {
		t.Errorf("Incorrect number of elements parsed: %d", len(parsed))
	}
	if !ok {
		t.Errorf("10.100.190.140 key does not exist in parsed map")
	}
	if vstr != "115.8GB" {
		t.Errorf("Map value %v not equal to expected: \"115.8 GB\"", vstr)
	}
}

var rawBean = `<?xml version="1.0" encoding="UTF-8"?>
<MBean classname="com.yammer.metrics.reporting.JmxReporter$Timer" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.metrics:type=ColumnFamily,keyspace=yourkeyspace,scope=node,name=RangeLatency">
  <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="50thPercentile" strinit="true" type="double" value="2.4563736922E7"/>
  <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="75thPercentile" strinit="true" type="double" value="4.9114161986E7"/>
  <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="95thPercentile" strinit="true" type="double" value="4.9114161986E7"/>
  <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="98thPercentile" strinit="true" type="double" value="4.9114161986E7"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="999thPercentile" strinit="true" type="double" value="4.9114161986E7"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="99thPercentile" strinit="true" type="double" value="4.9114161986E7"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="Count" strinit="true" type="long" value="23"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="EventType" strinit="true" type="java.lang.String" value="calls"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="FifteenMinuteRate" strinit="true" type="double" value="4.44659081257E-313"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="FiveMinuteRate" strinit="true" type="double" value="1.4821969375E-313"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="LatencyUnit" strinit="false" type="java.util.concurrent.TimeUnit" value="MICROSECONDS"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="Max" strinit="true" type="double" value="4.9114161986E7"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="Mean" strinit="true" type="double" value="4411217.701"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="MeanRate" strinit="true" type="double" value="2.3773582427257664E-6"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="Min" strinit="true" type="double" value="13311.858"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="OneMinuteRate" strinit="true" type="double" value="2.964393875E-314"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="RateUnit" strinit="false" type="java.util.concurrent.TimeUnit" value="SECONDS"/>
  <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="StdDev" strinit="true" type="double" value="1.0582224553531626E7"/>
  <Operation description="Operation exposed for management" impact="unknown" name="objectName" return="javax.management.ObjectName"/>
  <Operation description="Operation exposed for management" impact="unknown" name="values" return="[D"/>
</MBean>`

func TestUnMarshalXML(t *testing.T) {

	x := Bean{ObjectName: "neh"}

	err := xml.Unmarshal([]byte(rawBean), &x)
	if err != nil {
		t.Errorf("Error unmarshalling: %v\n", err)
	}

	//fmt.Printf("%#v\n", x)
	if len(x.Attributes) < 10 {
		t.Errorf("Error unmarshalling attributes")
	}

}

func TestUnmarshallingFunctions(t *testing.T) {
	b, err := beanUnmarshal([]byte(rawBean))
	if err != nil {
		t.Errorf("Error unmarshalling: %v\n", err)
	}

	//fmt.Printf("%#v\n", b)
	for _, x := range b.Attributes {
		fmt.Printf("%s %s\n", x.Name, x.Value)
	}
	if len(b.Attributes) < 5 {
		t.Errorf("Seems like number of Attributes unmarshalled is a little low..\n%#v\n", b.Attributes)
	}
}

func TestInterfaceAttributes(t *testing.T) {
	input := `<MBean classname="org.apache.cassandra.gms.FailureDetector" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.net:type=FailureDetector">
	<Attribute classname="java.util.Map" isnull="false" name="SimpleStates">
		<Map length="1">
			<Element element="UP" elementclass="java.lang.String" index="0" key="/127.0.0.1" keyclass="java.lang.String"/>
		</Map>
	</Attribute>
</MBean>`

	var bean Bean
	err := xml.Unmarshal([]byte(input), &bean)
	if err != nil {
		t.Errorf("Error unmarshalling MX4J data: %v", err)
	}
	_, err = beanUnmarshal([]byte(rawBean))
	if err != nil {
		t.Errorf("Error unmarshalling: %v\n", err)
	}

	// Enable for helpfull debugging
	//katoptron.Display("MBean", bean)
	//katoptron.Display("MBean", *b)
}

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

	mbean, err := readHttp(f, beanUnmarshal)
	if err != nil {
		t.Errorf("Error reading tmp file in getAttributes: %#v\n", err)
	}

	attr0 := mbean.Attributes[0]
	if attr0.Name != "Max" {
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

	x := Bean{ObjectName: "neh"}

	err := xml.Unmarshal([]byte(input), &x)
	if err != nil {
		t.Errorf("Error unmarshalling xml: %#v", err)
	}

	if x.ObjectName == "neh" {
		fmt.Printf("%#v\n", x)
		t.Errorf("Incorrect default value not overwritten.")
	}

	if x.Attributes[0].Name != "Max" {
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

	x, err := beanUnmarshal([]byte(input))
	if err != nil {
		t.Errorf("Error running GetAttrUnmarshal: %v\n", err)
	}

	if x.ObjectName == "neh" {
		fmt.Printf("%#v\n", x)
		t.Errorf("Incorrect default value not overwritten.")
	}

	if x.Attributes[0].Name != "Max" {
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

	x, err := beanUnmarshal([]byte(input))
	if err != nil {
		t.Errorf("Error running GetAttrUnmarshal: %v\n", err)
	}

	if x.ObjectName != "org.apache.cassandra.net:type=FailureDetector" {
		fmt.Printf("%#v\n", x)
		t.Errorf("Parsing failure of objectname")
	}

	if x.Attributes[0].Name != "SimpleStates" {
		t.Errorf("Attribute 'Name' was not unmarshalled correctly")
	}

	if x.Attributes[0].Map.Length != "1" {
		t.Errorf("Map Lenght incorrect: %s\n", x.Attributes[0].Map.Length)
	}

	e0 := x.Attributes[0].Map.Elements[0]
	if e0.Element != "UP" {
		t.Errorf("Value 'element' not 'UP': %#v", e0)
	}
	if e0.Key != "/127.0.0.1" {
		t.Errorf("Key of map incorrect")
	}
}
