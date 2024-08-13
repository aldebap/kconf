////////////////////////////////////////////////////////////////////////////////
//	main.go  -  Jul-5-2024  -  aldebap
//
//	Entry point for Kong Configuration tool
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	versionInfo string = "kconf 0.1"
)

// execution options
type Options struct {
	jsonOutput bool
	verbose    bool
}

// main entry point for kconf
func main() {
	var (
		version bool

		//	Kong server configuration
		kongAddress string
		kongPort    int

		options Options
	)

	//	CLI arguments
	flag.BoolVar(&version, "version", false, "show kconf version")

	flag.StringVar(&kongAddress, "kong-address", "localhost", "Kong configuration address")
	flag.IntVar(&kongPort, "port", 8001, "Kong configuration port")
	flag.BoolVar(&options.jsonOutput, "json-output", false, "use json output for every command")
	flag.BoolVar(&options.verbose, "verbose", false, "run in verbose mode")

	flag.Parse()

	//	version option
	if version {
		fmt.Printf("%s\n", versionInfo)
		return
	}

	//	connect and send command
	kongServer := NewKongServer(kongAddress, kongPort)
	if kongServer == nil {
		fmt.Fprintf(os.Stderr, "[error] fail attempting to alocate Kong server\n")
		os.Exit(-1)
	}

	err := kconf(kongServer, flag.Args(), options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[error] %s\n", err.Error())
		os.Exit(-1)
	}
}

// kconf utility
func kconf(myKongServer *KongServer, command []string, options Options) error {

	if len(command) == 0 {
		return errors.New("missing command: available commands: status, add, query, list")
	}

	//	command to get Kong status
	if command[0] == "status" {
		return myKongServer.CheckStatus(options)
	}

	//	command add
	if command[0] == "add" {
		return commandAdd(myKongServer, command[1:], options)
	}

	//	command query
	if command[0] == "query" {
		return commandQuery(myKongServer, command[1:], options)
	}

	//	command list
	if command[0] == "list" {
		return commandList(myKongServer, command[1:], options)
	}

	//	command delete
	if command[0] == "delete" {
		return commandDelete(myKongServer, command[1:], options)
	}

	return errors.New("invalid command: " + command[0])
}

// command add
func commandAdd(myKongServer *KongServer, command []string, options Options) error {

	if len(command) == 0 {
		return errors.New("missing entity for command add: available entities: service, route")
	}

	//	compile all regex required to extract parameters for command add
	nameRegEx, err := regexp.Compile(`^--name\s*=\s*(\S+)\s*$`)
	if err != nil {
		return err
	}

	urlRegEx, err := regexp.Compile(`^--url\s*=\s*(\S.*)\s*$`)
	if err != nil {
		return err
	}

	protocolsRegEx, err := regexp.Compile(`^--protocols\s*=\s*(\S.*)\s*$`)
	if err != nil {
		return err
	}

	methodsRegEx, err := regexp.Compile(`^--methods\s*=\s*(\S.*)\s*$`)
	if err != nil {
		return err
	}

	pathsRegEx, err := regexp.Compile(`^--paths\s*=\s*(\S.*)\s*$`)
	if err != nil {
		return err
	}

	serviceIdRegEx, err := regexp.Compile(`^--service-id\s*=\s*(\S.*)\s*$`)
	if err != nil {
		return err
	}

	if command[0] == "service" {
		var name string
		var url string
		var enabled bool = true

		for i := 1; i < len(command); i++ {
			match := nameRegEx.FindAllStringSubmatch(command[i], -1)
			if len(match) == 1 {
				name = match[0][1]
			}

			match = urlRegEx.FindAllStringSubmatch(command[i], -1)
			if len(match) == 1 {
				url = match[0][1]
			}
		}
		newService := NewKongService(name, url, enabled)

		return myKongServer.AddService(newService, options)
	}

	if command[0] == "route" {
		const valuesDelim = ","
		var name string
		var protocols []string
		var methods []string
		var paths []string
		var serviceId string

		for i := 1; i < len(command); i++ {
			match := nameRegEx.FindAllStringSubmatch(command[i], -1)
			if len(match) == 1 {
				name = match[0][1]
			}

			match = protocolsRegEx.FindAllStringSubmatch(command[i], -1)
			if len(match) == 1 {
				protocols = strings.Split(match[0][1], valuesDelim)
			}

			match = methodsRegEx.FindAllStringSubmatch(command[i], -1)
			if len(match) == 1 {
				methods = strings.Split(match[0][1], valuesDelim)
			}

			match = pathsRegEx.FindAllStringSubmatch(command[i], -1)
			if len(match) == 1 {
				paths = strings.Split(match[0][1], valuesDelim)
			}

			match = serviceIdRegEx.FindAllStringSubmatch(command[i], -1)
			if len(match) == 1 {
				serviceId = match[0][1]
			}
		}
		newKongRoute := NewKongRoute(name, protocols, methods, paths, serviceId)

		return myKongServer.AddRoute(newKongRoute, options)
	}

	return errors.New("invalid entity for command add: " + command[0])
}

// command query
func commandQuery(myKongServer *KongServer, command []string, options Options) error {

	if len(command) == 0 {
		return errors.New("missing entity for command query: available entities: service, route")
	}

	//	compile all regex required to extract parameters for command query
	idRegEx, err := regexp.Compile(`^--id\s*=\s*(\S.*)\s*$`)
	if err != nil {
		return err
	}

	if command[0] == "service" {
		var id string

		for i := 1; i < len(command); i++ {
			match := idRegEx.FindAllStringSubmatch(command[i], -1)
			if len(match) == 1 {
				id = match[0][1]
			}
		}

		if len(id) == 0 {
			return errors.New("missing service id: option --id={id} required for this command")
		}

		return myKongServer.QueryService(id, options)
	}

	if command[0] == "route" {
		var id string

		for i := 1; i < len(command); i++ {
			match := idRegEx.FindAllStringSubmatch(command[i], -1)
			if len(match) == 1 {
				id = match[0][1]
			}
		}

		if len(id) == 0 {
			return errors.New("missing route id: option --id={id} required for this command")
		}

		return myKongServer.QueryRoute(id, options)
	}

	return errors.New("invalid entity for command query: " + command[0])
}

// command list
func commandList(myKongServer *KongServer, command []string, options Options) error {

	if len(command) == 0 {
		return errors.New("missing entity for command list: available entities: service, route")
	}

	if command[0] == "service" {
		return myKongServer.ListServices(options)
	}

	if command[0] == "route" {
		return myKongServer.ListRoutes(options)
	}

	return errors.New("invalid entity for command list: " + command[0])
}

// command delete
func commandDelete(myKongServer *KongServer, command []string, options Options) error {

	if len(command) == 0 {
		return errors.New("missing entity for command delete: available entities: service, route")
	}

	//	compile all regex required to extract parameters for command query
	idRegEx, err := regexp.Compile(`^--id\s*=\s*(\S.*)\s*$`)
	if err != nil {
		return err
	}

	if command[0] == "service" {
		var id string

		for i := 1; i < len(command); i++ {
			match := idRegEx.FindAllStringSubmatch(command[i], -1)
			if len(match) == 1 {
				id = match[0][1]
			}
		}

		if len(id) == 0 {
			return errors.New("missing service id: option --id={id} required for this command")
		}

		return myKongServer.DeleteService(id, options)
	}

	/*
		if command[0] == "route" {
			return myKongServer.ListRoutes(options)
		}
	*/

	return errors.New("invalid entity for command list: " + command[0])
}
