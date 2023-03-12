package main

type Event struct {
	Exabgp   string   `json:"exabgp"`
	Time     float64  `json:"time"`
	Host     string   `json:"host"`
	PID      int      `json:"pid"`
	PPID     int      `json:"ppid"`
	Counter  int      `json:"counter"`
	Type     string   `json:"type"`
	Neighbor Neighbor `json:"neighbor"`
}

type Neighbor struct {
	Address   Address `json:"address"`
	ASN       ASN     `json:"asn"`
	Direction string  `json:"direction"`
	Message   Message `json:"message"`
}

type Address struct {
	Local string `json:"local"`
	Peer  string `json:"peer"`
}

type ASN struct {
	Local uint32 `json:"local"`
	Peer  uint32 `json:"peer"`
}

type Message struct {
	Update Update `json:"update"`
}

type Update struct {
	Attribute Attribute `json:"attribute"`
	Announce  Announce  `json:"announce"`
	Withdraw  Withdraw  `json:"withdraw"`
}

type Attribute struct {
	Origin          string     `json:"origin"`
	LocalPreference int        `json:"local-preference"`
	Community       [][]uint16 `json:"community"`
}

type Announce struct {
	IPv4Unicast map[string][]NLRI `json:"ipv4 unicast"`
}

type Withdraw struct {
	IPv4Unicast []NLRI `json:"ipv4 unicast"`
}

// BGP Network Layer Reachability Information (NLRI)

type NLRI struct {
	NLRI string `json:"nlri"`
}
