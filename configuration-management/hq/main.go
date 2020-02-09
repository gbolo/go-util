package main

import "github.com/spf13/viper"

func main() {

	initViper("")
	log.Infof("config file being used: %v", viper.ConfigFileUsed())

	targetHosts := getTargets()
	if len(targetHosts) < 1 {
		log.Fatalf("no target hosts have been specified")
	}
	log.Infof("found %d target host(s)", len(targetHosts))

	tasks := getTasks()
	if len(tasks) < 1 {
		log.Fatalf("no tasks have been specified")
	}
	log.Infof("found %d task(s) to submit", len(tasks))

	for _, host := range targetHosts {
		for _, task := range tasks {
			log.Infof("submitting task (%s) to target %s", task.Module, host)
			result, err := submitTask(host+taskEndpoint, task)
			if err != nil {
				log.Errorf("there was an error submitting task: %v", err)
			}
			if result.Success {
				log.Infof("task result success: %s", result.Message)
			} else {
				log.Errorf("task result error: %s", result.Message)
			}
		}
	}

}
