package main

import (
	"context"
	"fmt"
	"log"

	"hata/internal/orm"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	client, err := orm.Open("sqlite3", "file:data/dev.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Query all sys__kv records
	kvs, err := client.SysKV.Query().All(ctx)
	if err != nil {
		log.Fatalf("failed querying sys__kv: %v", err)
	}

	fmt.Println("sys__kv records:")
	for _, kv := range kvs {
		fmt.Printf("  %s = %s\n", kv.Key, kv.Value)
	}
}
