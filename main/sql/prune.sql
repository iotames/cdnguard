-- 删除qiniu_cdnauth_block_requests表中30天前的数据
DELETE FROM qiniu_cdnauth_block_requests WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '30 days';

-- 删除qiniu_cdnauth_requests表中7天前的数据
DELETE FROM qiniu_cdnauth_requests WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '7 days';