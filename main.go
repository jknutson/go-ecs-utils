package main

import (
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var (
	ecsCluster          string = "tst-novu-apps"
	serviceFilter       *regexp.Regexp
	serviceFilterString string = `.*qa7.*` // TODO: get from flag
	verbose             bool   = true      // TODO: get from flag
)

func main() {
	var clusterServices []string

	// TODO: get string from flag
	serviceFilter = regexp.MustCompile(serviceFilterString)

	sess, err := session.NewSession()
	if err != nil {
		log.Panicf("%s\n", err)
	}
	svc := ecs.New(sess)
	listServicesInput := &ecs.ListServicesInput{
		Cluster: aws.String(ecsCluster),
	}

	for {
		listServicesOutput, err := svc.ListServices(listServicesInput)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				log.Println(aerr.Error())
			} else {
				log.Println(err.Error())
			}
		}
		for _, serviceArn := range aws.StringValueSlice(listServicesOutput.ServiceArns) {
			if serviceFilter.MatchString(serviceArn) {
				clusterServices = append(clusterServices, serviceArn)
			}
		}

		if listServicesOutput.NextToken != nil {
			listServicesInput.NextToken = listServicesOutput.NextToken
			if verbose {
				log.Printf("NextToken found, will make another request\n")
			}
		} else {
			if verbose {
				log.Printf("NextToken NOT found, breaking request loop\n")
			}
			break
		}
	}

	for _, serviceArn := range clusterServices {
		log.Printf("%s\n", serviceArn)
	}
}
