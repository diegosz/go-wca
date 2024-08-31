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

var version = "latest"
var revision = "latest"

type DeviceFlag struct {
	Value string
	IsSet bool
}

func (f *DeviceFlag) Set(value string) (err error) {
	if value == "" {
		err = fmt.Errorf("a device name is required")
		return
	}
	f.Value = value
	f.IsSet = true
	return
}

func (f *DeviceFlag) String() string {
	return f.Value
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
	f.Var(&deviceFlag, "device", "Specify the device display name (required)")
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
		err = fmt.Errorf("a device name is required")
		return
	}
	return toggleMicrophone(deviceFlag.Value)
}

type micDevice struct {
	id   string
	name string
}

type micState struct {
	muted  *bool
	volume *float32
}

func (m *micState) IsEqual(other *micState) bool {
	if m == nil && other == nil { // both are nil, which is equal
		return true
	}

	// one is nil, but not both
	if m == nil || other == nil {
		return false
	}

	// one has nil muted, but the other doesn't
	if m.muted == nil && other.muted != nil {
		return false
	}
	if m.muted != nil && other.muted == nil {
		return false
	}

	// one has nil volume, but the other doesn't
	if m.volume == nil && other.volume != nil {
		return false
	}
	if m.volume != nil && other.volume == nil {
		return false
	}

	// both are not nil, so compare the values and return the result
	return *m.muted == *other.muted && *m.volume == *other.volume
}

// IsMuted returns true if the mic is muted or the volume is 0 (nil values are not considered muted)
func (m *micState) IsMuted() bool {
	if m == nil {
		return false
	}
	return (m.muted != nil && *m.muted) || (m.volume != nil && *m.volume == 0)
}

// IsUnmuted returns true if the mic is not muted or the volume is greater than 0 (nil values are not considered unmuted)
func (m *micState) IsUnmuted() bool {
	if m == nil {
		return false
	}
	return (m.muted != nil && !*m.muted) || (m.volume != nil && *m.volume > 0)
}

// IsUndetermined returns true if the mic state is undetermined (nil values)
func (m *micState) IsUndetermined() bool {
	return m == nil || (m.muted == nil && m.volume == nil)
}

func toggleMicrophone(name string) error {
	if name == "" {
		return nil
	}
	// The device matcher logic, which simply uses the device's display name as
	// shown in the Windows UI, but there's no reason you couldn't use
	// completely different matching logic.
	deviceMatcherByName := func(name string) func(micDevice) bool {
		return func(device micDevice) bool {
			return device.name == name
		}
	}
	// This performs the actual mic state toggle, which is used in.
	toggleMicMute := func(device micDevice, state micState) *micState {
		fmt.Printf("Toggling mute for: %s\n", device.name)
		if state.muted != nil {
			// fmt.Printf("Current mute state: %v\n", *state.muted)
			wantState := !*state.muted
			// fmt.Printf("Wanted mute state: %v\n", wantState)
			return &micState{muted: &wantState}
		}
		return nil
	}
	// How the `wcaECaptureDeviceInfo` function is used to locate the mic,
	// including the matcher function and the callback that performs the toggle.
	return wcaECaptureDeviceInfo(deviceMatcherByName(name), toggleMicMute)
}

func wcaECaptureDeviceInfo(matcher func(micDevice) bool, callback func(micDevice, micState) *micState) error {
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

	// Enumerate audio endpoint devices.
	var deviceCollection *wca.IMMDeviceCollection
	if err := enumerator.EnumAudioEndpoints(wca.ECapture, wca.DEVICE_STATE_ACTIVE, &deviceCollection); err != nil {
		return fmt.Errorf("failed to enumerate audio endpoints: %v", err)
	}
	defer deviceCollection.Release()

	// Get device count.
	var count uint32
	if err := deviceCollection.GetCount(&count); err != nil {
		return fmt.Errorf("failed to get device count: %v", err)
	}

	// Iterate over devices.
	for i := uint32(0); i < count; i++ {
		var item *wca.IMMDevice
		if err := deviceCollection.Item(i, &item); err != nil {
			return fmt.Errorf("failed to get device item: %v", err)
		}
		defer item.Release()

		// Get the device ID.
		var deviceId string
		if err := item.GetId(&deviceId); err != nil {
			return fmt.Errorf("failed to get device ID: %v", err)
		}

		// Open the property store.
		var propertyStore *wca.IPropertyStore
		if err := item.OpenPropertyStore(wca.STGM_READ, &propertyStore); err != nil {
			return fmt.Errorf("failed to open property store: %v", err)
		}
		defer propertyStore.Release()

		// Get the friendly name of the device.
		var propVariant wca.PROPVARIANT
		if err := propertyStore.GetValue(&wca.PKEY_Device_FriendlyName, &propVariant); err != nil {
			return fmt.Errorf("failed to get friendly name: %v", err)
		}
		deviceName := propVariant.String()
		defer func() {
			_ = propVariant.Clear()
		}()

		// Check if the device is the one we're looking for.
		if matcher(micDevice{id: deviceId, name: deviceName}) {

			// Get the IAudioEndpointVolume interface.
			var aev *wca.IAudioEndpointVolume
			if err := item.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev); err != nil {
				return fmt.Errorf("failed to activate audio endpoint volume: %v", err)
			}
			defer aev.Release()

			getMicState := func(aev *wca.IAudioEndpointVolume) *micState {
				// Build the micState object, starting with nil values to
				// indicate undetermined state.
				state := &micState{}
				// Get the mute state.
				for j := 0; j < 5; j++ {
					var muted bool
					if err := aev.GetMute(&muted); err != nil {
						fmt.Printf("failed to get mute state: %v", err)
					} else {
						state.muted = &muted
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
				// Get the volume level.
				for j := 0; j < 5; j++ {
					var volume float32
					if err := aev.GetMasterVolumeLevelScalar(&volume); err != nil {
						fmt.Printf("failed to get volume level: %v", err)
					} else {
						state.volume = &volume
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
				return state
			}

			setMicState := func(aev *wca.IAudioEndpointVolume, state *micState) {
				if state == nil {
					return
				}
				// Existing state.
				currentState := getMicState(aev)
				// Set the mute state, if necessary.
				if state.muted != nil && currentState.muted != state.muted {
					fmt.Printf("Setting mute state to %v\n", *state.muted)
					if err := aev.SetMute(*state.muted, nil); err != nil {
						fmt.Printf("Failed to set mute state: %v\n", err)
					}
				}
				// Set the volume level, if necessary.
				if state.volume != nil && currentState.volume != state.volume {
					fmt.Printf("Setting volume level to %v\n", *state.volume)
					if err := aev.SetMasterVolumeLevelScalar(*state.volume, nil); err != nil {
						fmt.Printf("Failed to set volume level: %v\n", err)
					}
				}
			}
			// Get the initial state.
			state := getMicState(aev)
			// Call the callback, if provided.
			if callback != nil {
				dev := micDevice{
					id:   deviceId,
					name: deviceName,
				}
				newState := callback(dev, *state)
				// Set the new state, if provided.
				if newState != nil {
					// The callback asked for a change, so update the state.
					setMicState(aev, newState)
				}
			}
		}
	}
	return nil
}
