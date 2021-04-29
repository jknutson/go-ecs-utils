package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	// "regexp"
	/*
		"github.com/aws/aws-sdk-go/aws"
		"github.com/aws/aws-sdk-go/aws/awserr"
		"github.com/aws/aws-sdk-go/aws/session"
		"github.com/aws/aws-sdk-go/service/ecs"
	*/)

var (
	//BuildVersion is set during `go build` by `VERSION.txt`
	BuildVersion     string
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
	flag.StringVar(&serviceArn, "serviceArn", "", "ECS Service ARN. '-' reads from STDIN (required")
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
		log.Printf("serviceArn: %s\n", serviceArn)
		log.Printf("confirmUpdate: %t\n", confirmUpdate)
	}

	// TODO: exit/error if serviceArn == ""
}

func updateService(serviceArn string) error {
	// TODO: validate ARN regex
	if verbose {
		log.Printf("%s\n", serviceArn)
	}
	return nil
}

func main() {
	initFlags()

	/*
		sess, err := session.NewSession()
		if err != nil {
			log.Panicf("%s\n", err)
		}
	*/

	if serviceArn == "-" { // read from STDIN
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			err := updateService(scanner.Text())
			if err != nil {
				log.Fatalf("%s\n", err)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	} else {
		err := updateService(serviceArn)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
	}
}
