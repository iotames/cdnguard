package migrate

// 使用事务保证操作的原子性
// https://developer.qiniu.com/kodo/1250/batch
func CopyFiles() {

}
func MoveFiles() {

}

func DeleteFiles() {

}

// file_migrate_list
//
// id SERIAL PRIMARY KEY,
// file_url varchar(1000) NOT NULL,
// status SMALLINT DEFAULT NULL,
// from_table VARCHAR(64) NOT NULL,
// from_column VARCHAR(64) NOT NULL,
// file_key VARCHAR(500) NOT NULL,
// from_bucket VARCHAR(32),
// created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// updated_at timestamp DEFAULT CURRENT_TIMESTAMP

// file_opt_log
//
// id SERIAL PRIMARY KEY,
// file_key VARCHAR(500) NOT NULL,
// opt_type SMALLINT NOT NULL,
// 	state boolean NOT NULL DEFAULT false,
// file_size int8 NULL,
// upload_time TIMESTAMP,
// created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// qiniu_etag varchar(64) NOT NULL,
// md5 VARCHAR(32) NOT NULL,
// from_bucket VARCHAR(32),
// to_bucket VARCHAR(32)
