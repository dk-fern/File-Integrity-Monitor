package main

type IP struct {
	Data Data `json:"data"`
}

type Data struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Links      Links      `json:"links"`
	Attributes Attributes `json:"attributes"`
}

type Links struct {
	Self string `json:"self"`
}

type Attributes struct {
	Country             string                `json:"country"`
	LastAnalysisStats   LastAnalysisStats     `json:"last_analysis_stats"`
	WhoIs               string                `json:"whois"`
	Categories          Categories            `json:"categories"`
	CrowdsourcedContext []CrowdsourcedContext `json:"crowdsourced_context"`
	TotalVotes          TotalVotes            `json:"total_votes"`
}

type LastAnalysisStats struct {
	Malicious  int `json:"malicious"`
	Suspicious int `json:"suspicious"`
	Undetected int `json:"undetected"`
	Harmless   int `json:"harmless"`
	Timeout    int `json:"timeout"`
}

type Domain struct {
	Data Data `json:"data"`
}

type Categories struct {
	AlphaMountainAi        string `json:"alphaMountain.ai"`
	DrWeb                  string `json:"Dr.Web"`
	Sophos                 string `json:"Sophos"`
	Webroot                string `json:"Webroot"`
	ForcepointThreatSeeker string `json:"Forecepoint ThreatSeeker"`
}

type CrowdsourcedContext struct {
	Severity  string `json:"severity"`
	Timestamp int64  `json:"timestamp"`
	Details   string `json:"details"`
	Title     string `json:"title"`
	Source    string `json:"source"`
}

type TotalVotes struct {
	Harmless  int `json:"harmless"`
	Malicious int `json:"malicious"`
}
