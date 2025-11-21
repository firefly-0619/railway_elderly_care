package models

type AccountEvaluation struct {
	BaseModel
	AccountID   uint    `gorm:"uniqueIndex"`
	Score       float64 `json:"score"` // 总得分
	AssignCount int     // 接收的任务数量
}

func (*AccountEvaluation) TableName() string {
	return "account_evaluation"
}
