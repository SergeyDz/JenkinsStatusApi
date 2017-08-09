package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"fmt"
	"os"
)

var engine *xorm.Engine

// Give us some seed data
func init() {
	var err error
	var dataSource = os.Getenv("MYSQL_CONNECTION")
	engine, err = xorm.NewEngine("mysql", dataSource)
	if(err != nil){
		panic(err)
	}

	engine.ShowSQL(true)
}

func RepoFindBuild(id int) Build {
	return Build{}
}

//this is bad, I don't think it passes race condtions
func RepoCreateBuild(t Build) Build {
	if _, err := engine.Insert(&t); err != nil{
		panic(err)
	}

	return t
}

func RepoShowAllBuilds() Builds{
	builds := make([]Build, 0)
	var err = engine.Desc("Id").Limit(500).Find(&builds)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(builds)
	return builds
}

