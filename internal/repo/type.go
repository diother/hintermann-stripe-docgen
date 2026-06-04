package repo

type Payout struct {
	Id      string
	Created string
	Gross   string
	Fee     string
	Net     string
}

type Invoice struct {
	Id          string
	Created     string
	ClientName  string
	ClientEmail string
	PayoutId    string
	Gross       string
	Fee         string
	Net         string
}
