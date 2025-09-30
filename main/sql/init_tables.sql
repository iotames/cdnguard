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
CREATE INDEX IF NOT EXISTS "IDX_client_ip_block_requests" ON qiniu_cdnauth_block_requests USING btree (client_ip);
CREATE INDEX IF NOT EXISTS "IDX_block_requests_block_type" ON qiniu_cdnauth_block_requests USING btree (block_type);

----------- ip_white_list, ip_black_list
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
