# Mute Microphone

This example shows how to mute the default capturing device.

## Prerequisites

- Go 1.23 or later

## Build

```console
go build
```

That's it. Then you'll get `MuteMicrophone.exe`. Note that your platform is not Windows, you need set the environment variable `GOOS='windows'`.

## Usage

```console
./MuteMicrophone --mute true
```

Available flags are:

- `-m, --mute` sets mute state

The Windows Core Audio API was introduced Windows vista, so that the later than that version of Windows could run this example. However, I'm not sure because I've just tested this example on Windows 10 version 1607 at the moment. Operation verification including bug report are welcome.
