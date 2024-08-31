package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/diegosz/go-wca/pkg/wca"
	"github.com/go-ole/go-ole"
)

var version = "latest"
var revision = "latest"

type MuteFlag struct {
	Value bool
	IsSet bool
}

func (f *MuteFlag) Set(value string) (err error) {
	if value != "true" && value != "false" {
		err = fmt.Errorf("set 'true' or 'false'")
		return
	}
	if value == "true" {
		f.Value = true
	}
	f.IsSet = true
	return
}

func (f *MuteFlag) String() string {
	return fmt.Sprintf("%v", f.Value)
}

func main() {
	var err error
	if err = run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) (err error) {
	var muteFlag MuteFlag
	var versionFlag bool

	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.Var(&muteFlag, "mute", "Specify mute state (default is false)")
	f.Var(&muteFlag, "m", "Alias of --mute")
	f.BoolVar(&versionFlag, "version", false, "Show version")
	if err = f.Parse(args[1:]); err != nil {
		return
	}
	if versionFlag {
		fmt.Printf("%s-%s\n", version, revision)
		return
	}
	if err = muteMicrophone(muteFlag); err != nil {
		return
	}
	fmt.Println("Successfully done")
	return
}

func muteMicrophone(muteFlag MuteFlag) (err error) {
	fmt.Println(muteFlag)
	if err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return
	}
	defer ole.CoUninitialize()

	var mmde *wca.IMMDeviceEnumerator
	if err = wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return
	}
	defer mmde.Release()

	var mmd *wca.IMMDevice
	if err = mmde.GetDefaultAudioEndpoint(wca.ECapture, wca.EConsole, &mmd); err != nil {
		return
	}
	defer mmd.Release()

	var ps *wca.IPropertyStore
	if err = mmd.OpenPropertyStore(wca.STGM_READ, &ps); err != nil {
		return
	}
	defer ps.Release()

	var pv wca.PROPVARIANT
	if err = ps.GetValue(&wca.PKEY_Device_FriendlyName, &pv); err != nil {
		return
	}
	fmt.Printf("%s\n", pv.String())

	var aev *wca.IAudioEndpointVolume
	if err = mmd.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev); err != nil {
		return
	}
	defer aev.Release()

	var mute bool
	if err = aev.GetMute(&mute); err != nil {
		return
	}

	if muteFlag.IsSet {
		if err = aev.SetMute(muteFlag.Value, nil); err != nil {
			return
		}
	}

	if err = aev.GetMute(&mute); err != nil {
		return
	}

	var channels uint32
	if err = aev.GetChannelCount(&channels); err != nil {
		return
	}

	var masterVolumeLevel float32
	if err = aev.GetMasterVolumeLevel(&masterVolumeLevel); err != nil {
		return
	}

	var masterVolumeLevelScalar float32
	if err = aev.GetMasterVolumeLevelScalar(&masterVolumeLevelScalar); err != nil {
		return
	}

	fmt.Println("--------")
	fmt.Printf("Mute state: %v\n", mute)
	fmt.Printf("Channels: %d\n", channels)
	fmt.Println("Master volume level:")
	fmt.Printf("  %v [dB]\n", masterVolumeLevel)
	fmt.Printf("  %v [scalar]\n", masterVolumeLevelScalar)
	fmt.Println("--------")

	return
}
