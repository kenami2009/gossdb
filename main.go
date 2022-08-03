package main

import (
	"fmt"
	"go-ssdb/ssdb"
)

func main() {
	addr := "127.0.0.1:8888"
	conn, err := ssdb.Connect(addr)
	if err != nil {
		fmt.Println(err)
	}
	conn.Do("set", "key", "value")
	replay, _ := conn.Do("get", "key")
	val, _ := ssdb.String(replay, nil)
	fmt.Println(val)

}
