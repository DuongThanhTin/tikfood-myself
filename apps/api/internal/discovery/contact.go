package discovery

// WithContact returns a copy of the venue carrying phone/website for ingestion.
// These fields are unexported and never appear in the JSON payload.
func (v Venue) WithContact(phone string, website string) Venue {
	v.contactPhone = phone
	v.contactWebsite = website
	return v
}

func (v Venue) Phone() string   { return v.contactPhone }
func (v Venue) Website() string { return v.contactWebsite }
