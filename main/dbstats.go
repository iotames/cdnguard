package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func showDbStats() bool {
	requrl := fmt.Sprintf("http://127.0.0.1:%d/api/local/dbstats", WebPort)
	req, err := http.NewRequest("GET", requrl, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	type StatsResult struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data sql.DBStats `json:"data"`
	}
	var result StatsResult
	stats := sql.DBStats{}
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		panic(err)
	}
	if result.Code != 200 {
		log.Println("请求失败:", result.Msg)
		return true
	}
	stats = result.Data
	// stats := gdb.Stats()
	log.Printf("数据库连接统计:\n")
	log.Printf("最大打开连接数: %d\n", stats.MaxOpenConnections)
	log.Printf("打开连接数: %d\n", stats.OpenConnections)
	log.Printf("使用中连接数: %d\n", stats.InUse)
	log.Printf("空闲连接数: %d\n", stats.Idle)
	log.Printf("等待新连接的数量: %d\n", stats.WaitCount)
	log.Printf("因超时关闭的连接数: %d\n", stats.MaxLifetimeClosed)
	log.Printf("----rawjson(%s)---------\n", buf.String())
	return true
}
