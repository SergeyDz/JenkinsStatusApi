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

	//engine.ShowSQL(true)
}

func RepoFindBuild(id int) Build {
	return Build{}
}

//this is bad, I don't think it passes race condtions
func RepoCreateBuild(t Build) Build {
	if _, err := engine.Insert(&t); err != nil{
		panic(err)
	}
	var repository Repository

	has, err := engine.Where("Name = ?", t.RepositoryName).Desc("Id").Get(&repository)

	if(err != nil){
		panic(err)
	}

	if(has)	{
		engine.Update(&Repository{ Name:t.RepositoryName, Url:t.RepositoryUrl, Status: t.Status }, &Repository{Name:t.RepositoryName})
	} else {
		engine.Insert(&Repository{ Name:t.RepositoryName, Url:t.RepositoryUrl, Status: t.Status })
	}

	return t
}

func RepoShowAllBuilds(size int) Builds{
	builds := make([]Build, 0)
	var err = engine.Desc("Id").Limit(size).Find(&builds)
	if err != nil {
		fmt.Println(err)
	}
	return builds
}

func RepoShowAllRepos() Repositories{
	repos := make([]Repository, 0)
	var err = engine.Asc("Status").Asc("Name").Limit(500).Find(&repos)
	if err != nil {
		fmt.Println(err)
	}

	return repos
}

