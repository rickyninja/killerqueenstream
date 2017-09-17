// camcontrol discovers the killer queen cabs on the LAN, and then sends commands to an api on the cabs to control streaming cameras.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/rickyninja/killerqueenstream"
)

var ErrAPILost = errors.New("Failed to find API")

var (
	port            int
	debug           bool
	action          string
	device          string
	apihost         string
	name            string
	videoResolution string
)

var actionDispatch map[string]func()

func init() {
	flag.IntVar(&port, "port", 14000, "Set the service port")
	flag.BoolVar(&debug, "debug", false, "print debugging info")
	flag.StringVar(&action, "action", "", "Send a cab camera command")
	flag.StringVar(&device, "device", "", "a camera device, like: /dev/video0")
	flag.StringVar(&name, "name", "", "the name of a stream")
	flag.StringVar(&videoResolution, "video-resolution", "", "1920x1080, 1280x720, etc.")
	flag.Usage = help
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	actionDispatch = map[string]func(){
		"on": func() {
			startCam(name)
		},
		"off": func() {
			stopCam(name)
		},
		"getconfig": func() {
			getStreamConfig()
		},
		"setconfig": func() {
			configureStream(device, name)
		},
		"caminfo": func() {
			getCaminfo()
		},
		"probe": func() {
			var err error
			apihost, err = probe()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("API host is: ", apihost)
		},
	}

}

func main() {
	flag.Parse()
	f, ok := actionDispatch[action]
	if !ok {
		fmt.Fprintf(os.Stderr, "%s is not a valid action!\n\n", action)
		help()
	}
	f()
}

func getStreamConfig() {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/streamconfig", apihost, port), nil)
	if err != nil {
		log.Fatal(err)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Println(buf.String())
}

func getCaminfo() {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/caminfo", apihost, port), nil)
	if err != nil {
		log.Fatal(err)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Println(buf.String())
}

func startCam(name string) {
	if name == "" {
		help()
	}
	args := map[string]string{
		"name": name,
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(args)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s:%d/startcam", apihost, port), buf)
	if err != nil {
		log.Fatal(err)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	buf = new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Println(buf.String())
}

func stopCam(name string) {
	if name == "" {
		help()
	}
	args := map[string]string{
		"name": name,
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(args)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s:%d/stopcam", apihost, port), buf)
	if err != nil {
		log.Fatal(err)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	buf = new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Println(buf.String())
}

func configureStream(device, name string) {
	if device == "" {
		help()
	}
	if name == "" {
		help()
	}
	stream := kq.NewStream()
	stream.Camera.Device = device
	stream.Name = name
	stream.VideoResolution = videoResolution
	b, err := json.Marshal(stream)
	if err != nil {
		log.Fatal(err)
	}
	buf := bytes.NewBuffer(b)
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s:%d/streamconfig", apihost, port), buf)
	if err != nil {
		log.Fatal(err)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	buf = new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Println(buf.String())
}

func probe() (string, error) {
	addrs, err := getLANAddrs()
	if err != nil {
		log.Fatal(err)
	}
	return findApiHost(addrs)
}

// find ApiHost probes a list of addresses to see if any are the api host.
func findApiHost(addrs []string) (string, error) {
	ch := make(chan string)
	for _, addr := range addrs {
		go testHost(ch, addr, port)
	}
	tic := time.NewTicker(10 * time.Second)
L:
	for {
		select {
		case <-tic.C:
			tic.Stop()
			break L
		case host := <-ch:
			return host, nil
		}
	}
	return "", ErrAPILost
}

// tesHost checks to see if a host is the api host.
func testHost(ch chan string, addr string, port int) {
	uri := fmt.Sprintf("http://%s:%d/ping", addr, port)
	resp, err := http.Get(uri)
	if err != nil {
		if debug {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body) // don't need body
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		ch <- addr
	} else {
		if debug {
			fmt.Fprintf(os.Stderr, "%s: %d\n", addr, resp.StatusCode)
		}
	}
}

/* jeremys@jeremys-desktop> cat /proc/net/arp
IP address       HW type     Flags       HW address            Mask     Device
10.0.0.147       0x1         0x2         1c:56:fe:ca:82:de     *        enp5s0
10.0.0.1         0x1         0x2         58:6d:8f:7c:c3:aa     *        enp5s0
10.0.0.133       0x1         0x2         6c:ad:f8:7f:68:ab     *        enp5s0
10.0.0.107       0x1         0x2         b8:27:eb:19:d6:c3     *        enp5s0
10.0.0.2         0x1         0x2         00:18:f8:d5:4c:eb     *        enp5s0
*/

// getLANAddrs returns a list of IP addresses observed by ARP.
func getLANAddrs() ([]string, error) {
	fd, err := os.Open("/proc/net/arp")
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	addrs := make([]string, 0)
	addrs = append(addrs, "127.0.0.1") // helps with testing
	var sawHeader bool
	for scanner.Scan() {
		if !sawHeader {
			sawHeader = true
			continue
		}
		fields := strings.Fields(scanner.Text())
		addrs = append(addrs, fields[0])
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return addrs, nil
}

func help() {
	type helpVar struct {
		Program string
		Actions string
	}
	fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
	flag.PrintDefaults()
	hv := helpVar{Program: os.Args[0]}
	for a := range actionDispatch {
		hv.Actions += fmt.Sprintf("%s, ", a)
	}
	hv.Actions = strings.TrimRight(hv.Actions, ", ")
	msg := ` 
Examples:
# get info on cameras attached to the server
{{.Program}} -action caminfo

# configure a stream
{{.Program}} -action setconfig -device /dev/video0 -name gold -video-resolution 1920x1080

# show stream configs
{{.Program}} -action getconfig

# start camera output
{{.Program}} -action on -name gold

# stop camera output
{{.Program}} -action off -name gold

# probe the LAN to find the API host
{{.Program}} -action probe

# list of possible actions
{{.Actions}}
`
	tmpl, err := template.New("help").Parse(msg)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(os.Stderr, hv)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(1)
}
