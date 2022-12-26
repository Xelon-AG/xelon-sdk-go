package xelon

type NIC struct {
	ControllerKey int             `json:"niccontrollerkey1,omitempty"`
	Key           int             `json:"nickey1,omitempty"`
	IPs           map[string][]IP `json:"ips,omitempty"`
	Name          string          `json:"nicname,omitempty"`
	Networks      []NICNetwork    `json:"networks,omitempty"`
	Number        int             `json:"nicnumber,omitempty"`
	Unit          int             `json:"nicunit1,omitempty"`
}

type NICNetwork struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Value int    `json:"value,omitempty"`
}

type IP struct {
	ID      int    `json:"value,omitempty"`
	Address string `json:"text,omitempty"`
}
