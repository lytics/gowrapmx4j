package gowrapmx4j

/*
//Handles reading of the http.Body and passes bytes of io.ReadCloser
//to getAttrUnmarshal() for unmarshaling XML.
func getAttributes(httpBody io.ReadCloser, unmarshalFunc func([]byte) (*MBean, error)) (*MBean, error) {
	xmlBytes, err := ioutil.ReadAll(httpBody)
	if err != nil {
		log.Errorf("Failed to read http response: %#v", err)
		return nil, err
	}

	return unmarshalFunc(xmlBytes)
}

//Unmarshals XML and returns an MBean struct
func getAttrUnmarshal(xmlBytes []byte) (*MBean, error) {
	var mb MBean
	err := xml.Unmarshal([]byte(xmlBytes), &mb)
	if err != nil {
		log.Errorf("Failed to Unmarshal xml: %#v", err)
		log.Errorf("Bytes failed to be unmarshalled: \n%s", xmlBytes)
		return nil, err
	}
	return &mb, nil
}*/
