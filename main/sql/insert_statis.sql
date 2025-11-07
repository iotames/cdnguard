-- 判断是否存在昨天的统计数据，如果不存在则进行统计插入
INSERT INTO qiniu_cdnauth_statis (
    bucket_id, 
    statis_date, 
    request_count, 
    blocked_count, 
    blocked_black_count, 
    blocked_scanvul_count, 
    blocked_webspider_count, 
    blocked_useragent_count,
    request_size,
    blocked_size
)
SELECT 
    COALESCE(r.bucket_id, b.bucket_id) as bucket_id,
    CURRENT_DATE - INTERVAL '1 day' as statis_date,
    COALESCE(r.request_count, 0) as request_count,
    COALESCE(b.blocked_count, 0) as blocked_count,
    COALESCE(b.blocked_black_count, 0) as blocked_black_count,
    COALESCE(b.blocked_scanvul_count, 0) as blocked_scanvul_count,
    COALESCE(b.blocked_webspider_count, 0) as blocked_webspider_count,
    COALESCE(b.blocked_useragent_count, 0) as blocked_useragent_count,
    COALESCE(r.request_size, 0) as request_size,
    COALESCE(b.blocked_size, 0) as blocked_size
FROM 
    -- 统计正常请求
    (SELECT 
        f.bucket_id,
        COUNT(*) as request_count,
        COALESCE(SUM(f.file_size), 0) as request_size
    FROM qiniu_cdnauth_requests req
    INNER JOIN qiniu_cdnauth_files f ON f.file_key = 
    regexp_replace(split_part(req.request_url, '?', 1), '^https?://[^/]+/', '')
    WHERE 
        req.created_at >= CURRENT_DATE - INTERVAL '1 day'
        AND req.created_at < CURRENT_DATE
        AND req.deleted_at IS NULL
    GROUP BY f.bucket_id) r
FULL OUTER JOIN
    -- 统计拦截请求
    (SELECT 
        f.bucket_id,
        COUNT(*) as blocked_count,
        SUM(CASE WHEN br.block_type = 0 THEN 1 ELSE 0 END) as blocked_black_count,
        SUM(CASE WHEN br.block_type = 1 THEN 1 ELSE 0 END) as blocked_scanvul_count,
        SUM(CASE WHEN br.block_type = 2 THEN 1 ELSE 0 END) as blocked_webspider_count,
        SUM(CASE WHEN br.block_type = 3 THEN 1 ELSE 0 END) as blocked_useragent_count,
        COALESCE(SUM(f.file_size), 0) as blocked_size
    FROM qiniu_cdnauth_block_requests br
    INNER JOIN qiniu_cdnauth_files f ON f.file_key = 
    regexp_replace(split_part(br.request_url, '?', 1), '^https?://[^/]+/', '')
    WHERE 
        br.created_at >= CURRENT_DATE - INTERVAL '1 day'
        AND br.created_at < CURRENT_DATE
        AND br.deleted_at IS NULL
    GROUP BY f.bucket_id) b ON r.bucket_id = b.bucket_id
WHERE NOT EXISTS (
    SELECT 1 
    FROM qiniu_cdnauth_statis 
    WHERE statis_date = CURRENT_DATE - INTERVAL '1 day'
    LIMIT 1
);