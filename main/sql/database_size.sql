-- SELECT client_ip, COUNT(*) AS request_count
-- FROM public.qiniu_cdnauth_requests
-- WHERE created_at >= NOW() - INTERVAL '10 minutes'
-- --WHERE created_at >= NOW() - INTERVAL '1 hour'
-- GROUP BY client_ip
-- ORDER BY request_count DESC
-- LIMIT 10;
-- --最近10分钟内网络请求最频繁的前10名IP
-- --最近1小时内网络请求最频繁的前10名IP

-- 查看数据库大小
SELECT pg_size_pretty(pg_database_size('qiniudb')) AS dbsize;
-- SELECT pg_database_size('qiniudb');

-- -- 查看表大小
-- SELECT pg_size_pretty(pg_total_relation_size('qiniu_cdnauth_requests')) AS total_size, current_timestamp;