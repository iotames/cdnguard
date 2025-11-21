-- 记录所有网络请求，为防止存储空间爆炸，默认清理7天前的数据。可修改prune.sql文件以覆盖默认配置。
CREATE TABLE IF NOT EXISTS qiniu_cdnauth_requests (
        id SERIAL PRIMARY KEY,
		request_id int8,
		client_ip VARCHAR(45),
		x_forwarded_for VARCHAR(255),
		user_agent VARCHAR(500),
		http_referer VARCHAR(255),
		request_url varchar(1000) NOT NULL,
		request_headers json NOT NULL,
		raw_url varchar(1000) NOT NULL,
		deleted_at timestamp NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP
    );
CREATE INDEX IF NOT EXISTS "IDX_client_ip" ON qiniu_cdnauth_requests USING btree (client_ip);
CREATE INDEX IF NOT EXISTS "IDX_created_at_client_ip" ON qiniu_cdnauth_requests USING btree (created_at, client_ip);
CREATE INDEX IF NOT EXISTS "IDX_http_referer" ON qiniu_cdnauth_requests USING btree (http_referer);

-- 记录被拦截的网络请求，为防止存储空间爆炸，默认清理30天前的数据。可修改prune.sql文件以覆盖默认配置。
CREATE TABLE IF NOT EXISTS qiniu_cdnauth_block_requests (
        id SERIAL PRIMARY KEY,
		request_id int8,
		client_ip VARCHAR(45),
		x_forwarded_for VARCHAR(255),
		user_agent VARCHAR(500),
		http_referer VARCHAR(255),
		request_url varchar(1000) NOT NULL,
		request_headers json NOT NULL,
		raw_url varchar(1000) NOT NULL,
		block_type SMALLINT NOT NULL DEFAULT 0,
		deleted_at timestamp NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON COLUMN qiniu_cdnauth_block_requests.block_type IS '拦截阻断的理由类别。0=IP黑名单拦截 1=漏洞扫描拦截 2=网络爬虫 3=异常的UserAgent';
CREATE INDEX IF NOT EXISTS "IDX_client_ip_block_requests" ON qiniu_cdnauth_block_requests USING btree (client_ip);
CREATE INDEX IF NOT EXISTS "IDX_block_requests_block_type" ON qiniu_cdnauth_block_requests USING btree (block_type);

-- 记录归档的网络请求。不会被自动删除。
CREATE TABLE IF NOT EXISTS qiniu_cdnauth_archived_requests (
        id SERIAL PRIMARY KEY,
		request_id int8,
		client_ip VARCHAR(45),
		x_forwarded_for VARCHAR(255),
		user_agent VARCHAR(500),
		http_referer VARCHAR(255),
		request_url varchar(1000) NOT NULL,
		request_headers json NOT NULL,
		raw_url varchar(1000) NOT NULL,
		archived_type SMALLINT NOT NULL DEFAULT 0,
		remark VARCHAR(64),
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS "IDX_client_ip_archived_requests" ON qiniu_cdnauth_archived_requests USING btree (client_ip);
-- 添加字段注释
COMMENT ON COLUMN qiniu_cdnauth_archived_requests.archived_type IS '归档类型：0默认，1漏洞扫描，2网络爬虫';

--IP白名单
CREATE TABLE IF NOT EXISTS qiniu_cdnauth_ip_white_list (
        id SERIAL PRIMARY KEY,
		ip VARCHAR(45) NOT NULL,
		title VARCHAR(64) DEFAULT NULL,
		deleted_at timestamp NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);
-- 为ip字段添加唯一约束（此操作会自动创建唯一索引）
CREATE UNIQUE INDEX IF NOT EXISTS "UQE_ip_white_list_ip" ON qiniu_cdnauth_ip_white_list USING btree (ip);

-- IP黑名单
CREATE TABLE IF NOT EXISTS qiniu_cdnauth_ip_black_list (
        id SERIAL PRIMARY KEY,
		ip VARCHAR(45) NOT NULL,
		title VARCHAR(64) DEFAULT NULL,
		black_type SMALLINT NOT NULL DEFAULT 0,
		deleted_at timestamp NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);
-- 为ip字段添加唯一约束（此操作会自动创建唯一索引）
CREATE UNIQUE INDEX IF NOT EXISTS "UQE_ip_black_list_ip" ON qiniu_cdnauth_ip_black_list USING btree (ip);
-- 为black_type字段创建普通索引（非唯一）
CREATE INDEX IF NOT EXISTS "IDX_ip_black_list_black_type" ON qiniu_cdnauth_ip_black_list (black_type);

-- 文件列表
CREATE TABLE IF NOT EXISTS qiniu_cdnauth_files (
	id SERIAL PRIMARY KEY,
	file_key varchar(500) NULL,
	file_size int8 NULL,
	file_hash varchar(64) NULL,
	md5 varchar(32) NULL,
	mime_type varchar(128) NULL,
	file_type SMALLINT NOT NULL DEFAULT 0,
	upload_time timestamp NULL,
	bucket_id SMALLINT NOT NULL DEFAULT 0,
	status SMALLINT NOT NULL DEFAULT 0,
	request_count int8 NOT NULL DEFAULT 0,
	data_raw json NOT NULL
);
-- 添加字段注释
COMMENT ON COLUMN qiniu_cdnauth_files.bucket_id IS '空间名：0bucket123，1bucket567';
COMMENT ON COLUMN qiniu_cdnauth_files.file_type IS '资源的存储类型，0表示标准存储，1 表示低频存储，2 表示归档存储，3 表示深度归档存储，4 表示归档直读存储，5 表示智能分层存储。';
COMMENT ON COLUMN qiniu_cdnauth_files.status IS '文件的存储状态：0启用，1禁用';
COMMENT ON COLUMN qiniu_cdnauth_files.request_count IS '请求次数';
-- 为file_hash字段添加唯一约束（此操作会自动创建唯一索引）
-- CREATE UNIQUE INDEX IF NOT EXISTS "UQE_files_file_hash" ON qiniu_cdnauth_files USING btree (file_hash);
CREATE INDEX IF NOT EXISTS "IDX_files_file_key" ON qiniu_cdnauth_files (file_key);
CREATE INDEX IF NOT EXISTS "IDX_files_request_count" ON qiniu_cdnauth_files (request_count);

-- 同步记录
CREATE TABLE IF NOT EXISTS qiniu_cdnauth_file_sync_log (
	id SERIAL PRIMARY KEY,
	bucket_id SMALLINT NOT NULL DEFAULT 0,
	has_next boolean NOT NULL DEFAULT false,
	cursor_marker varchar(255) NULL,
	size_len SMALLINT NOT NULL DEFAULT 0,
	sort int8 NOT NULL DEFAULT 0,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP	 
);

-- statis 数据统计
CREATE TABLE IF NOT EXISTS qiniu_cdnauth_statis (
	id SERIAL PRIMARY KEY,
	bucket_id SMALLINT NOT NULL DEFAULT 0,
	statis_date DATE NOT NULL DEFAULT CURRENT_DATE,
	request_count int8 NOT NULL DEFAULT 0,
	blocked_count int8 NOT NULL DEFAULT 0,
	blocked_black_count int8 NOT NULL DEFAULT 0,
	blocked_scanvul_count int8 NOT NULL DEFAULT 0,
	blocked_webspider_count int8 NOT NULL DEFAULT 0,
	blocked_useragent_count int8 NOT NULL DEFAULT 0,
	request_size int8 NOT NULL DEFAULT 0,
	blocked_size int8 NOT NULL DEFAULT 0,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON COLUMN qiniu_cdnauth_statis.bucket_id IS '存储空间ID';
COMMENT ON COLUMN qiniu_cdnauth_statis.statis_date IS '统计日期';
COMMENT ON COLUMN qiniu_cdnauth_statis.request_count IS '请求次数';
COMMENT ON COLUMN qiniu_cdnauth_statis.blocked_count IS '拦截请求次数';
COMMENT ON COLUMN qiniu_cdnauth_statis.blocked_black_count IS '黑名单拦截请求次数';
COMMENT ON COLUMN qiniu_cdnauth_statis.blocked_scanvul_count IS '漏洞扫描拦截请求次数';
COMMENT ON COLUMN qiniu_cdnauth_statis.blocked_webspider_count IS '网络爬虫拦截请求次数';
COMMENT ON COLUMN qiniu_cdnauth_statis.blocked_useragent_count IS '异常用户代理拦截请求次数';
COMMENT ON COLUMN qiniu_cdnauth_statis.request_size IS '请求消耗流量大小';
COMMENT ON COLUMN qiniu_cdnauth_statis.blocked_size IS '拦截流量大小';
CREATE INDEX IF NOT EXISTS "IDX_statis_bucket_id" ON qiniu_cdnauth_statis (bucket_id);
CREATE INDEX IF NOT EXISTS "IDX_statis_statis_date" ON qiniu_cdnauth_statis (statis_date);


-- 要迁移的文件列表
CREATE TABLE IF NOT EXISTS qiniu_cdnauth_file_migrate_list (
    id SERIAL PRIMARY KEY,
	file_url varchar(1000) DEFAULT NULL,
    file_key VARCHAR(500) NOT NULL,
	status SMALLINT NOT NULL DEFAULT 0,
    from_table VARCHAR(64) DEFAULT NULL,
    from_column_name VARCHAR(64) DEFAULT NULL,
	from_column_value TEXT,
	from_bucket VARCHAR(32),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS "IDX_file_migrate_list_file_key" ON qiniu_cdnauth_file_migrate_list (file_key);
COMMENT ON COLUMN qiniu_cdnauth_file_migrate_list.file_url IS '文件原始URL';
COMMENT ON COLUMN qiniu_cdnauth_file_migrate_list.status IS '迁移状态：-1操作失败0未开始，1copy成功，2move成功，3原文件已删除';
COMMENT ON COLUMN qiniu_cdnauth_file_migrate_list.from_table IS '来源表名';
COMMENT ON COLUMN qiniu_cdnauth_file_migrate_list.from_column_name IS '来源字段名';
COMMENT ON COLUMN qiniu_cdnauth_file_migrate_list.from_column_value IS '来源字段值';
COMMENT ON COLUMN qiniu_cdnauth_file_migrate_list.file_key IS '文件存储的key';
COMMENT ON COLUMN qiniu_cdnauth_file_migrate_list.from_bucket IS '来源bucket文件空间';


-- 文件操作日志表
CREATE TABLE IF NOT EXISTS qiniu_cdnauth_file_opt_log (
    id SERIAL PRIMARY KEY,
    file_key VARCHAR(500) NOT NULL,
    opt_type SMALLINT NOT NULL,
	state boolean NOT NULL DEFAULT false,
    file_size int8 DEFAULT NULL,
    upload_time TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    qiniu_etag varchar(64) DEFAULT NULL,
    md5 VARCHAR(32) DEFAULT NULL,
    from_bucket VARCHAR(32),
    to_bucket VARCHAR(32)
);
CREATE INDEX IF NOT EXISTS "IDX_file_opt_log_file_key" ON qiniu_cdnauth_file_opt_log (file_key);
COMMENT ON COLUMN qiniu_cdnauth_file_opt_log.opt_type IS '操作类型：1copy, 2move, 3delete';
COMMENT ON COLUMN qiniu_cdnauth_file_opt_log.state IS '操作状态：1成功|0失败';
COMMENT ON COLUMN qiniu_cdnauth_file_opt_log.file_size IS '文件大小';
COMMENT ON COLUMN qiniu_cdnauth_file_opt_log.upload_time IS '上传时间';
COMMENT ON COLUMN qiniu_cdnauth_file_opt_log.qiniu_etag IS '七牛返回的ETag';
COMMENT ON COLUMN qiniu_cdnauth_file_opt_log.md5 IS '文件的MD5';
COMMENT ON COLUMN qiniu_cdnauth_file_opt_log.from_bucket IS '来源bucket文件空间';
COMMENT ON COLUMN qiniu_cdnauth_file_opt_log.to_bucket IS '目标bucket文件空间';
