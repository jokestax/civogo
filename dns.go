package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// DNSDomain represents a domain registered within Civo's infrastructure
type DNSDomain struct {
	// The ID of the domain
	ID string `json:"id"`

	// The ID of the account
	AccountID string `json:"account_id"`

	// The Name of the domain
	Name string `json:"name"`
}

type dnsDomainConfig struct {
	Name string `form:"name"`
}

// DNSRecordType represents the allowed record types: a, cname, mx or txt
type DNSRecordType string

// DNSRecord represents a DNS record registered within Civo's infrastructure
type DNSRecord struct {
	ID          string        `json:"id"`
	AccountID   string        `json:"account_id"`
	DNSDomainID string        `json:"domain_id"`
	Name        string        `json:"name"`
	Value       string        `json:"value"`
	Type        DNSRecordType `json:"type"`
	Priority    int           `json:"priority"`
	TTL         int           `json:"ttl"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// DNSRecordConfig describes the parameters for a new DNS record
// none of the fields are mandatory and will be automatically
// set with default values
type DNSRecordConfig struct {
	DNSDomainID string        `form:"-"`
	Type        DNSRecordType `form:"type"`
	Name        string        `form:"name"`
	Value       string        `form:"value"`
	Priority    int           `form:"priority"`
	TTL         int           `form:"ttl"`
}

const (
	// DNSRecordTypeA represents an A record
	DNSRecordTypeA = "a"

	// DNSRecordTypeCName represents an CNAME record
	DNSRecordTypeCName = "cname"

	// DNSRecordTypeMX represents an MX record
	DNSRecordTypeMX = "mx"

	// DNSRecordTypeTXT represents an TXT record
	DNSRecordTypeTXT = "txt"
)

var (
	// ErrDNSDomainNotFound is returned when the domain is not found
	ErrDNSDomainNotFound = fmt.Errorf("domain not found")

	// ErrDNSRecordNotFound is returned when the record is not found
	ErrDNSRecordNotFound = fmt.Errorf("record not found")
)

// ListDNSDomains returns all Domains owned by the calling API account
func (c *Client) ListDNSDomains() ([]DNSDomain, error) {
	url := "/v2/dns"

	resp, err := c.SendGetRequest(url)
	if err != nil {
		return nil, err
	}

	var ds = make([]DNSDomain, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&ds); err != nil {
		return nil, err

	}

	return ds, nil
}

// CreateDNSDomain registers a new Domain
func (c *Client) CreateDNSDomain(name string) (*DNSDomain, error) {
	url := "/v2/dns"
	d := &dnsDomainConfig{Name: name}
	body, err := c.SendPostRequest(url, d)
	if err != nil {
		return nil, err
	}

	var n = &DNSDomain{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(n); err != nil {
		return nil, err
	}

	return n, nil
}

// GetDNSDomain returns the DNS Domain that matches the name
func (c *Client) GetDNSDomain(name string) (*DNSDomain, error) {
	ds, err := c.ListDNSDomains()
	if err != nil {
		return nil, err
	}

	for _, d := range ds {
		if d.Name == name {
			return &d, nil
		}
	}

	return nil, ErrDNSDomainNotFound
}

// UpdateDNSDomain updates the provided domain with name
func (c *Client) UpdateDNSDomain(d *DNSDomain, name string) (*DNSDomain, error) {
	url := fmt.Sprintf("/v2/dns/%s", d.ID)
	dc := &dnsDomainConfig{Name: name}
	body, err := c.SendPutRequest(url, dc)
	if err != nil {
		return nil, err
	}

	var r = &DNSDomain{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(r); err != nil {
		return nil, err
	}

	return r, nil
}

// DeleteDNSDomain deletes the Domain that matches the name
func (c *Client) DeleteDNSDomain(d *DNSDomain) (*SimpleResponse, error) {
	url := fmt.Sprintf("/v2/dns/%s", d.ID)
	resp, err := c.SendDeleteRequest(url)
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}

// CreateDNSRecord creates a new DNS record
func (c *Client) CreateDNSRecord(r *DNSRecordConfig) (*DNSRecord, error) {
	if len(r.DNSDomainID) == 0 {
		return nil, fmt.Errorf("r.DomainID is empty")
	}

	url := fmt.Sprintf("/v2/dns/%s/records", r.DNSDomainID)
	body, err := c.SendPostRequest(url, r)
	if err != nil {
		return nil, err
	}

	var record = &DNSRecord{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(record); err != nil {
		return nil, err
	}

	return record, nil
}

// ListDNSRecords returns all the records associated with domainID
func (c *Client) ListDNSRecords(dnsDomainID string) ([]DNSRecord, error) {
	url := fmt.Sprintf("/v2/dns/%s/records", dnsDomainID)
	resp, err := c.SendGetRequest(url)
	if err != nil {
		return nil, err
	}

	var rs = make([]DNSRecord, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&rs); err != nil {
		return nil, err

	}

	return rs, nil
}

// GetDNSRecord returns the Record that matches the name and the domainID
func (c *Client) GetDNSRecord(domainID, name string) (*DNSRecord, error) {
	rs, err := c.ListDNSRecords(domainID)
	if err != nil {
		return nil, err
	}

	for _, r := range rs {
		if r.Name == name {
			return &r, nil
		}
	}

	return nil, ErrDNSRecordNotFound
}

// UpdateDNSRecord updates the DNS record
func (c *Client) UpdateDNSRecord(rc *DNSRecordConfig, r *DNSRecord) (*DNSRecord, error) {
	url := fmt.Sprintf("/v2/dns/%s/records/%s", r.DNSDomainID, r.ID)
	body, err := c.SendPutRequest(url, rc)
	if err != nil {
		return nil, err
	}

	var dnsRecord = &DNSRecord{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(dnsRecord); err != nil {
		return nil, err
	}

	return dnsRecord, nil
}

// DeleteDNSRecord deletes the DNS record
func (c *Client) DeleteDNSRecord(r *DNSRecord) (*SimpleResponse, error) {
	if len(r.ID) == 0 {
		return nil, fmt.Errorf("ID is empty")
	}

	if len(r.DNSDomainID) == 0 {
		return nil, fmt.Errorf("DNSDomainID is empty")
	}

	url := fmt.Sprintf("/v2/dns/%s/records/%s", r.DNSDomainID, r.ID)
	resp, err := c.SendDeleteRequest(url)
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
