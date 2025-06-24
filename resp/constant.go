package resp

// response constants
var (
	OK   = SimpleString("OK")
	NULL = Null{}
)

// command constants
var (
	SET     = BulkString("SET")
	EX      = BulkString("EX")
	PX      = BulkString("PX")
	EXAT    = BulkString("EXAT")
	PXAT    = BulkString("PXAT")
	NX      = BulkString("NX")
	XX      = BulkString("XX")
	KEEPTTL = BulkString("KEEPTTL")

	GET    = BulkString("GET")
	DEL    = BulkString("DEL")
	EXISTS = BulkString("EXISTS")
	EXPIRE = BulkString("EXPIRE")
	TTL    = BulkString("TTL")
	TYPE   = BulkString("TYPE")
	KEYS   = BulkString("KEYS")
	SCAN   = BulkString("SCAN")

	INCR   = BulkString("INCR")
	DECR   = BulkString("DECR")
	INCRBY = BulkString("INCRBY")
	DECRBY = BulkString("DECRBY")
	MSET   = BulkString("MSET")
	MGET   = BulkString("MGET")
	APPEND = BulkString("APPEND")

	// hash commands
	HGET    = BulkString("HGET")
	HSET    = BulkString("HSET")
	HGETALL = BulkString("HGETALL")
	HDEL    = BulkString("HDEL")
	HLEN    = BulkString("HLEN")
	HKEYS   = BulkString("HKEYS")
	HVALS   = BulkString("HVALS")
	HEXISTS = BulkString("HEXISTS")

	// list commands
	LPUSH  = BulkString("LPUSH")
	RPUSH  = BulkString("RPUSH")
	LPOP   = BulkString("LPOP")
	RPOP   = BulkString("RPOP")
	LRANGE = BulkString("LRANGE")
	LLEN   = BulkString("LLEN")

	// set commands
	SADD      = BulkString("SADD")
	SREM      = BulkString("SREM")
	SMEMBERS  = BulkString("SMEMBERS")
	SISMEMBER = BulkString("SISMEMBER")
	SCARD     = BulkString("SCARD")
	SINTER    = BulkString("SINTER")

	// sorted set commands
	ZADD      = BulkString("ZADD")
	ZREM      = BulkString("ZREM")
	ZRANGE    = BulkString("ZRANGE")
	ZREVRANGE = BulkString("ZREVRANGE")
	ZSCORE    = BulkString("ZSCORE")
	ZCARD     = BulkString("ZCARD")

	COMMAND = BulkString("COMMAND")
)
