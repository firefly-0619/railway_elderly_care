package constants

import "time"

const (
	WEEKLY_RANKING  = "weekly_ranking"
	MONTHLY_RANKING = "monthly_ranking"

	// cache 缓存
	CACHE_MUSIC_RANK_PREFIX     = "cache:rank:music_list"
	CACHE_SOURCE_AND_LYRICS     = "cache:music:source_lyrics"
	CACHE_SOURCE_AND_LYRICS_TTL = 12 * time.Hour
	//锁
	LOCK_MUSIC_WEEK_RANK         = "lock:music_week_rank"
	LOCK_MUSIC_SOURCE_AND_LYRICS = "lock:source_lyrics"

	//redis计数器
	CHAT_ID_COUNT = "chat_id_count"

	//账户信息
	ACCOUNT_LOGINTYPE = "account:login_type"
)
