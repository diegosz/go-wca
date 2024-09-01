package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/diegosz/go-wca/pkg/wca"
	"github.com/go-ole/go-ole"
)

const (
	wcaRetryInterval  = 100 * time.Millisecond
	wcaRetryCount     = 5
	defaultDeviceName = "Default"
)

var version = "latest"
var revision = "latest"

type ListenFlag struct {
	Value string
	IsSet bool
}

func (f *ListenFlag) Set(value string) (err error) {
	if value == "" {
		value = "true"
	}
	f.Value = value
	f.IsSet = true
	return
}

func (f *ListenFlag) String() string {
	return f.Value
}

func (f *ListenFlag) Bool() bool {
	return f.Value == "true"
}

type InFlag struct {
	Value string
	IsSet bool
}

func (f *InFlag) Set(value string) (err error) {
	f.Value = value
	f.IsSet = true
	return
}

func (f *InFlag) String() string {
	return f.Value
}

type OutFlag struct {
	Value string
	IsSet bool
}

func (f *OutFlag) Set(value string) (err error) {
	f.Value = value
	f.IsSet = true
	return
}

func (f *OutFlag) String() string {
	return f.Value
}

func main() {
	var err error
	if err = run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) (err error) {
	var listenFlag ListenFlag
	var inFlag InFlag
	var outFlag OutFlag
	var versionFlag bool

	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.Var(&listenFlag, "listen", "Enable or Disable listening (default: true)")
	f.Var(&inFlag, "in", "Specify the input device display name (required)")
	f.Var(&outFlag, "out", "Specify the output device display name")
	f.BoolVar(&versionFlag, "version", false, "Show version")
	if err = f.Parse(args[1:]); err != nil {
		return
	}
	if versionFlag {
		fmt.Printf("%s-%s\n", version, revision)
		return
	}
	if !listenFlag.IsSet {
		_ = listenFlag.Set("true")
	}
	if !inFlag.IsSet {
		err = fmt.Errorf("an input device name is required")
		return
	}
	if listenFlag.Value == "true" && outFlag.Value == "" {
		_ = outFlag.Set(defaultDeviceName)
	}
	return listenMicrophone(listenFlag.Bool(), inFlag.Value, outFlag.Value)
}

func listenMicrophone(enable bool, in, out string) error {
	// Initialize COM library.
	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return fmt.Errorf("failed to initialize COM library: %v", err)
	}
	defer ole.CoUninitialize()

	// Initialize Windows Core Audio.
	var enumerator *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &enumerator); err != nil {
		return fmt.Errorf("failed to create device enumerator: %v", err)
	}
	defer enumerator.Release()

	// Enumerate audio capturing endpoints.
	var cdc *wca.IMMDeviceCollection
	if err := enumerator.EnumAudioEndpoints(wca.ECapture, wca.DEVICE_STATE_ACTIVE, &cdc); err != nil {
		return fmt.Errorf("failed to enumerate audio capturing endpoints: %v", err)
	}
	defer cdc.Release()
	inDevice, inInfo, err := findDevice(cdc, DeviceMatcherByName(in))
	if err != nil {
		return fmt.Errorf("failed to find input device: %v", err)
	}
	defer func() {
		if inDevice != nil {
			_ = inDevice.Release()
		}
	}()
	// Enumerate audio rendering endpoints.
	var rdc *wca.IMMDeviceCollection
	if err := enumerator.EnumAudioEndpoints(wca.ERender, wca.DEVICE_STATE_ACTIVE, &rdc); err != nil {
		return fmt.Errorf("failed to enumerate audio rendering endpoints: %v", err)
	}

	var ps *wca.IPropertyStore
	if err = inDevice.OpenPropertyStore(wca.STGM_READ_WRITE, &ps); err != nil {
		return fmt.Errorf("failed to open input device property store: %v", err)
	}
	defer ps.Release()

	ls, err := getListenSetting(ps)
	if err != nil {
		return fmt.Errorf("failed to get listen setting: %v", err)
	}
	fmt.Printf("%s\n", ls.String())

	switch enable {
	case true:
		var outID, outLog string
		if out != defaultDeviceName {
			outDevice, outInfo, err := findDevice(rdc, DeviceMatcherByName(out))
			if err != nil {
				return fmt.Errorf("failed to find output device: %v", err)
			}
			defer func() {
				if outDevice != nil {
					_ = outDevice.Release()
				}
			}()
			outID = outInfo.ID
			outLog = outInfo.String()
		} else {
			outID = ""
			outLog = defaultDeviceName
		}
		if err = enableListening(ps, outID); err != nil {
			return fmt.Errorf("failed to enable listening: %v", err)
		}
		fmt.Printf("Listening to %s with %s\n", in, out)
		fmt.Printf("Input device: %s\n", inInfo)
		fmt.Printf("Output device: %s\n", outLog)
	case false:
		if err = disableListening(ps); err != nil {
			return fmt.Errorf("failed to disable listening: %v", err)
		}
		fmt.Printf("Listening %s disabled\n", in)
		fmt.Printf("Input device: %s\n", inInfo)
	}

	return nil
}

func enableListening(ps *wca.IPropertyStore, outputID string) (err error) {
	if ps == nil {
		return fmt.Errorf("property store is null")
	}
	// Enable the listen setting.
	enabled := wca.NewBoolPropVariant(true)
	if err = ps.SetValue(&wca.PKEY_Listen_Setting_Enable, &enabled); err != nil {
		return fmt.Errorf("failed to enable listen setting: %v", err)
	}
	// Disable the save power setting.
	savePower := wca.NewBoolPropVariant(false)
	if err = ps.SetValue(&wca.PKEY_Listen_Setting_SavePower, &savePower); err != nil {
		return fmt.Errorf("failed to disable listen save power: %v", err)
	}
	if outputID != "" {
		// Set the listen device.
		var npv wca.PROPVARIANT
		if npv, err = wca.NewStringPropVariant(outputID); err != nil {
			return fmt.Errorf("failed to create listen device setting: %v", err)
		}
		if err = ps.SetValue(&wca.PKEY_Listen_Setting_Device, &npv); err != nil {
			return fmt.Errorf("failed to set listen device: %v", err)
		}
	} else {
		// Set the listen device to default.
		empty := wca.NewEmptyPropVariant()
		if err = ps.SetValue(&wca.PKEY_Listen_Setting_Device, &empty); err != nil {
			return fmt.Errorf("failed to set listen device: %v", err)
		}
	}
	return nil
}

func disableListening(ps *wca.IPropertyStore) (err error) {
	if ps == nil {
		return fmt.Errorf("property store is null")
	}
	// Disable the listen setting.
	disabled := wca.NewBoolPropVariant(false)
	if err = ps.SetValue(&wca.PKEY_Listen_Setting_Enable, &disabled); err != nil {
		return fmt.Errorf("failed to disable listen setting: %v", err)
	}
	// Disable the save power setting.
	savePower := wca.NewBoolPropVariant(false)
	if err = ps.SetValue(&wca.PKEY_Listen_Setting_SavePower, &savePower); err != nil {
		return fmt.Errorf("failed to disable listen save power: %v", err)
	}
	// Set the listen device to default.
	empty := wca.NewEmptyPropVariant()
	if err = ps.SetValue(&wca.PKEY_Listen_Setting_Device, &empty); err != nil {
		return fmt.Errorf("failed to set listen device: %v", err)
	}
	return nil
}

func getListenSetting(ps *wca.IPropertyStore) (state *ListenState, err error) {
	if ps == nil {
		return nil, fmt.Errorf("property store is null")
	}
	state = &ListenState{}
	// Get the liste setting.
	var found bool
	for j := 0; j < wcaRetryCount; j++ {
		var lpv wca.PROPVARIANT
		if err := ps.GetValue(&wca.PKEY_Listen_Setting_Enable, &lpv); err != nil {
			time.Sleep(wcaRetryInterval)
			continue
		}
		v, err := lpv.Bool()
		if err != nil {
			return nil, fmt.Errorf("failed to convert listen setting: %v", err)
		}
		state.Enabled = v
		found = true
		break
	}
	if !found {
		return nil, fmt.Errorf("failed to retrieve listen setting")
	}
	// Get the save power setting.
	found = false
	for j := 0; j < wcaRetryCount; j++ {
		var lpv wca.PROPVARIANT
		if err := ps.GetValue(&wca.PKEY_Listen_Setting_SavePower, &lpv); err != nil {
			time.Sleep(wcaRetryInterval)
			continue
		}
		v, err := lpv.Bool()
		if err != nil {
			return nil, fmt.Errorf("failed to convert listen save power: %v", err)
		}
		state.SavePower = v
		found = true
		break
	}
	if !found {
		return nil, fmt.Errorf("failed to retrieve listen save power")
	}
	// Get the listen device setting.
	found = false
	for j := 0; j < wcaRetryCount; j++ {
		var lpv wca.PROPVARIANT
		if err := ps.GetValue(&wca.PKEY_Listen_Setting_Device, &lpv); err != nil {
			time.Sleep(wcaRetryInterval)
			continue
		}
		state.DeviceID = lpv.String()
		found = true
		break
	}
	if !found {
		return nil, fmt.Errorf("failed to retrieve listen save power")
	}
	return state, nil
}

// findDevice finds a device in a collection that matches the given matcher. The
// matcher is a function that returns true if the device is the one we're
// looking for.
//
// IMPORTANT: The caller is responsible for releasing the device object.
func findDevice(collection *wca.IMMDeviceCollection, matcher deviceMatcher) (device *wca.IMMDevice, info *DeviceInfo, err error) {
	info = &DeviceInfo{}
	if collection == nil {
		return device, info, fmt.Errorf("collection is nil")
	}
	// Get device count.
	var count uint32
	if err := collection.GetCount(&count); err != nil {
		return device, info, fmt.Errorf("failed to get device count: %v", err)
	}
	// Iterate over devices.
	for i := uint32(0); i < count; i++ {
		if err := collection.Item(i, &device); err != nil {
			return device, info, fmt.Errorf("failed to get device device: %v", err)
		}
		// Get the device ID.
		var deviceId string
		if err := device.GetId(&deviceId); err != nil {
			return device, info, fmt.Errorf("failed to get device ID: %v", err)
		}
		// Open the property store.
		var ps *wca.IPropertyStore
		if err := device.OpenPropertyStore(wca.STGM_READ, &ps); err != nil {
			return device, info, fmt.Errorf("failed to open property store: %v", err)
		}
		defer ps.Release()
		// Get the friendly name of the device.
		var propVariant wca.PROPVARIANT
		if err := ps.GetValue(&wca.PKEY_Device_FriendlyName, &propVariant); err != nil {
			return device, info, fmt.Errorf("failed to get friendly name: %v", err)
		}
		deviceName := propVariant.String()
		defer func() {
			_ = propVariant.Clear()
		}()
		// Check if the device is the one we're looking for.
		di := &DeviceInfo{ID: deviceId, Name: deviceName}
		if matcher(di) {
			return device, di, nil
		}
	}
	return device, info, fmt.Errorf("device not found")
}

type deviceMatcher func(device *DeviceInfo) bool

// DeviceMatcherByName uses the device's display name as shown in the Windows UI
// to find a device.
func DeviceMatcherByName(name string) deviceMatcher {
	return func(device *DeviceInfo) bool {
		if device == nil {
			return false
		}
		return device.Name == name
	}
}

// DeviceMatcherByID uses the device's ID to find a device.
func DeviceMatcherByID(id string) deviceMatcher {
	return func(device *DeviceInfo) bool {
		if device == nil {
			return false
		}
		return device.ID == id
	}
}

type DeviceInfo struct {
	ID   string
	Name string
}

func (d DeviceInfo) String() string {
	return d.Name + " [" + d.ID + "]"
}

type CapturingState struct {
	Muted  *bool
	Volume *float32
}

// IsMuted returns true if the mic is muted or the volume is 0 (nil values are not considered muted)
func (m *CapturingState) IsMuted() bool {
	if m == nil {
		return false
	}
	return (m.Muted != nil && *m.Muted) || (m.Volume != nil && *m.Volume == 0)
}

// IsUnmuted returns true if the mic is not muted or the volume is greater than 0 (nil values are not considered unmuted)
func (m *CapturingState) IsUnmuted() bool {
	if m == nil {
		return false
	}
	return (m.Muted != nil && !*m.Muted) || (m.Volume != nil && *m.Volume > 0)
}

// IsUndetermined returns true if the mic state is undetermined (nil values)
func (m *CapturingState) IsUndetermined() bool {
	return m == nil || (m.Muted == nil && m.Volume == nil)
}

type ListenState struct {
	Enabled   bool
	SavePower bool
	DeviceID  string
}

func (s *ListenState) String() string {
	return fmt.Sprintf("Listen state: %v, Save power: %v, Device ID: %s", s.Enabled, s.SavePower, s.DeviceID)
}
