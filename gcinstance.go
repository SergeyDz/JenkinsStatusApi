package main

type Instance struct {
	ID 	string	`json:"ID"`
	NAME	string	`json:"NAME"`
	ZONE	string	`json:"ZONE" xorm:"ZONE"`
	MACHINE_TYPE	string	`json:"MACHINE_TYPE" xorm:"MACHINE_TYPE"`
	INTERNAL_IP	string	`json:"INTERNAL_IP" xorm:"INTERNAL_IP"`
	EXTERNAL_IP	string	`json:"EXTERNAL_IP" xorm:"EXTERNAL_IP"`
	STATUS	string	`json:"STATUS" xorm:"STATUS"`
}

type Instances []Repository

