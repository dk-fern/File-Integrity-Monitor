package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

func main() {
	// Create .env file in same directory or use an environmental variable for api key
	godotenv.Load(".env")
	apiKey := os.Getenv("apiKey")

	//---------FLAGS---------//
	ipAddr := flag.String("ip", "", "ip address to lookup")
	domain := flag.String("domain", "", "domain address to lookup")

	flag.Parse()
	//-----------------------//

	// Check if -ip flag is used
	if *ipAddr != "" {
		ipPattern := `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
		match, _ := regexp.MatchString(ipPattern, *ipAddr)
		if !match {
			log.Fatal("\nip address is not correctly formatted or is out of a valid ip range: ", *ipAddr)
		}

		// Pull ip report data
		ipReport, err := getIPReport(apiKey, *ipAddr)
		if err != nil {
			log.Fatal("error pulling ip report: ", err)
		}

		// Print ip report to the console
		fmt.Println("Results for:", *ipAddr)
		fmt.Println("~~~Overall~~~")
		fmt.Println("Malicious:", ipReport.Data.Attributes.LastAnalysisStats.Malicious)
		fmt.Println("Suspicious:", ipReport.Data.Attributes.LastAnalysisStats.Suspicious)
		fmt.Println("Undetected:", ipReport.Data.Attributes.LastAnalysisStats.Undetected)
		fmt.Println("Harmless:", ipReport.Data.Attributes.LastAnalysisStats.Harmless)
		fmt.Println("Timeout:", ipReport.Data.Attributes.LastAnalysisStats.Timeout)

		fmt.Println("\n~~~Community Votes~~~")
		fmt.Println("Harmless:", ipReport.Data.Attributes.TotalVotes.Harmless)
		fmt.Println("Malicious:", ipReport.Data.Attributes.TotalVotes.Malicious)

		fmt.Println("\nWhoIs Results:")
		fmt.Println(ipReport.Data.Attributes.WhoIs)

		fmt.Println("Access full report here:", ipReport.Data.Links.Self)
	}

	// Check if -domain flag is used
	if *domain != "" {
		domainReport, err := getDomainReport(apiKey, *domain)
		if err != nil {
			log.Fatal("error pulling domain report: ", err)
		}

		// Print domain report to the console
		fmt.Println("Results for:", *domain)
		fmt.Println("~~~Overall~~~")
		fmt.Println("Malicious:", domainReport.Data.Attributes.LastAnalysisStats.Malicious)
		fmt.Println("Suspicious:", domainReport.Data.Attributes.LastAnalysisStats.Malicious)
		fmt.Println("Undetected:", domainReport.Data.Attributes.LastAnalysisStats.Malicious)
		fmt.Println("Harmless:", domainReport.Data.Attributes.LastAnalysisStats.Malicious)
		fmt.Println("Timeout:", domainReport.Data.Attributes.LastAnalysisStats.Malicious)

		fmt.Println("\n~~~Community Votes~~~")
		fmt.Println("Harmless:", domainReport.Data.Attributes.TotalVotes.Harmless)
		fmt.Println("Malicious:", domainReport.Data.Attributes.TotalVotes.Malicious)

		if domainReport.Data.Attributes.Categories.AlphaMountainAi != "" || domainReport.Data.Attributes.Categories.DrWeb != "" {
			fmt.Println("\n~~~Reported Categories~~~")
			fmt.Println("AlphaMountainAi:", domainReport.Data.Attributes.Categories.AlphaMountainAi)
			fmt.Println("Dr Web:", domainReport.Data.Attributes.Categories.DrWeb)
			fmt.Println("Sophos:", domainReport.Data.Attributes.Categories.Sophos)
			fmt.Println("Webroot:", domainReport.Data.Attributes.Categories.Webroot)
			fmt.Println("Forcepoint ThreatSeeker:", domainReport.Data.Attributes.Categories.ForcepointThreatSeeker)
		}

		if domainReport.Data.Attributes.CrowdsourcedContext != nil {
			fmt.Println("\n~~~Crowdourced Info~~~")
			for _, context := range domainReport.Data.Attributes.CrowdsourcedContext {
				fmt.Println("Severity:", context.Severity)
				fmt.Println("Title:", context.Title)
				fmt.Println("Details:", context.Details)
				fmt.Println("Source:", context.Source)
			}
		}

		fmt.Println("\nAccess full report here:", domainReport.Data.Links.Self)
	}

}
