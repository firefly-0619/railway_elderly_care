package constants

import "time"

const (
	DEFAULT_PAGE_SIZE = 10
)

// 重试
const (
	BATCH_SIZE     = 100
	MAX_TIME_OUT   = 10 * time.Second
	RETRY_INTERVAL = 500 * time.Millisecond
)
