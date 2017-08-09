package main


import "time"

type Repository struct {
	Id        			int       `json:"id"`
	Name      			string    `json:"name" xorm:"Name"`
	Url      			string    `json:"url" xorm:"Url"`
	Status				string    `json:"status" xorm:"Status"`
	CreatedOn       	time.Time `json:"createdon" xorm:"CreatedOn"`
	UpdatedOn       	time.Time `json:"updatedon" xorm:"UpdatedOn"`
}

type Repositories []Repository

