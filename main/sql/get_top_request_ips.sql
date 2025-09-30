-- WHERE created_at >= NOW() - INTERVAL '10 minutes'
-- WHERE created_at > CURRENT_DATE + INTERVAL '1 hours'
--   AND created_at < CURRENT_DATE + INTERVAL '5 hours'


-- 当天凌晨01:00到05:00的请求数统计 1 hours,5 hours
-- r.http_referer is NULL 为异常请求
-- 排除了白名单里面的IP
SELECT
    r.client_ip ip, 
    COUNT(*) AS request_count
FROM public.qiniu_cdnauth_requests r
WHERE r.created_at > $1
  AND r.created_at < $2
  AND r.http_referer is NULL
  AND NOT EXISTS (
    SELECT 1 
    FROM public.qiniu_cdnauth_ip_white_list w 
    WHERE w.ip = r.client_ip
  )
GROUP BY r.client_ip
ORDER BY request_count DESC
LIMIT 30;
