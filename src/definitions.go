package main

// GeneralStats is a struct to unmarshal the stats results into
type GeneralStats struct {
	Bytes                  *int     `mapstructure:"bytes"                 metric_name:"bytesUsedServerInBytes"            source_type:"gauge"`
	CurrItems              *int     `mapstructure:"curr_items"            metric_name:"currentItemsStoredServer"          source_type:"gauge"`
	TotalItems             *int     `mapstructure:"total_items"           metric_name:"itemsStoredPerSecond"              source_type:"rate"`
	AverageItemSize        *float64 `                                     metric_name:"avgItemSizeInBytes"                source_type:"gauge"`
	LimitMaxbytes          *int     `mapstructure:"limit_maxbytes"        metric_name:"limitBytesStorage"                 source_type:"gauge"`
	PercentMaxUsed         *float64 `                                     metric_name:"storingItemsPercentMemory"         source_type:"gauge"`
	BytesRead              *int     `mapstructure:"bytes_read"            metric_name:"bytesReadServerPerSecond"          source_type:"rate"`
	BytesWritten           *int     `mapstructure:"bytes_written"         metric_name:"bytesWrittenServerPerSecond"       source_type:"rate"`
	DeleteHits             *int     `mapstructure:"delete_hits"           metric_name:"deleteCmdRemovedPerSecond"         source_type:"rate"`
	DeleteMisses           *int     `mapstructure:"delete_misses"         metric_name:"deleteCmdNoneRemovedPerSecond"     source_type:"rate"`
	TimeInListenDisabledUs *int     `mapstructure:"time_in_listen_disabled_us"`
	ExpiredUnfetched       *int     `mapstructure:"expired_unfetched"`
	CasBadval              *int     `mapstructure:"cas_badval"            metric_name:"casWrongRatePerSecond"             source_type:"rate"`
	CasMisses              *int     `mapstructure:"cas_misses"            metric_name:"casMissRatePerSecond"              source_type:"rate"`
	CasHits                *int     `mapstructure:"cas_hits"              metric_name:"casHitRatePerSecond"               source_type:"rate"`
	GetHits                *int     `mapstructure:"get_hits"              metric_name:"getHitPerSecond"                   source_type:"rate"`
	GetMisses              *int     `mapstructure:"get_misses"            metric_name:"getMissPerSecond"                  source_type:"rate"`
	GetExpired             *int     `mapstructure:"get_expired"`
	GetFlushed             *int     `mapstructure:"get_flushed"`
	GetHitPercent          *float64 `                                     metric_name:"getHitPercent"                     source_type:"gauge"`
	CmdGet                 *int     `mapstructure:"cmd_get"               metric_name:"cmdGetRatePerSecond"               source_type:"rate"`
	CmdSet                 *int     `mapstructure:"cmd_set"               metric_name:"cmdSetRatePerSecond"               source_type:"rate"`
	CmdFlush               *int     `mapstructure:"cmd_flush"             metric_name:"cmdFlushRatePerSecond"             source_type:"rate"`
	CmdTouch               *int     `mapstructure:"cmd_touch"`
	TouchHits              *int     `mapstructure:"touch_hits"`
	TouchMisses            *int     `mapstructure:"touch_misses"`
	ListenDisabledNum      *int     `mapstructure:"listen_disabled_num"`
	HashPowerLevel         *int     `mapstructure:"hash_power_level"`
	AcceptingConns         *int     `mapstructure:"accepting_conns"`
	HashIsExpanding        *int     `mapstructure:"hash_is_expanding"`
	Evictions              *int     `mapstructure:"evictions"             metric_name:"evictionsPerSecond"                source_type:"rate"`
	EvictedUnfetched       *int     `mapstructure:"evicted_unfetched"`
	ConnectionStructures   *int     `mapstructure:"connection_structures" metric_name:"connectionStructuresAllocated"     source_type:"gauge"`
	CurrConnections        *int     `mapstructure:"curr_connections"      metric_name:"openConnectionsServer"             source_type:"gauge"`
	RusageUser             *float64 `mapstructure:"rusage_user"           metric_name:"executionTime"                     source_type:"prate"`
	RusageSystem           *float64 `mapstructure:"rusage_system"         metric_name:"usageRate"                         source_type:"prate"`
	PID                    *int     `mapstructure:"pid"`
	LogWatcherSkipped      *int     `mapstructure:"log_watcher_skipped"`
	CrawlerItemsChecked    *int     `mapstructure:"crawler_items_checked"`
	AuthErrors             *int     `mapstructure:"auth_errors"`
	HashBytes              *int     `mapstructure:"hash_bytes"`
	Uptime                 *int     `mapstructure:"uptime"`
	UptimeMilliseconds     *int     `                                     metric_name:"uptimeInMilliseconds"              source_type:"gauge"`
	IncrHits               *int     `mapstructure:"incr_hits"`
	DecrMisses             *int     `mapstructure:"decr_misses"`
	AuthCmds               *int     `mapstructure:"auth_cmds"`
	Reclaimed              *int     `mapstructure:"reclaimed"`
	LogWatcherSent         *int     `mapstructure:"log_watcher_sent"`
	Time                   *int     `mapstructure:"time"`
	PointerSize            *int     `mapstructure:"pointer_size"          metric_name:"pointerSize"                       source_type:"gauge"`
	ReservedFds            *int     `mapstructure:"reserved_fds"`
	Libevent               *string  `mapstructure:"libevent"`
	DecrHits               *int     `mapstructure:"decr_hits"`
	Threads                *int     `mapstructure:"threads"               metric_name:"threads"                           source_type:"gauge"`
	CrawlerReclaimed       *int     `mapstructure:"crawler_reclaimed"`
	LrutailReflocked       *int     `mapstructure:"lrutail_reflocked"`
	Version                *string  `mapstructure:"version"`
	IncrMisses             *int     `mapstructure:"incr_misses"`
	MallocFails            *int     `mapstructure:"malloc_fails"`
	LogWorkerDropped       *int     `mapstructure:"log_worker_dropped"`
	LogWorkerWritten       *int     `mapstructure:"log_worker_written"`
}

// ItemStats is a struct which is used to marshal metrics into a metric set
type ItemStats struct {
	EvictedTime         *int `mapstructure:"evicted_time"      metric_name:"itemsTimeSinceEvictionInMilliseconds"       source_type:"gauge"`
	Evicted             *int `mapstructure:"evicted"           metric_name:"evictionsBeforeExpirationPerSecond"         source_type:"rate"`
	EvictedNonzero      *int `mapstructure:"evicted_nonzero"   metric_name:"evictionsBeforeExplicitExpirationPerSecond" source_type:"rate"`
	ExpiredUnfetched    *int `mapstructure:"expired_unfetched" metric_name:"validItemsEvictedPerSecond"                 source_type:"rate"`
	EvictedUnfetched    *int `mapstructure:"evicted_unfetched" metric_name:"expiredItemsReclaimedPerSecond"             source_type:"rate"`
	OutOfMemory         *int `mapstructure:"outofmemory"       metric_name:"outOfMemoryPerSecond"                       source_type:"rate"`
	Number              *int `mapstructure:"number"            metric_name:"itemsSlabClass"                             source_type:"gauge"`
	NumberHot           *int `mapstructure:"number_hot"        metric_name:"itemsHot"                                   source_type:"gauge"`
	NumberCold          *int `mapstructure:"number_cold"       metric_name:"itemsCold"                                  source_type:"gauge"`
	NumberWarm          *int `mapstructure:"number_warm"       metric_name:"itemsWarm"                                  source_type:"gauge"`
	TailRepairs         *int `mapstructure:"tailrepairs"       metric_name:"selfHealedSlabPerSecond"                    source_type:"rate"`
	CrawlerReclaimed    *int `mapstructure:"crawler_reclaimed" metric_name:"itemsFreedCrawlerPerSecond"                 source_type:"rate"`
	CrawlerItemsChecked *int `mapstructure:"crawler_items_checked"`
	Reclaimed           *int `mapstructure:"reclaimed"         metric_name:"entriesReclaimedPerSecond"                  source_type:"rate"`
	LrutailReflocked    *int `mapstructure:"lrutail_reflocked" metric_name:"itemsRefcountLockedPerSecond"               source_type:"rate"`
	Age                 *int `mapstructure:"age"               metric_name:"itemsOldestInMilliseconds"                  source_type:"gauge"`
	DirectReclaims      *int `mapstructure:"direct_reclaims"   metric_name:"itemsDirectReclaimedPerSecond"              source_type:"rate"`
	MovesToCold         *int `mapstructure:"moves_to_cold"     metric_name:"itemsColdPerSecond"                         source_type:"rate"`
	MovesToWarm         *int `mapstructure:"moves_to_warm"     metric_name:"itemsWarmPerSecond"                         source_type:"rate"`
	MovesWithinLRU      *int `mapstructure:"moves_within_lru"  metric_name:"activeItemsBumpedPerSecond"                 source_type:"rate"`
}

// SlabStats is a struct which is used to marshal metrics into a metric set
type SlabStats struct {
	GetHits             *int `mapstructure:"get_hits"        metric_name:"getHitRateSlabPerSecond"          source_type:"rate"`
	CmdSet              *int `mapstructure:"cmd_set"         metric_name:"cmdSetRateSlabPerSecond"          source_type:"rate"`
	DeleteHits          *int `mapstructure:"delete_hits"     metric_name:"deleteRateSlabPerSecond"          source_type:"rate"`
	IncrHits            *int `mapstructure:"incr_hits"       metric_name:"incrsModifySlabPerSecond"         source_type:"rate"`
	DecrHits            *int `mapstructure:"decr_hits"       metric_name:"decrsModifySlabPerSecond"         source_type:"rate"`
	CasHits             *int `mapstructure:"cas_hits"        metric_name:"casModifiedSlabPerSecond"         source_type:"rate"`
	CasBadval           *int `mapstructure:"cas_badval"      metric_name:"casBadValPerSecond"               source_type:"rate"`
	TouchHits           *int `mapstructure:"touch_hits"      metric_name:"touchHitSlabPerSecond"            source_type:"rate"`
	UsedChunks          *int `mapstructure:"used_chunks"     metric_name:"usedChunksItems"                  source_type:"gauge"`
	UsedChunksPerSecond *int `mapstructure:"used_chunks"     metric_name:"usedChunksPerSecond"              source_type:"rate"`
	ChunkSize           *int `mapstructure:"chunk_size"      metric_name:"chunkSizeInBytes"                 source_type:"gauge"`
	ChunksPerPage       *int `mapstructure:"chunks_per_page" metric_name:"chunksPerPage"                    source_type:"gauge"`
	TotalPages          *int `mapstructure:"total_pages"     metric_name:"totalPagesSlab"                   source_type:"gauge"`
	TotalChunks         *int `mapstructure:"total_chunks"    metric_name:"totalChunksSlab"                  source_type:"gauge"`
	FreeChunks          *int `mapstructure:"free_chunks"     metric_name:"freedChunks"                      source_type:"gauge"`
	FreeChunksEnd       *int `mapstructure:"free_chunks_end" metric_name:"freedChunksEnd"                   source_type:"gauge"`
	MemRequested        *int `mapstructure:"mem_requested"   metric_name:"memRequestedSlabInBytesPerSecond" source_type:"rate"`
}

// ClusterSlabStats is a struct which is used to marshal metrics into a metric set
type ClusterSlabStats struct {
	ActiveSlabs   *int `mapstructure:"active_slabs"   metric_name:"activeSlabs"              source_type:"gauge"`
	TotalMalloced *int `mapstructure:"total_malloced" metric_name:"memAllocatedSlabsInBytes" source_type:"gauge"`
}
