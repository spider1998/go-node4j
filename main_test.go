package main

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"testing"
)

func helloWorld(uri, username, password string) (string, error) {
	var (
		err      error
		driver   neo4j.Driver
		session  neo4j.Session
		result   neo4j.Result
		greeting interface{}
	)

	// 创建neo4j驱动
	driver, err = neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return "", err
	}
	defer driver.Close()

	// 获取neo4j session
	session, err = driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return "", err
	}
	defer session.Close()

	// 通过事务操作neo4j
	greeting, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err = transaction.Run(
			"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
			map[string]interface{}{"message": "hello, world"})
		if err != nil {
			return nil, err
		}

		// 获取neo4j操作结果
		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return "", err
	}

	return greeting.(string), nil
}

func TestHelloWorld(t *testing.T) {
	str, err :=	helloWorld("bolt://localhost:7687", "neo4j", "123456")

	if err != nil {
		panic(err)
	}

	fmt.Println(str)
}

//https://github.com/neo4j-drivers/seabolt