package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	dac "github.com/xinsnake/go-http-digest-auth-client"
)

const (
	method    = "GET"
	checkName = "CheckVOX"
	host      = "%s"
	schema    = "%s://"
	uri       = "/%s/service?action=get_gsminfo"
)

func main() {
	username := flag.String("user", "admin", "username for digest auth")
	password := flag.String("pass", "admin", "password for digest auth")
	hostPtr := flag.String("host", "", "192.168.1.254 (required)")
	schemaPtr := flag.String("schema", "http", "http or https")
	slotPtr := flag.String("slot", "", "like: /2/service?action=get_gsminfo (2 - is a slot)")
	modemPtr := flag.String("modem", "1", "min: 1, max: 4")
	critSigLevel := flag.Int("crit", 7, "critical signal")
	help := flag.Bool("h", false, "-h for help")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(2)
	}

	u := fmt.Sprintf(schema+host+uri, *schemaPtr, *hostPtr, *slotPtr)

	// get response
	service := make(map[string]interface{})
	if err := json.Unmarshal(readWithDigest(u, *username, *password), &service); err != nil {
		fmt.Printf("%s: prepare JOSN error | %s\n", checkName, err)
		os.Exit(2)
	}

	// modem number validation
	modemNum, err := strconv.Atoi(*modemPtr)
	if err != nil {
		fmt.Printf("%s: modem number error | %s\n", checkName, err)
		os.Exit(2)
	}
	if modemNum > 4 || modemNum < 1 {
		fmt.Printf("%s: modem number error | %s\n", checkName, "must be between 1 - 4")
		os.Exit(2)
	}

	// get modem state from server
	getStateByModem(*modemPtr, service, *critSigLevel)
}

func readWithDigest(url, username, password string) []byte {
	t := dac.NewTransport(username, password)
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Printf("%s: request error | %s\n", checkName, err)
		os.Exit(2)
	}

	resp, err := t.RoundTrip(req)
	if err != nil {
		fmt.Printf("%s: trip error | %s\n", checkName, err)
		os.Exit(2)
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s: read body error | %s\n", checkName, err)
		os.Exit(2)
	}

	return r
}

func getStateByModem(modemID string, service map[string]interface{}, critLevel int) {

	modemContent := service[modemID].([]interface{})[0].(map[string]interface{})

	sig, err := strconv.Atoi(modemContent["signal"].(string))
	if err != nil {
		fmt.Printf("%s: %s; value: %v\n", checkName, err, modemContent["signal"])
		os.Exit(2)
	}

	// unregistered device
	if !strings.Contains(modemContent["register"].(string), "Registered") {
		fmt.Printf("%s: %s - %s | signal: %s\n", checkName, modemContent["operator"], modemContent["register"], modemContent["signal"])
		os.Exit(2)
	}

	// signal level
	if sig >= critLevel {
		fmt.Printf("%s: %s - %s | signal: %s\n", checkName, modemContent["operator"], modemContent["register"], modemContent["signal"])
		os.Exit(0)
	} else {
		fmt.Printf("%s: %s - %s | low signal: %s(%d)\n", checkName, modemContent["operator"], modemContent["register"], modemContent["signal"], critLevel)
		os.Exit(2)
	}
}
