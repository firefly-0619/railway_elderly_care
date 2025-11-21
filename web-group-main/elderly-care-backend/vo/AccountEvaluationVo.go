package vo

type AccountEvaluationVo struct {
	Score       float64 `json:"score"` // 总得分
	AssignCount int     // 获取任务数量
}
