/*
The MIT License (MIT)

Copyright (c) 2014 DutchCoders [https://github.com/dutchcoders/]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/dutchcoders/gonest"
	"log"
	"os"
	"strings"
)

var clientid string
var token string
var secret string

func init() {
	// ok this should actually be a secret, but there is no way to keep this a real secret, same as when it
	// is in a executable or apk. So just have it plain here.
	clientid = "7348d6f1-1437-4935-a6e6-2a5567c96036"
	secret = "xUVHKSKw8RTzxQrVH2UKiCnCb"
	token = os.Getenv("NEST_ACCESS_TOKEN")
}

func main() {
	var err error

	var nest *gonest.Nest
	if nest, err = gonest.Connect(clientid, token); err != nil {
		if ue, ok := err.(*gonest.UnauthorizedError); ok {
			for {
				fmt.Printf("Login to nest using url: %s\n\n%s", ue.Url, "Enter pincode: ")

				reader := bufio.NewReader(os.Stdin)
				code, _ := reader.ReadString('\n')

				code = strings.Replace(code, "\n", "", -1)

				if err = nest.Authorize(secret, code); err != nil {
					fmt.Printf("%s\n%s", err, "Enter pincode: ")
				}

				break
			}

			fmt.Printf("Successfully authorized.\nYou can now persist the accesstoken using: \nexport NEST_ACCESS_TOKEN=%s\n", nest.Token)
		} else {
			fmt.Println("Could not login to nest.")
			os.Exit(1)
		}

	}

	structureid := flag.String("structure", "", "operate on structure")
	thermostatid := flag.String("thermostat", "", "operate on thermostat")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Printf("Usage:\n\n")
		fmt.Printf("nest --structure {structureid} [home|away]\n")
		fmt.Printf("nest structures\n")
		fmt.Printf("nest thermostats\n")
		os.Exit(0)
	}

	switch flag.Args()[0] {
	case "away", "home":
		if *structureid == "" {
			fmt.Println("Structure not set\n\nUsage: nest --structure {structureid} [home|away]")
			os.Exit(1)
		}

		away := flag.Args()[0]

		if err = nest.Set(fmt.Sprintf("structures/%s", *structureid), map[string]interface{}{"away": away}); err != nil {
			log.Panic(err)
		}

		fmt.Printf("Set to: %s\n", away)
	case "structures":
		var structures map[string]gonest.Structure
		if err = nest.Structures(&structures); err != nil {
			log.Panic(err)
		}

		fmt.Printf("Structures:\n")
		for _, structure := range structures {
			fmt.Printf("Id:             %s\n", structure.StructureId)
			fmt.Printf("Name:           %s\n", structure.Name)
			fmt.Printf("Away:           %s\n", structure.Away)
		}
	case "thermostats":
		var devices gonest.Devices
		if err = nest.Devices(&devices); err != nil {
			log.Panic(err)
		}

		fmt.Printf("Thermostats:\n")
		for _, device := range devices.Thermostats {
			fmt.Printf("Id:             %s\n", device.DeviceId)
			fmt.Printf("Name:           %s\n", device.Name)
			fmt.Printf("Name(long):     %s\n", device.NameLong)
			fmt.Printf("Temperature:    %f C\n", device.AmbientTemperatureC)
			fmt.Printf("Humidity:       %f\n", device.Humidity)
			if device.IsOnline {
				fmt.Printf("Status:         online\n")
			} else {
				fmt.Printf("Status:         offline\n")
			}
		}

	default:
		if *thermostatid != "" {
			var devices gonest.Devices
			if err = nest.Devices(&devices); err != nil {
				log.Panic(err)
			}

			var thermostat gonest.Thermostat
			var match bool
			if thermostat, match = devices.Thermostats[*thermostatid]; !match {
				fmt.Printf("Thermostat %s not found.", thermostatid)
				os.Exit(1)
			}

			args := map[string]interface{}{
				"deviceid":              thermostat.DeviceId,
				"version":               thermostat.SoftwareVersion,
				"name":                  thermostat.Name,
				"online":                thermostat.IsOnline,
				"away-temperature-high": thermostat.AwayTemperatureHighC,
				"away-temperature-low":  thermostat.AwayTemperatureLowC,
				"ambient-temperature":   thermostat.AmbientTemperatureC,
				"target-temperature":    thermostat.TargetTemperatureC,
			}

			if val, ok := args[flag.Args()[0]]; ok {
				fmt.Println(val)
			} else {
				fmt.Printf("Unsupported argument %s\n", flag.Args()[0])
			}

		} else if *structureid != "" {
			fmt.Printf("Not implemented yet\n")
		} else {
			fmt.Printf("Unknown command: %s\n", flag.Args()[0])
		}
	}
}
