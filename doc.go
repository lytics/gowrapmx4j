//   Copyright 2016 Lytics
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

/*
gowrapmx4j is a base library of types to assist UnMarshalling and Querying MX4J data.

MX4J is a very useful service which makes JMX data accessible via HTTP. Unfortunately little is done to
improve the data's representation and it is returned as dense raw XML via an API frought with perilous
query variables which are poorly documented.

The types and unmarshalling structures defined here have sorted out some of the XML saddness
returned from MX4J and makes it easier to operate on the data stuctures.

Basic Primer:
The Registry is a concurrent safe map of MX4J data which is updated when queries are made.
This is to reduce the number of calls to MX4J if multiple goroutines want to access the data.

MX4J Unmarshalling types:
Sadly MX4J likes to reuse XML tag names despite different data structures.
eg: XML "MBean" tag. This leaves few options to keep the library's API clean and readable.

Type: "Bean"; a root level Map of MX4JAttributes
Type: "MBean"; a single effective MX4J variable path has a nested K-V data type

*/
package gowrapmx4j
