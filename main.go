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
	BuildVersion         string
	clusterName          string
	servicesFilter       *regexp.Regexp
	servicesFilterString string
	verbose, version     bool
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
	flag.StringVar(&servicesFilterString, "servicesFilter", `.*`, "Service Filter (RegExp pattern)")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&version, "version", false, "display version and exit")

	flag.Usage = usage
	flag.Parse()

	if version {
		log.Printf("version: %s\n", BuildVersion)
		os.Exit(0)
	}

	if verbose {
		log.Printf("BuildVersion: %s\n", BuildVersion)
		log.Printf("clusterName: %s\n", clusterName)
		log.Printf("servicesFilter: %s\n", servicesFilter)
	}

}

func listServices(sess *session.Session, filter *regexp.Regexp) ([]string, error) {
	var services []string

	svc := ecs.New(sess)
	listServicesInput := &ecs.ListServicesInput{
		Cluster: aws.String(clusterName),
	}

	for {
		listServicesOutput, err := svc.ListServices(listServicesInput)
		if err != nil {
			return nil, err
		}
		for _, serviceArn := range aws.StringValueSlice(listServicesOutput.ServiceArns) {
			if servicesFilter.MatchString(serviceArn) {
				services = append(services, serviceArn)
			}
		}

		if listServicesOutput.NextToken != nil {
			listServicesInput.NextToken = listServicesOutput.NextToken
			if verbose {
				log.Printf("listServices NextToken found, will make another request\n")
			}
		} else {
			if verbose {
				log.Printf("listServices NextToken not found, breaking request loop\n")
			}
			break
		}
	}
	return services, nil
}

func main() {
	initFlags()

	servicesFilter = regexp.MustCompile(servicesFilterString)

	sess, err := session.NewSession()
	if err != nil {
		log.Panicf("%s\n", err)
	}

	clusterServices, err := listServices(sess, servicesFilter)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Fatalf("%s\n", aerr.Error())
		} else {
			log.Fatalf("%s\n", err.Error())
		}
	}

	for _, serviceArn := range clusterServices {
		log.Printf("%s\n", serviceArn)
	}
}
