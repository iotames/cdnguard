-- SELECT 
--     client_ip, 
--     COUNT(*) AS request_count
-- FROM public.qiniu_cdnauth_requests
-- WHERE created_at > CURRENT_DATE + INTERVAL '1 hours'
--   AND created_at < CURRENT_DATE + INTERVAL '5 hours'
-- GROUP BY client_ip
-- ORDER BY request_count DESC
-- LIMIT 20;
-- -- WHERE created_at >= NOW() - INTERVAL '10 minutes'

-- 当天凌晨01:00到05:00的请求数统计 1 hours,5 hours
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
LIMIT 15;