package main

import (
	"flag"
	"log"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var (
	//BuildVersion is set during `go build` by `VERSION.txt`
	BuildVersion        string
	clusterName         string
	serviceFilter       *regexp.Regexp
	serviceFilterString string
	verbose, version    bool
)

func usage() {
	println(`Usage: ecs-utils [options]
Work with AWS ECS
Options:`)
	flag.PrintDefaults()
	println(`For more information, see https://github.com/jknutson/one-wire-temp-go`)
}

func initFlags() {
	flag.StringVar(&clusterName, "clusterName", "default", "ECS cluster name")
	flag.StringVar(&serviceFilterString, "serviceFilter", `.*`, "Service Filter (RegExp pattern)")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&version, "version", false, "display version and exit")

	flag.Usage = usage
	flag.Parse()

	if version {
		log.Printf("version: %s\n", BuildVersion)
		os.Exit(0)
	}
}

func main() {
	initFlags()
	var clusterServices []string

	serviceFilter = regexp.MustCompile(serviceFilterString)

	sess, err := session.NewSession()
	if err != nil {
		log.Panicf("%s\n", err)
	}
	svc := ecs.New(sess)
	listServicesInput := &ecs.ListServicesInput{
		Cluster: aws.String(clusterName),
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
