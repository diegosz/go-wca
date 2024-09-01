# CHANGELOG

## v0.5.0

- fix: set default device property with NewEmptyPropVariant
- feat: add NewEmptyPropVariant
- feat: add ListenMicrophone example
- chore: fix makefile clean
- feat: add NewStringPropVariant and NewBoolPropVariant
- feat: implement SetValue
- feat: add PKEY_Listen_Setting
- tidy GUIDs
- feat: add IPolicyConfigVista and SetDefaultDeviceByID example
- feat: add ObserveMicrophones example WIP
- feat: add ToggleMicrophone example
- chore: lint
- feat: add MuteMicrophone example
- update go version requirements to 1.23
- fix: Fix BIsOffload type in AudioClientProperties struct
- fix: VTable IAudioMeterInformationVtbl struct order
- feat: Add constants for speaker flag positions
- chore: fix call to SyscallN
- chore: change module to match github path
- feat: Add support for initializing audio client with extended format
- chore: upgrade go to version 1.23 and update deps
- chore: change makefile build directory
- Fork from <https://github.com/arkadijs/go-wca>, based on PR: <https://github.com/moutend/go-wca/pull/21>

## v0.4.0

- IMMDeviceEnumerator.GetDevice impl
- IMMNotificationClient.OnDeviceStateChanged map dwNewState proper
- Proper Make deps
- Shared util function
- IAudioSessionNotification, IAudioSessionEvents example
- Bump Readme
- Sync type stubs
- Tidy build

## v0.3.0

- Merge pull request #16 from moutend/update-all-go-mod-files-under-example-directory
- Update all go.mod files under example directory
- Merge pull request #15 from moutend/remove-circle-ci-status-badge
- Remove CircleCI status badge
- Merge pull request #13 from moutend/update-go-ole-to-v1.2.6
- Update go-ole to v1.2.6
- Apply go fmt
- Merge pull request #9 from ThiefMaster/fix-ptr
- Merge pull request #10 from emarj/fix-typo-constants
- Merge pull request #11 from DarkMetalMouse/AudioMeterInformation
- implemented other functions
- updated README
- implemented IAudioMeterInformation
- fixing typo in AudioSessionState constants
- Use uint64 for pointer return vars

## v0.2.0

- Support device events (IMMNotificationClient)
- Support IMMNotificationClient registration.
- Add `_example/DeviceEvents`.
- Change project layout.
- Use go.mod to manage dependencies.
- Rename `.go` files.

The golang convention doesn't allow the capitalized file name (e.g. `IMMNotificationClient.go`). However, the lower-cased file name (e.g. `immnotificationclient.go`) is very hard to read, I decided to break the convention.

## 0.1.1

- Update README
- Update README
- Add LICENSE
- Replace executable names to `ls example`
- Use REFERENCE_TIME type instead of int64
- Add REFERENCE_TIME type
- Add example of exclusive timer driven mode rendering
- Add example which uses IAudioClient3 interface
- Use go-wav
- Use go-wave
- Use go 1.8.2 on CircleCI

## 0.1.0

Initial version
