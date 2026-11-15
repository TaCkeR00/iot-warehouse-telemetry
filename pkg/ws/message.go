package ws

type Message struct {
	Type  string `json:"type"`
	Token string `json:"token,omitempty"`
	Date  string `json:"data,omitempty"`
}
