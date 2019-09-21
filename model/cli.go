package model

// 客户端
type Cli struct {
	ID    string `json:"id"`   // 客户端 id
	App   string `json:"app"`  // 产品
	Uuid  string `json:"uuid"` // 客户端 随机token
	IP    string `json:"ip"`
	Start int64  `json:"start"`
	Lease int64  `json:"lease"`
}
