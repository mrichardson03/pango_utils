package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/commit"
)

type Credentials struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	Password string `json:"password"`
	ApiKey   string `json:"api_key"`
	Protocol string `json:"protocol"`
	Port     uint   `json:"port"`
	Timeout  int    `json:"timeout"`
}

func getCredentials(configFile, hostname, username, password, apiKey string) Credentials {
	var (
		config Credentials
		val    string
		ok     bool
	)

	// Auth from the config file.
	if configFile != "" {
		fd, err := os.Open(configFile)
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
		defer fd.Close()

		dec := json.NewDecoder(fd)
		err = dec.Decode(&config)
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
	}

	// Auth from env variables.
	if val, ok = os.LookupEnv("PANOS_HOSTNAME"); ok {
		config.Hostname = val
	}
	if val, ok = os.LookupEnv("PANOS_USERNAME"); ok {
		config.Username = val
	}
	if val, ok = os.LookupEnv("PANOS_PASSWORD"); ok {
		config.Password = val
	}
	if val, ok = os.LookupEnv("PANOS_API_KEY"); ok {
		config.ApiKey = val
	}

	// Auth from CLI args.
	if hostname != "" {
		config.Hostname = hostname
	}
	if username != "" {
		config.Username = username
	}
	if password != "" {
		config.Password = password
	}
	if apiKey != "" {
		config.ApiKey = apiKey
	}

	if config.Hostname == "" {
		log.Fatalf("ERROR: No hostname specified")
	} else if config.Username == "" && config.ApiKey == "" {
		log.Fatalf("ERROR: No username specified")
	} else if config.Password == "" && config.ApiKey == "" {
		log.Fatalf("ERROR: No password specified")
	}

	return config
}

func main() {
	var (
		err                                              error
		configFile, hostname, username, password, apiKey string
		job                                              uint
	)

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	flag.StringVar(&configFile, "config", "", "JSON config file with panos connection info")
	flag.StringVar(&hostname, "host", "", "PAN-OS hostname")
	flag.StringVar(&username, "user", "", "PAN-OS username")
	flag.StringVar(&password, "pass", "", "PAN-OS password")
	flag.StringVar(&apiKey, "key", "", "PAN-OS API key")
	flag.Parse()

	config := getCredentials(configFile, hostname, username, password, apiKey)

	fw := &pango.Firewall{Client: pango.Client{
		Hostname: config.Hostname,
		Username: config.Username,
		Password: config.Password,
		ApiKey:   config.ApiKey,
		Protocol: config.Protocol,
		Port:     config.Port,
		Timeout:  config.Timeout,
		Logging:  pango.LogOp | pango.LogAction,
	}}
	if err = fw.Initialize(); err != nil {
		log.Fatalf("Failed: %s", err)
	}

	// Build commit command
	cmd := commit.FirewallCommit{
		Description:             flag.Arg(0),
		ExcludeDeviceAndNetwork: false,
		ExcludeSharedObjects:    false,
		ExcludePolicyAndObjects: false,
		Force:                   false,
	}

	admins := strings.TrimSpace(config.Username)
	if admins != "" {
		cmd.Admins = strings.Split(admins, ",")
	}

	job, _, err = fw.Commit(cmd, "", nil)
	if err != nil {
		log.Fatalf("Error in commit: %s", err)
	} else if job == 0 {
		log.Printf("No commit needed")
	} else {
		log.Printf("Committed config successfully")
	}
}
