/* gowrapmx4j is a base library of types to assist UnMarshalling and Querying MX4J data

MX4J is a very useful layer which makes JMX accessible via HTTP. Unfortunately little is done to
improve the data's representation and it is returned as raw XML via an API frought with perilous
query variables which are poorly documented.

The types and unmarshalling structures defined here have sorted out some of the XML maddness
returned from MX4J and operating on (slightly)more sensible data structures easier.

*/
package gowrapmx4j
