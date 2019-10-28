package logic

// 客户端
type Cli struct {
	ID    string `json:"ID"`   // 客户端 ID
	App   string `json:"App"`  // 产品
	Uuid  string `json:"uuid"` // 客户端 随机token
	Start int64  `json:"start"`
	Lease int64  `json:"lease"`
}
