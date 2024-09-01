# Dev Notes

## Core Audio APIs

The core audio APIs provide the means for audio applications to access audio endpoint devices such as headphones and microphones.
The core audio APIs serve as the foundation for higher-level audio APIs such as Microsoft DirectSound and the Windows multimedia waveXxx functions.
Most applications communicate with the higher-level APIs, but some applications with special requirements might need to communicate directly with the core audio APIs.

## Audio Property System

The Windows Property System provides common interfaces for accessing metadata referencing various areas of the Windows operating system.
These areas may be devices, windows, the file system, and more.
The Audio Property System is a subset of the Windows Property System that provides access to audio-related settings.

Apparently, `IPropertyStore::SetValue` takes care of notifying to the `AudioSrv` and `IPropertyStore::Commit` is seemingly not required for events to fire.

## Access Denied

When running an example that needs to change a value in the `IPropertyStore`, it has open it with `STGM_READ_WRITE` storage mode flag.
You will encounter an "Access Denied" error if you run the example without elevated privileges.

You need to either elevate privileges to at least `Power User` or run your executable `As Administrator` or disable `UAC`.
An alternative solution is to change the permissions on the registry key.
Open the key using `Regedit`, right click and select `Permissions`.
Either add the profile that your application runs under (ie a `service account`) or select an existing group (ie `Users`) and grant them full access.
This way you're not having to grant elevated privileges which is always a security concern.

## Resources

- [c# - Reverse engineer Listen to this device on Windows - Reverse Engineer Windows When There's No API - Stack Overflow](https://stackoverflow.com/questions/57929296/reverse-engineer-listen-to-this-device-on-windows-reverse-engineer-windows-whe)
- [dll - Reverse engineer Listen to this device on Windows - Reverse engineer Windows when there's no API - Reverse Engineering Stack Exchange](https://reverseengineering.stackexchange.com/questions/23454/reverse-engineer-listen-to-this-device-on-windows-reverse-engineer-windows-whe/31865#31865)
- [Core Audio APIs - Win32 apps | Microsoft LearnCalifornia Consumer Privacy Act (CCPA) Opt-Out IconCalifornia Consumer Privacy Act (CCPA) Opt-Out Icon](https://learn.microsoft.com/en-us/windows/win32/coreaudio/core-audio-apis-in-windows-vista)
- [Device Properties (Core Audio APIs) - Win32 apps | Microsoft LearnCalifornia Consumer Privacy Act (CCPA) Opt-Out IconCalifornia Consumer Privacy Act (CCPA) Opt-Out Icon](https://learn.microsoft.com/en-us/windows/win32/coreaudio/device-properties?redirectedfrom=MSDN)
- [IAudioSystemEffectsPropertyStore - Win32 apps | Microsoft LearnCalifornia Consumer Privacy Act (CCPA) Opt-Out IconCalifornia Consumer Privacy Act (CCPA) Opt-Out Icon](https://learn.microsoft.com/en-us/windows/win32/api/mmdeviceapi/nn-mmdeviceapi-iaudiosystemeffectspropertystore)
- [IMMDevice (mmdeviceapi.h) - Win32 apps | Microsoft LearnCalifornia Consumer Privacy Act (CCPA) Opt-Out IconCalifornia Consumer Privacy Act (CCPA) Opt-Out Icon](https://learn.microsoft.com/en-us/windows/win32/api/mmdeviceapi/nn-mmdeviceapi-immdevice)
- [IAudioSystemEffectsPropertyStore::OpenUserPropertyStore - Win32 apps | Microsoft LearnCalifornia Consumer Privacy Act (CCPA) Opt-Out IconCalifornia Consumer Privacy Act (CCPA) Opt-Out Icon](https://learn.microsoft.com/en-us/windows/win32/api/mmdeviceapi/nf-mmdeviceapi-iaudiosystemeffectspropertystore-openuserpropertystore)
- [PROPVARIANT (propidlbase.h) - Win32 apps | Microsoft LearnCalifornia Consumer Privacy Act (CCPA) Opt-Out IconCalifornia Consumer Privacy Act (CCPA) Opt-Out Icon](https://learn.microsoft.com/en-us/windows/win32/api/propidlbase/ns-propidlbase-propvariant?redirectedfrom=MSDN)
