package main

type Instance struct {
	ID 	string	`json:"ID"`
	NAME	string	`json:"NAME"`
	ZONE	string	`json:"ZONE" xorm:"ZONE"`
	MACHINE_TYPE	string	`json:"MACHINE_TYPE" xorm:"MACHINE_TYPE"`
	INTERNAL_IP	string	`json:"INTERNAL_IP" xorm:"INTERNAL_IP"`
	EXTERNAL_IP	string	`json:"EXTERNAL_IP" xorm:"EXTERNAL_IP"`
	STATUS	string	`json:"STATUS" xorm:"STATUS"`
	JENSTAT	string	`json:"JENSTAT" xorm:"JENSTAT"`
}

type JenkinsBuilds struct {
	Builds []struct {
		Actions []struct {
			Parameters []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"parameters,omitempty"`
		} `json:"actions"`
		Building        bool        `json:"building"`
		FullDisplayName string      `json:"fullDisplayName"`
		ID              string      `json:"id"`
		Result          interface{} `json:"result"`
	} `json:"builds"`
}

type Instances []Repository

