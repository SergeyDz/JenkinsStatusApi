package main


import "time"

type Build struct {
	Id        			int       `json:"id"`
	JobName      		string    `json:"job" xorm:"JobName"`
	JobNumber      		string    `json:"jobid" xorm:"JobNumber"`
	RepositoryName      string    `json:"repository" xorm:"RepositoryName"`
	Status				string 	  `json:"status" xorm:"Status"`
	CreatedOn       	time.Time `json:"createdon" xorm:"CreatedOn"`
}

type Builds []Build

