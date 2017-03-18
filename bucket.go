package oss

type Bucket struct {
	Name   string `json:"bucket"`
	client *Client
}

func (b *Bucket) NewTicket(verb, object string) *Ticket {
	return b.client.NewTicket(verb, b.Name, object)
}
