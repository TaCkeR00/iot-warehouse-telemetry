package warehouse

type Warehouse struct {
	Id      int    `json:"id"`
	State   string `json:"state"`
	City    string `json:"city"`
	Address string `json:"address"`
}
