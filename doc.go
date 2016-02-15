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

MX4J is a very useful layer which makes JMX accessible via HTTP. Unfortunately little is done to
improve the data's representation and it is returned as raw XML via an API frought with perilous
query variables which are poorly documented.

The types and unmarshalling structures defined here have sorted out some of the XML maddness
returned from MX4J and operating on (slightly)more sensible data structures easier.
*/
package gowrapmx4j
