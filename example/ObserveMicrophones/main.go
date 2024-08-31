package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/diegosz/go-wca/pkg/wca"
	"github.com/go-ole/go-ole"
)

var version = "latest"
var revision = "latest"

type IntervalFlag struct {
	Value uint16
	IsSet bool
}

func (f *IntervalFlag) Set(value string) (err error) {
	if value == "" {
		value = "1"
	}
	var v uint64
	if v, err = strconv.ParseUint(value, 10, 16); err != nil {
		return
	}
	f.Value = uint16(v)
	f.IsSet = true
	return
}

func (f *IntervalFlag) String() string {
	return fmt.Sprintf("%d", f.Value)
}

func main() {
	fmt.Println("WIP: This example is a work in progress")
	var err error
	if err = run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) (err error) {
	var intervalFlag IntervalFlag
	var versionFlag bool

	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.Var(&intervalFlag, "interval", "Specify the observe interval in seconds (default to 1)")
	f.Var(&intervalFlag, "i", "Alias of --interval")
	f.BoolVar(&versionFlag, "version", false, "Show version")
	if err = f.Parse(args[1:]); err != nil {
		return
	}
	if versionFlag {
		fmt.Printf("%s-%s\n", version, revision)
		return
	}
	i := 1 * time.Second
	if intervalFlag.IsSet {
		i = time.Second * time.Duration(intervalFlag.Value)
	}
	fmt.Println("Press Ctrl+C to exit")
	observeMicrophones(i)
	return nil
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

func observeMicrophones(interval time.Duration) {
	if interval == 0 {
		return
	}
	fmt.Printf("Observing microphone interval: %s\n", interval)

	done := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go observeMicState(ctx, interval, func(devices []micDevice, states []micState) {
	}, done)
	go func() {
		<-signalChan
		fmt.Println("Interrupted by SIGINT")
		cancel()
	}()
	<-done
}

func observeMicState(ctx context.Context, interval time.Duration, callback func([]micDevice, []micState), done chan struct{}) {
	defer close(done)

	deviceMatcherAll := func() func(micDevice) bool {
		return func(_ micDevice) bool {
			return true // Match all devices.
		}
	}

	captureMicState := func(matchedDevices *[]micDevice, capturedStates *[]micState) func(micDevice, micState) *micState {
		return func(device micDevice, state micState) *micState {
			*matchedDevices = append(*matchedDevices, device)
			*capturedStates = append(*capturedStates, state)
			return nil // No changes requested, so return nil.
		}
	}

	captureDeviceInfoOnce := func() ([]micDevice, []micState, error) {
		states := make([]micState, 0)
		devices := make([]micDevice, 0)
		err := wcaECaptureDeviceInfo(deviceMatcherAll(), captureMicState(&devices, &states))
		if err != nil {
			fmt.Printf("Failed to get state: %v\n", err)
		}
		return devices, states, err
	}

	devices, states, err := captureDeviceInfoOnce()
	if err != nil {
		fmt.Printf("Failed to get initial state: %v\n", err)
	}
	if len(devices) == 0 || len(states) == 0 {
		fmt.Printf("No devices found\n")
	} else {
		callback(devices, states)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			devicesNow, statesNow, err := captureDeviceInfoOnce()
			if err != nil {
				fmt.Printf("Failed to get state: %v\n", err)
			}
			if len(devicesNow) != len(statesNow) {
				fmt.Printf("Mismatched devices and states: %d != %d\n", len(devicesNow), len(statesNow))
				continue
			}

			// Turn the devices & states into maps for easier comparison.
			historyDevices := make(map[string]micDevice)
			for i := 0; i < len(devices); i++ {
				historyDevices[devices[i].id] = devices[i]
			}
			historyStates := make(map[string]micState)
			for i := 0; i < len(states); i++ {
				historyStates[devices[i].id] = states[i]
			}

			// Turn the devicesNow & statesNow into maps for easier comparison.
			devicesNowMap := make(map[string]micDevice)
			for i := 0; i < len(devicesNow); i++ {
				devicesNowMap[devicesNow[i].id] = devicesNow[i]
			}
			statesNowMap := make(map[string]micState)
			for i := 0; i < len(statesNow); i++ {
				statesNowMap[devicesNow[i].id] = statesNow[i]
			}

			// Find the devices and states that have changed.
			changedDevices := make([]micDevice, 0)
			changedStates := make([]micState, 0)

			// Check if the state has changed.
			for i := 0; i < len(devicesNow); i++ {
				historyState, ok := historyStates[devicesNow[i].id]
				if (ok && !statesNow[i].IsEqual(&historyState)) || !ok {
					// State has changed or is new.
					fmt.Printf("State changed for %s\n", devicesNow[i].name)
					historyStates[devicesNow[i].id] = statesNow[i]
					changedDevices = append(changedDevices, devicesNow[i])
					changedStates = append(changedStates, statesNow[i])
				}
			}

			// Check if a device has been removed.
			for id, device := range historyDevices {
				if _, ok := devicesNowMap[id]; !ok {
					// Device has been removed.
					fmt.Printf("Device removed: %s\n", device.name)
					delete(historyDevices, id)
					delete(historyStates, id)
					changedDevices = append(changedDevices, device)
					changedStates = append(changedStates, micState{})
				}
			}

			// check if a device has been added/removed
			if len(changedDevices) > 0 {
				// Update the history from the map.
				devices = make([]micDevice, 0)
				states = make([]micState, 0)
				for id, state := range historyStates {
					devices = append(devices, historyDevices[id])
					states = append(states, state)
				}
				// Call the callback with the changed devices and states.
				callback(changedDevices, changedStates)
			}
		}
	}
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
