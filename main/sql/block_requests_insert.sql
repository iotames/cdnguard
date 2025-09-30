INSERT INTO qiniu_cdnauth_block_requests (
request_id,client_ip,x_forwarded_for,user_agent,http_referer,request_url,request_headers,raw_url,block_type
)VALUES ($1,$2,$3,$4,$5,$6,$7, $8, $9);