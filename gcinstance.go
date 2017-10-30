package main

type Instance struct {
	ID           string `json:"ID"`
	NAME         string `json:"NAME"`
	ZONE         string `json:"ZONE" xorm:"ZONE"`
	MACHINE_TYPE string `json:"MACHINE_TYPE" xorm:"MACHINE_TYPE"`
	INTERNAL_IP  string `json:"INTERNAL_IP" xorm:"INTERNAL_IP"`
	EXTERNAL_IP  string `json:"EXTERNAL_IP" xorm:"EXTERNAL_IP"`
	STATUS       string `json:"STATUS" xorm:"STATUS"`
	JENBUILD     string `json:"JENBUILD" xorm:"JENBUILD"`
}

// https://mholt.github.io/json-to-go/ to get JSON struckture

type JenkinsBuilds struct {
	Builds []struct {
		Actions []struct {
			Parameters []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"parameters,omitempty"`
		} `json:"actions"`
		Url      string `json:"url"`
		Building bool   `json:"building"`
	} `json:"builds"`
}

type JenkinsJob struct {
	Timestamp   string `json:"Timestamp"`
	Building    bool `json:"Building"`
	Result      string `json:"Result" xorm:"RESULT"`
	DisplayName string `json:"DisplayName" xorm:"DISPLAYNAME"`
	URL         string `json:"URL" xorm:"URL"`
	UserName    string `json:"UserName" xorm:"USERNAME"`
}

type AppTerraformPOC struct {
	Builds []struct {
		Actions []struct {
			Class  string `json:"_class,omitempty"`
			Causes []struct {
				Class         string `json:"_class"`
				UpstreamBuild int    `json:"upstreamBuild"`
				UpstreamURL   string `json:"upstreamUrl"`
			} `json:"causes,omitempty"`
		} `json:"actions"`
		Building    bool        `json:"building"`
		DisplayName string      `json:"displayName"`
		Result      string `json:"result"`
		Timestamp   int64       `json:"timestamp"`
		URL         string      `json:"url"`
	} `json:"builds"`
}

type BuilsArtifacrtPOP struct {
	Actions []struct {
		Class  string `json:"_class,omitempty"`
		Causes []struct {
			Class         string `json:"_class"`
			UpstreamBuild int    `json:"upstreamBuild"`
			UpstreamURL   string `json:"upstreamUrl"`
		} `json:"causes,omitempty"`
	} `json:"actions"`
}

type CloudEnvPOC struct {
	Actions []struct {
		Class  string `json:"_class,omitempty"`
		Causes []struct {
			Class    string `json:"_class"`
			UserName string `json:"userName"`
		} `json:"causes,omitempty"`
	} `json:"actions"`
}

type Instances []Repository
