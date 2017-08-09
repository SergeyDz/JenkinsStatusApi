package main


import "time"

type Build struct {
	Id        			int       `json:"id"`
	JobName      		string    `json:"job" xorm:"JobName"`
	JobNumber      		string    `json:"jobid" xorm:"JobNumber"`
	JobDescription		string    `json:"jobdescription" xorm:"JobDescription"`
	JobUrl 				string    `json:"joburl" xorm:"JobUrl"`
	BuildDuration 		string    `json:"duration" xorm:"BuildDuration"`
	BranchName			string    `json:"branchname" xorm:"BranchName"`
	RepositoryName      string    `json:"repositoryname" xorm:"RepositoryName"`
	RepositoryUrl       string    `json:"repositoryurl" xorm:"RepositoryName"`
	Status				string 	  `json:"status" xorm:"Status"`
	CreatedOn       	time.Time `json:"createdon" xorm:"CreatedOn"`
	TriggedBy			string 	  `json:"triggedby" xorm:"TriggedBy"`
}

type Builds []Build

