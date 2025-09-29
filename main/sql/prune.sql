-- 删除qiniu_cdnauth_block_requests表中5天前的数据
DELETE FROM qiniu_cdnauth_block_requests WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '5 days';

-- 删除qiniu_cdnauth_requests表中5天前的数据
DELETE FROM qiniu_cdnauth_requests WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '5 days';