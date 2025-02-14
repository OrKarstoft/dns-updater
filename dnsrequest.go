package domain

import (
	"fmt"
)

type DNSRequest struct {
	recordName string
	domain     string
	zone       string
	ip         string
	recordType string
}

var ErrInvalidInput error = fmt.Errorf("invalid input")

func NewDNSRequest(recordName, domain, zone, ip, recordType string) (*DNSRequest, error) {
	if recordName == "" || domain == "" || zone == "" || ip == "" || recordType == "" {
		return nil, ErrInvalidInput
	}

	return &DNSRequest{
		recordName: recordName,
		domain:     domain,
		zone:       zone,
		ip:         ip,
		recordType: recordType,
	}, nil
}

func (r *DNSRequest) GetIP() string {
	return r.ip
}

func (r *DNSRequest) GetRecordName() string {
	return r.recordName
}

func (r *DNSRequest) GetDomain() string {
	return r.domain
}

func (r *DNSRequest) GetRecordType() string {
	return r.recordType
}

func (r *DNSRequest) GetZone() string {
	return r.zone
}
