INSERT INTO qiniu_cdnauth_ip_black_list (ip, title, black_type) 
VALUES ($1, $2, $3) 
ON CONFLICT (ip) 
DO NOTHING;