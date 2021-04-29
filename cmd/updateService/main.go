package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	// "regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var (
	//BuildVersion is set during `go build` by `VERSION.txt`
	BuildVersion     string
	clusterName      string
	serviceArn       string
	confirmUpdate    bool
	verbose, version bool
)

func usage() {
	println(`Usage: updateService [options]
Work with AWS ECS
Options:`)
	flag.PrintDefaults()
	println(`For more information, see https://github.com/jknutson/one-wire-temp-go`)
}

func initFlags() {
	flag.StringVar(&clusterName, "clusterName", "default", "ECS Cluster name")
	flag.StringVar(&serviceArn, "serviceArn", "", "ECS Service ARN; '-' reads from STDIN (required)")
	flag.BoolVar(&confirmUpdate, "confirmUpdate", false, "confirm update")
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
		log.Printf("serviceArn: %s\n", serviceArn)
		log.Printf("confirmUpdate: %t\n", confirmUpdate)
	}

	// TODO: exit/error if serviceArn == ""
}

func updateService(sess *session.Session, serviceArn string) error {
	if confirmUpdate {
		// TODO: validate ARN regex
		if verbose {
			log.Printf("updating %s\n", serviceArn)
		}

		svc := ecs.New(sess)
		updateServiceInput := &ecs.UpdateServiceInput{
			Cluster:            aws.String(clusterName),
			ForceNewDeployment: aws.Bool(true),
			Service:            aws.String(serviceArn),
		}

		_, err := svc.UpdateService(updateServiceInput)
		if err != nil {
			return err
		}
		log.Printf("updated %s\n", serviceArn)
		return nil
	}
	if verbose {
		log.Printf("skipping update: confirmUpdate = %t\n", confirmUpdate)
	}
	return nil
}

// TODO: move to pkg/ and share with other cmds
func awsError(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		log.Fatalf("%s\n", aerr.Error())
	} else {
		log.Fatalf("%s\n", err.Error())
	}
}

func main() {
	initFlags()

	sess, err := session.NewSession()
	if err != nil {
		log.Panicf("%s\n", err)
	}

	if serviceArn == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			err := updateService(sess, scanner.Text())
			if err != nil {
				awsError(err)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	} else {
		err := updateService(sess, serviceArn)
		if err != nil {
			awsError(err)
		}
	}
}
