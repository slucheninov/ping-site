package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"

	"gopkg.in/yaml.v2"
)

const ComfigYaml string = "config.yaml"

// type Config struct {
// 	Site []string `yaml:"site"`
// }

func main() {
	// package yaml
	var Groups map[string][]string
	// Read config
	cfgData, err := ioutil.ReadFile(ComfigYaml)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(cfgData, &Groups); err != nil {
		log.Fatal(err)
	}
	//
	//
	for group, v := range Groups {
		fmt.Printf("===>Check groups %v: \n", group)
		for _, site := range v {
			fmt.Printf("=>Test site https://%v\n", site)
			// httptrace
			// test to http get
			// Timeout to 10 sec
			client := &http.Client{
				Timeout: time.Second * 10,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					if len(via) > 2 {
						return fmt.Errorf("max 2 hops")
					} else if len(via) == 1 {
						fmt.Println("Redirect ===> ")
					}
					return nil
				},
			}
			req, _ := http.NewRequest("GET", "https://"+site, nil)
			// added a user agent to disguise itself as a browser :)
			req.Header.Add("user-agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36`)
			// code copy https://blog.golang.org/http-tracing
			var start, connect, dns, tlsHandshake time.Time

			trace := &httptrace.ClientTrace{
				DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
				DNSDone: func(ddi httptrace.DNSDoneInfo) {
					fmt.Printf("DNS Done: %v\n", time.Since(dns))
				},

				TLSHandshakeStart: func() { tlsHandshake = time.Now() },
				TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
					fmt.Printf("TLS Handshake: %v\n", time.Since(tlsHandshake))
				},

				ConnectStart: func(network, addr string) { connect = time.Now() },
				ConnectDone: func(network, addr string, err error) {
					fmt.Printf("Connect time: %v\n", time.Since(connect))
				},

				GotFirstResponseByte: func() {
					fmt.Printf("Time from start to first byte: %v\n", time.Since(start))
				},
			}

			req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
			start = time.Now()
			if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Total time: %v\n", time.Since(start))
			// end copy
			//
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			// to clean up resources
			defer resp.Body.Close()
			fmt.Printf("Respons status code: %v\n", resp.StatusCode)
			fmt.Println("--------------------------------------")

		}
	}

}
