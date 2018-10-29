package main

type GeneralStats struct {
	Bytes                  int     `mapstructure:"bytes" metric_name:"bytesUsedServerInBytes" source_type:"gauge"`
	CurrItems              int     `mapstructure:"curr_items" metric_name:"currentItemsStoredServer" source_type:"gauge"`
	TotalItems             int     `mapstructure:"total_items" metric_name:"itemsStoredPerSecond" source_type:"rate"`
	AverageItemSize        int     ` metric_name:"avgItemSizeInBytes" source_type:"gauge"` // TODO calculate
	LimitMaxbytes          int     `mapstructure:"limit_maxbytes" metric_name:"limitBytesStorage" source_type:"rate"`
	PercentMaxUsed         float64 ` metric_name:"storingItemsPercentMemory" source_type:"gauge"` // TODO calculate
	BytesRead              int     `mapstructure:"bytes_read" metric_name:"bytesReadServerPerSecond" source_type:"rate"`
	BytesWritten           int     `mapstructure:"bytes_written" metric_name:"bytesWrittenServerPerSecond" source_type:"rate"`
	DeleteHits             int     `mapstructure:"delete_hits" metric_name:"deleteCmdRemovedPerSecond" source_type:"rate"`
	DeleteMisses           int     `mapstructure:"delete_misses" metric_name:"deleteCmdNoneRemovedPerSecond" source_type:"rate"`
	TimeInListenDisabledUs int     `mapstructure:"time_in_listen_disabled_us"`
	ExpiredUnfetched       int     `mapstructure:"expired_unfetched"`
	CasBadval              int     `mapstructure:"cas_badval" metric_name:"casWrongRatePerSecond" source_type:"rate"`
	CasMisses              int     `mapstructure:"cas_misses" metric_name:"casMissRatePerSecond" source_type:"rate"`
	CasHits                int     `mapstructure:"cas_hits" metric_name:"casHitRatePerSecond" source_type:"rate"`
	GetHits                int     `mapstructure:"get_hits" metric_name:"getHitPerSecond" source_type:"rate"`
	GetMisses              int     `mapstructure:"get_misses" metric_name:"getMissPerSecond" source_type:"rate"`
	GetExpired             int     `mapstructure:"get_expired"`
	GetFlushed             int     `mapstructure:"get_flushed"`
	GetHitPercent          float64 ` metric_name:"getHitPercent" source_type:"gauge"` // TODO calculate
	CmdGet                 int     `mapstructure:"cmd_get" metric_name:"cmdGetRatePerSecond" source_type:"rate"`
	CmdSet                 int     `mapstructure:"cmd_set" metric_name:"cmdSetRatePerSecond" source_type:"rate"`
	CmdFlush               int     `mapstructure:"cmd_flush" metric_name:"cmdFlushRatePerSecond" source_type:"rate"`
	CmdTouch               int     `mapstructure:"cmd_touch"`
	TouchHits              int     `mapstructure:"touch_hits"`
	TouchMisses            int     `mapstructure:"touch_misses"`
	ListenDisabledNum      int     `mapstructure:"listen_disabled_num"`
	HashPowerLevel         int     `mapstructure:"hash_power_level"`
	AcceptingConns         int     `mapstructure:"accepting_conns"`
	HashIsExpanding        int     `mapstructure:"hash_is_expanding"`
	Evictions              int     `mapstructure:"evictions" metric_name:"evictionsPerSecond" source_type:"rate"`
	EvictedUnfetched       int     `mapstructure:"evicted_unfetched"`
	ConnectionStructures   int     `mapstructure:"connection_structures" metric_name:"connectionStructuresAllocated" source_type:"gauge"`
	CurrConnections        int     `mapstructure:"curr_connections" metric_name:"openConnectionsServer" source_type:"gauge"`
	TotalConnections       int     `mapstructure:"total_connections" metric_name:"connectionRateServerPerSecond" source_type:"rate"`
	ConnYields             int     `mapstructure:"conn_yields" metric_name:"serverMaxConnectionLimitPerSecond" source_type:"rate"`
	PID                    int     `mapstructure:"pid"`
	LogWatcherSkipped      int     `mapstructure:"log_watcher_skipped"`
	CrawlerItemsChecked    int     `mapstructure:"crawler_items_checked"`
	AuthErrors             int     `mapstructure:"auth_errors"`
	HashBytes              int     `mapstructure:"hash_bytes"`
	Uptime                 int     `mapstructure:"uptime"`
	UptimeMilliseconds     int     ` metric_name:"uptimeInMilliseconds" source_type:"gauge"` // TODO calculate
	IncrHits               int     `mapstructure:"incr_hits"`
	DecrMisses             int     `mapstructure:"decr_misses"`
	AuthCmds               int     `mapstructure:"auth_cmds"`
	Reclaimed              int     `mapstructure:"reclaimed"`
	LogWatcherSent         int     `mapstructure:"log_watcher_sent"`
	Time                   int     `mapstructure:"time"`
	PointerSize            int     `mapstructure:"pointer_size" metric_name:"pointerSize" source_type:"gauge"`
	RusageUser             int     `mapstructure:"rusage_user"`
	ReservedFds            int     `mapstructure:"reserved_fds"`
	Libevent               string  `mapstructure:"libevent"`
	RusageSystem           float64 `mapstructure:"rusage_system"`
	DecrHits               int     `mapstructure:"decr_hits"`
	Threads                int     `mapstructure:"threads" metric_name:"threads" source_type:"gauge"`
	CrawlerReclaimed       int     `mapstructure:"crawler_reclaimed"`
	LrutailReflocked       int     `mapstructure:"lrutail_reflocked"`
	Version                string  `mapstructure:"version"`
	IncrMisses             int     `mapstructure:"incr_misses"`
	MallocFails            int     `mapstructure:"malloc_fails"`
	LogWorkerDropped       int     `mapstructure:"log_worker_dropped"`
	LogWorkerWritten       int     `mapstructure:"log_worker_written"`
}

type ItemStats struct {
	EvictedNonzero      int `mapstructure:"evicted_nonzero"`
	OutOfMemory         int `mapstructure:"outofmemory"`
	Number              int `mapstructure:"number"`
	EvictedUnfetched    int `mapstructure:"evicted_unfetched"`
	TailRepairs         int `mapstructure:"tailrepairs"`
	EvictedTime         int `mapstructure:"evicted_time"`
	CrawlerReclaimed    int `mapstructure:"crawler_reclaimed"`
	CrawlerItemsChecked int `mapstructure:"crawler_items_checked"`
	Reclaimed           int `mapstructure:"reclaimed"`
	Evicted             int `mapstructure:"evicted"`
	LrutailReflocked    int `mapstructure:"lrutail_reflocked"`
	Age                 int `mapstructure:"age"`
	ExpiredUnfetched    int `mapstructure:"expired_unfetched"`
}

type SlabStats struct {
	IncrHits      int `mapstructure:"incr_hits"`
	TotalChunks   int `mapstructure:"total_chunks"`
	ChunksPerPage int `mapstructure:"chunks_per_page"`
	DeleteHits    int `mapstructure:"delete_hits"`
	CasBadval     int `mapstructure:"cas_badval"`
	ChunkSize     int `mapstructure:"chunk_size"`
	FreeChunksEnd int `mapstructure:"free_chunks_end"`
	CasHits       int `mapstructure:"cas_hits"`
	TouchHits     int `mapstructure:"touch_hits"`
	TotalPages    int `mapstructure:"total_pages"`
	UsedChunks    int `mapstructure:"used_chunks"`
	FreeChunks    int `mapstructure:"free_chunks"`
	CmdSet        int `mapstructure:"cmd_set"`
	MemRequested  int `mapstructure:"mem_requested"`
	DecrHits      int `mapstructure:"decr_hits"`
	GetHits       int `mapstructure:"get_hits"`
}
