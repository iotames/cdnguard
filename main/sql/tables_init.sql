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
		deleted_at timestamp NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP
)