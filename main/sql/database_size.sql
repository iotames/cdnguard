-- 查看数据库大小
SELECT pg_size_pretty(pg_database_size($1)) AS dbsize;
-- SELECT pg_database_size('qiniudb');

-- -- 查看表大小
-- SELECT pg_size_pretty(pg_total_relation_size('qiniu_cdnauth_requests')) AS total_size, current_timestamp;

-- 大概估算表记录行数
-- SELECT reltuples AS estimate_count FROM pg_class WHERE relname = 'qiniu_cdnauth_requests';