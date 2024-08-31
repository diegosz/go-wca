package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/diegosz/go-wca/pkg/wca"
	"github.com/go-ole/go-ole"
)

var version = "latest"
var revision = "latest"

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type DeviceFlag struct {
	Value string
	IsSet bool
}

func (f *DeviceFlag) Set(value string) (err error) {
	if value == "" {
		err = fmt.Errorf("set the required device ID")
		return
	}
	if !uuidRegex.MatchString(value) {
		err = fmt.Errorf("invalid ID format")
		return
	}
	f.Value = value
	f.IsSet = true
	return
}

func (f *DeviceFlag) String() string {
	return fmt.Sprintf("%v", f.Value)
}

func main() {
	var err error
	if err = run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) (err error) {
	var deviceFlag DeviceFlag
	var versionFlag bool

	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.Var(&deviceFlag, "device", "Specify device ID (required)")
	f.Var(&deviceFlag, "d", "Alias of --device")
	f.BoolVar(&versionFlag, "version", false, "Show version")
	if err = f.Parse(args[1:]); err != nil {
		return
	}
	if versionFlag {
		fmt.Printf("%s-%s\n", version, revision)
		return
	}
	if !deviceFlag.IsSet {
		err = fmt.Errorf("a device ID is required")
		return
	}
	if err = setAudioDeviceByID(deviceFlag.Value); err != nil {
		return
	}
	fmt.Println("Successfully done")
	return
}

func setAudioDeviceByID(deviceID string) (err error) {
	GUID_IPolicyConfigVista := ole.NewGUID("{568b9108-44bf-40b4-9006-86afe5b5a620}")
	GUID_CPolicyConfigVistaClient := ole.NewGUID("{294935CE-F637-4E7C-A41B-AB255460B862}")

	if err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return
	}
	defer ole.CoUninitialize()

	var pc *wca.IPolicyConfigVista
	if err = wca.CoCreateInstance(GUID_CPolicyConfigVistaClient, 0, wca.CLSCTX_ALL, GUID_IPolicyConfigVista, &pc); err != nil {
		return
	}
	defer pc.Release()

	id := "{0.0.0.00000000}.{" + deviceID + "}"
	if err = pc.SetDefaultEndpoint(id, wca.EConsole); err != nil {
		return
	}
	return
}
