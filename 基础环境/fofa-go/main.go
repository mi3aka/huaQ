package main

import (
	"fofa/fofa"
	"fmt"
)

func main() {
	email := ""
	key := ""
	fofaClient := fofa.NewFofaClient([]byte(email), []byte(key))
	if fofaClient == nil {
		fmt.Println("create fofa client")
		return
	}
	fields := []byte("domain,host,ip,port,protocol")
	targets := []string{`server=="nginx"`}
	for index, target := range targets {
		result, err := fofaClient.Query([]byte(target), fields, 1)
		if err != nil {
			fmt.Println(err.Error())
		}
		fofaClient.WriteFile(result, fields, fmt.Sprintf("result/%d.json", index))
	}
}