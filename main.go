package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cli/go-gh"
	graphql "github.com/cli/shurcooL-graphql"
)

func main() {
	var organization string
	var affiliation string
	var samlIdentities bool

	validAffiliation := []string{"ALL", "OUTSIDE", "DIRECT"}

	flag.StringVar(&organization, "organization", "octodemo", "organization")
	flag.StringVar(&affiliation, "affiliation", "direCT", "affiliation")
	flag.BoolVar(&samlIdentities, "saml-identities", false, "saml-identities")

	flag.Parse()

	if !contains(validAffiliation, affiliation) {
		fmt.Printf("Affiliation value %s is invalid. Valid options for affiliation are All, Direct, and Outside.", affiliation)
		os.Exit(0)
	}

	fmt.Println("Organization: ", organization)
	fmt.Println("Affiliation: ", affiliation)
	fmt.Println("SAML Identities: ", samlIdentities)

	OrgMembership(organization, affiliation, samlIdentities)
	// createFile()
}

// a function that creates a file in the current directory
func createFile() {

	datetime := getDateTime()

	f, err := os.Create(datetime + "_audit.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString("test")
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// a function that gets the current date and time
func getDateTime() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

// a function that is a graphql query to return all the users in an organization
func OrgMembership(organization string, affiliation string, samlIdentities bool) {
	client, err := gh.GQLClient(nil)
	if err != nil {
		log.Fatal(err)
	}
	var query struct {
		Organization struct {
			MembersWithRole struct {
				Edges []struct {
					Node struct {
						Id    string
						Login string
					}
				}
			} `graphql:"membersWithRole(first: $first)"`
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"first": graphql.Int(10),
		"login": graphql.String(organization),
	}

	err = client.Query("OrgMembers", &query, variables)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(query)
}

// a function to check if a value is part of an array
func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}
