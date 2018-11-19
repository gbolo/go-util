package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/viper"
)

const (
	// defines environment states
	// means that ALL instances for an env are in "running" state
	ENV_RUNNING = "running"
	// means that AT LEAST ONE instance for an env is NOT in a "running" state
	ENV_DOWN = "stopped"
	// means that AT LEAST ONE instance for an env is NOT in a "running" state or "stopped" state
	ENV_CHANGING = "changing-state"
)

var (
	// global aws client
	awsClient *ec2.EC2
	// global cached env list
	cachedTable envList
	// aws region
	awsRegion string
	// aws tags
	requiredTagKey, requiredTagValue, environmentTagKey string
	// safety, will refuse to shutdown if more than this amount of instances is requested
	maxInstancesToShutdown int
	// ignore these instance types
	instanceTypeIgnore []string
)

type ec2Instance struct {
	Id          string
	Type        string
	Name        string
	State       string
	Environment string
}

type ec2Environment struct {
	Name      string
	Instances []ec2Instance
	CPUs      int
	MemoryMb  int
	State     string
}

type ec2EnvironmentSummary struct {
	Name             string
	RunningInstances int
	TotalInstances   int
	State            string
}

// for global cached table
type envList []ec2Environment

// determineEnvStates
func determineEnvStates() {
	for i, env := range cachedTable {
		cachedTable[i].State = ENV_RUNNING
		for _, instance := range env.Instances {
			if instance.State != "running" && instance.State != "stopped" {
				cachedTable[i].State = ENV_CHANGING
				break
			} else if instance.State == "stopped" {
				cachedTable[i].State = ENV_DOWN
			}
		}
	}
}

// checks if an instance should be included
// true if its OK, false to ignore
func checkInstanceType(instanceType string) (ok bool) {
	ok = true
	for _, ignoredType := range instanceTypeIgnore {
		if ignoredType == instanceType {
			ok = false
			break
		}
	}
	return
}

// adds and instance to cachedTable
func addInstance(instance *ec2Instance) {
	if !checkInstanceType(instance.Type) {
		fmt.Printf("instance is being ignored: %s\n", instance.Name)
		return
	}
	envExists := false
	for i, env := range cachedTable {
		if env.Name == instance.Environment {
			envExists = true
			cachedTable[i].Instances = append(cachedTable[i].Instances, *instance)
		}
	}
	if !envExists {
		ec2env := ec2Environment{
			Name:      instance.Environment,
			Instances: []ec2Instance{*instance},
		}
		cachedTable = append(cachedTable, ec2env)
	}
}

// polls aws for updates to cachedTable
func refreshTable() (err error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []ec2.Filter{
			{
				Name: aws.String("tag:" + requiredTagKey),
				Values: []string{
					requiredTagValue,
				},
			},
		},
	}

	req := awsClient.DescribeInstancesRequest(params)
	resp, err := req.Send()
	if err != nil {
		log.Errorf("failed to describe instances, %s, %v", awsRegion, err)
		return
	}
	log.Infof("aws poll was successful, clearing old cached table")
	cachedTable = cachedTable[:0]

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			instanceObj := ec2Instance{Id: *instance.InstanceId, State: string(instance.State.Name), Type: string(instance.InstanceType)}
			for _, tag := range instance.Tags {
				if *tag.Key == environmentTagKey && *tag.Value != "" {
					instanceObj.Environment = *tag.Value
				}
				if *tag.Key == "Name" {
					instanceObj.Name = *tag.Value
				}
			}

			if instanceObj.Environment != "" {
				addInstance(&instanceObj)
			}
		}
	}
	determineEnvStates()
	log.Debugf("valid enviornments in cache: %d", len(cachedTable))
	return
}

// get instance ids for an environment
func getInstanceIds(envName string) (instanceIds []string) {
	for _, env := range cachedTable {
		if env.Name == envName {
			for _, instance := range env.Instances {
				instanceIds = append(instanceIds, instance.Id)
			}
		}
	}
	return
}

// shuts down an env
func shutdownEnv(envName string) (response []byte, err error) {
	instanceIds := getInstanceIds(envName)
	if len(instanceIds) > maxInstancesToShutdown {
		err = fmt.Errorf("SAFETY: env [%s] has too many associated instances to shutdown %d", envName, len(instanceIds))
		log.Debugf("SAFETY: instances: %v", instanceIds)
	} else if len(instanceIds) > 0 {
		input := &ec2.StopInstancesInput{
			InstanceIds: instanceIds,
			DryRun:      aws.Bool(false),
		}

		req := awsClient.StopInstancesRequest(input)
		resp, reqErr := req.Send()
		response, _ = json.MarshalIndent(resp, "", "  ")
		if reqErr != nil {
			log.Errorf("aws api error [%s]: %v", envName, reqErr)
			err = reqErr
		} else {
			log.Infof("successfully stopped env %s", envName)
		}
	} else {
		err = fmt.Errorf("env [%s] has no associated instances", envName)
	}
	return
}

// starts up an env
func startupEnv(envName string) (response []byte, err error) {
	instanceIds := getInstanceIds(envName)
	if len(instanceIds) > 0 {
		input := &ec2.StartInstancesInput{
			InstanceIds: instanceIds,
			DryRun:      aws.Bool(false),
		}

		req := awsClient.StartInstancesRequest(input)
		resp, reqErr := req.Send()
		response, _ = json.MarshalIndent(resp, "", "  ")
		if reqErr != nil {
			log.Errorf("aws api error [%s]: %v", envName, reqErr)
			err = reqErr
		} else {
			log.Infof("successfully started env %s", envName)
		}
	} else {
		err = fmt.Errorf("env [%s] has no associated instances", envName)
		log.Errorf("env [%s] has no associated instances", envName)
	}
	return
}

// returns an env Summary
func getEnvSummary() (envSummary []ec2EnvironmentSummary) {
	for _, env := range cachedTable {
		running := 0
		for _, instance := range env.Instances {
			if instance.State == "running" {
				running++
			}
		}
		envSummary = append(envSummary, ec2EnvironmentSummary{
			Name:             env.Name,
			State:            env.State,
			RunningInstances: running,
			TotalInstances:   len(env.Instances),
		})
	}
	return
}

// starts the poller
func StartPoller() {
	// build the initial cache
	refreshTable()

	pollInterval := time.Minute * time.Duration(viper.GetInt("aws.polling_interval"))
	log.Infof("start polling with interval %v", pollInterval)

	t := time.Tick(pollInterval)
	// start polling forever...
	for {
		select {
		// interval reached
		case <-t:
			refreshTable()
		}
	}
}

func main() {
	// init config and logging
	ConfigInit("")

	// start http server
	go StartServer()

	// init the aws client
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	// set aws region
	cfg.Region = awsRegion

	// pass aws client config
	awsClient = ec2.New(cfg)

	// start the poller
	StartPoller()
}
