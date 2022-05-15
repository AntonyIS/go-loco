package app

type Locomotive struct {
	LocoID      string `json:"loco_id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	CreatedAT   int64  `json:"created_at"`
}
