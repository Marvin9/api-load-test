package pkg

type RequestData struct {
	// generally the index used in for loop is sufficient
	RequestSequence int64
	// time it took, in millisecond
	// Note: this value is close to how much server took to process request
	// we do not include connection time (dns, tls if any)
	// and we measure until first byte of response
	Performance int64
	// TODO: response status, body
}

type Report = []RequestData
