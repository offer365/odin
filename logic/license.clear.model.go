package logic

type Clear struct {
	Lic    *License `json:"lic"`
	Cipher string   `json:"cipher"`
	Date   int64    `json:"date"`
}
