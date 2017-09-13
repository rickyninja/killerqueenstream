package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/rickyninja/killerqueenstream"
)

func main() {
	router := httprouter.New()
	router.GET("/", Help)
	router.GET("/caminfo", CamInfo)
	router.GET("/ping", Ping)
	router.PUT("/streamconfig", SetStreamConfig)
	router.GET("/streamconfig", ViewStreamConfig)
	router.PUT("/startcam", StartCam)
	router.PUT("/stopcam", StopCam)
	log.Fatal(http.ListenAndServe(":14000", router))
}

func Ping(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	pong := map[string]bool{
		"pong": true,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	err := enc.Encode(pong)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Help(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	type end struct {
		Description string
		Route       string
		Method      string
	}
	ends := []end{
		end{
			Description: "The help page",
			Route:       "/",
			Method:      "GET",
		},
		end{
			Description: "Get info about camera devices.",
			Route:       "/caminfo",
			Method:      "GET",
		},
		end{
			Description: "Responds to api pings to aid LAN discovery of host api runs on.",
			Route:       "/ping",
			Method:      "GET",
		},
		end{
			Description: "Set values that will be used to stream camera output to an rtmp service.",
			Route:       "/streamconfig",
			Method:      "PUT",
		},
		end{
			Description: "Get stream configuration values that were previously configured.",
			Route:       "/streamconfig",
			Method:      "GET",
		},
		end{
			Description: "Start a camera stream.",
			Route:       "/startcam",
			Method:      "PUT",
		},
		end{
			Description: "Stop a camera stream.",
			Route:       "/stopcam",
			Method:      "PUT",
		},
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	err := enc.Encode(ends)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ViewStreamConfig(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	err := enc.Encode(running)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SetStreamConfig(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	dec := json.NewDecoder(r.Body)
	stream := kq.NewStream()
	err := dec.Decode(stream)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	running.Lock()
	running.Stream[stream.Name] = stream
	running.Unlock()

	cam, err := getCamDetail(stream.Camera.Device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stream.Camera = cam

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	err = enc.Encode(stream)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// track which cameras are running by client provided name (gold, blue, etc.).
type streamMap map[string]*kq.Stream

type streamInfo struct {
	sync.Mutex
	Stream streamMap `json:"stream"`
}

var running streamInfo = streamInfo{
	Stream: make(streamMap),
}

func StartCam(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	name, err := getname(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stream, ok := running.Stream[name]
	if !ok {
		http.Error(w, fmt.Sprintf("Failed to find stream by %s.  Did you configure the stream?", name), http.StatusInternalServerError)
		return
	}
	uriStr := stream.Start()
	uri, err := url.Parse(uriStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uri.Host = r.Host

	resp := map[string]string{
		"status": "on",
		"url":    uri.String(),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	err = enc.Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func StopCam(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	name, err := getname(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stream, ok := running.Stream[name]
	if !ok {
		http.Error(w, fmt.Sprintf("Failed to find stream by %s.  Did you configure the stream?", name), http.StatusInternalServerError)
		return
	}
	err = stream.Stop()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"status": "off",
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	err = enc.Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getname(r *http.Request) (string, error) {
	dec := json.NewDecoder(r.Body)
	p := map[string]string{}
	err := dec.Decode(&p)
	if err != nil {
		return "", err
	}
	name, ok := p["name"]
	if !ok {
		return "", errors.New("name missing from request")
	}
	return name, nil
}

func CamInfo(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	cams, err := probeCameras()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	err = enc.Encode(cams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/* jeremys@jeremys-desktop> udevadm info --query=all --name=/dev/video0
P: /devices/pci0000:00/0000:00:1a.7/usb1/1-6/1-6.1/1-6.1.2/1-6.1.2:1.2/video4linux/video0
N: video0
S: v4l/by-id/usb-046d_HD_Webcam_C615_1563CF40-video-index0
S: v4l/by-path/pci-0000:00:1a.7-usb-0:6.1.2:1.2-video-index0
E: COLORD_DEVICE=1
E: COLORD_KIND=camera
E: DEVLINKS=/dev/v4l/by-path/pci-0000:00:1a.7-usb-0:6.1.2:1.2-video-index0 /dev/v4l/by-id/usb-046d_HD_Webcam_C615_1563CF40-video-index0
E: DEVNAME=/dev/video0
E: DEVPATH=/devices/pci0000:00/0000:00:1a.7/usb1/1-6/1-6.1/1-6.1.2/1-6.1.2:1.2/video4linux/video0
E: ID_BUS=usb
E: ID_FOR_SEAT=video4linux-pci-0000_00_1a_7-usb-0_6_1_2_1_2
E: ID_MODEL=HD_Webcam_C615
E: ID_MODEL_ENC=HD\x20Webcam\x20C615
E: ID_MODEL_ID=082c
E: ID_PATH=pci-0000:00:1a.7-usb-0:6.1.2:1.2
E: ID_PATH_TAG=pci-0000_00_1a_7-usb-0_6_1_2_1_2
E: ID_REVISION=0011
E: ID_SERIAL=046d_HD_Webcam_C615_1563CF40
E: ID_SERIAL_SHORT=1563CF40
E: ID_TYPE=video
E: ID_USB_DRIVER=uvcvideo
E: ID_USB_INTERFACES=:010100:010200:0e0100:0e0200:
E: ID_USB_INTERFACE_NUM=02
E: ID_V4L_CAPABILITIES=:capture:
E: ID_V4L_PRODUCT=HD Webcam C615
E: ID_V4L_VERSION=2
E: ID_VENDOR=046d
E: ID_VENDOR_ENC=046d
E: ID_VENDOR_ID=046d
E: MAJOR=81
E: MINOR=0
E: SUBSYSTEM=video4linux
E: TAGS=:seat:uaccess:
E: USEC_INITIALIZED=400948839519
*/

type caminfo map[string]map[string]string

func probeCameras() ([]kq.Camera, error) {
	vdevs, err := getVideoDevices()
	if err != nil {
		return nil, err
	}
	cameras := make([]kq.Camera, 0)
	for _, dev := range vdevs {
		cam, err := getCamDetail(dev)
		if err != nil {
			return nil, err
		}
		cameras = append(cameras, cam)
	}
	return cameras, nil
}

func getCamDetail(dev string) (kq.Camera, error) {
	cam := kq.Camera{}
	attr := make(map[string]string)
	cmd := exec.Command("udevadm", "info", "--query=all", "--name="+dev)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return cam, err
	}
	err = cmd.Start()
	if err != nil {
		return cam, err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || !strings.Contains(line, "=") {
			continue
		}
		fields := strings.Fields(line)
		keyvals := strings.Split(fields[1], "=")
		attr[keyvals[0]] = keyvals[1]
	}
	if err := scanner.Err(); err != nil {
		return cam, err
	}
	err = cmd.Wait()
	if err != nil {
		return cam, err
	}
	cam = kq.Camera{
		Serial: attr["ID_SERIAL_SHORT"],
		Model:  attr["ID_MODEL"],
		Vendor: attr["ID_VENDOR"],
		IdPath: attr["ID_PATH"],
		Device: dev,
	}
	return cam, nil
}

func getVideoDevices() ([]string, error) {
	vdevs := make([]string, 0)
	fd, err := os.Open("/dev")
	if err != nil {
		return vdevs, err
	}
	files, err := fd.Readdir(0)
	if err != nil {
		return vdevs, err
	}
	for _, f := range files {
		if strings.Contains(f.Name(), "video") {
			vdevs = append(vdevs, "/dev/"+f.Name())
		}
	}
	return vdevs, nil
}
