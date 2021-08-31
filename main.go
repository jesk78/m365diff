package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Services []struct {
	ID                     int      `json:"id"`
	ServiceArea            string   `json:"serviceArea"`
	ServiceAreaDisplayName string   `json:"serviceAreaDisplayName"`
	Urls                   []string `json:"urls,omitempty"`
	Ips                    []string `json:"ips,omitempty"`
	TCPPorts               string   `json:"tcpPorts,omitempty"`
	ExpressRoute           bool     `json:"expressRoute"`
	Category               string   `json:"category"`
	Required               bool     `json:"required"`
	Notes                  string   `json:"notes,omitempty"`
	UDPPorts               string   `json:"udpPorts,omitempty"`
}

func main() {
	resp, err := http.Get("https://endpoints.office.com/endpoints/Worldwide?ClientRequestId=5336be43-f4fc-48ac-b4df-789e7466f6b8")
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var result Services
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte into slice of structs (type Services)
		fmt.Println("Can not unmarshal JSON")
	}
	// map of services (eg. tcp/123) to slice of CIDR networks
	smap := make(map[string][]string)
	for _, v := range result {
		// TCP keys(services)
		if len(v.Ips) > 0 && v.TCPPorts != "" {
			tcpPorts := strings.Split(v.TCPPorts, ",")
			for _, p := range tcpPorts {
				service := "tcp/" + p
				for _, ip := range v.Ips {
					smap[service] = append(smap[service], ip)
				}
			}
			// UDP keys(services)
		} else if len(v.Ips) > 0 && v.UDPPorts != "" {
			udpPorts := strings.Split(v.UDPPorts, ",")
			for _, p := range udpPorts {
				service := "udp/" + p
				for _, ip := range v.Ips {
					smap[service] = append(smap[service], ip)
				}
			}
		}
	}
	for k, v := range smap {
		fmt.Print(k, "->", v, "\n\n\n")
	}
}
