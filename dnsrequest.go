package domain

type DNSRequest struct {
	recordName string
	domain     string
	zone       string
	ip         string
	recordType string
}

func NewDNSRequest(recordName, domain, zone, ip, recordType string) *DNSRequest {
	if recordName == "" || domain == "" || zone == "" || ip == "" || recordType == "" {
		return nil
	}

	return &DNSRequest{
		recordName: recordName,
		domain:     domain,
		zone:       zone,
		ip:         ip,
		recordType: recordType,
	}
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
