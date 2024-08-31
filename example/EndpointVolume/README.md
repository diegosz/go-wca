# Endpoint volume

This example shows that changing volume for default rendering device.

## Prerequisites

- Go 1.23 or later

## Build

```console
go build
```

That's it. Then you'll get `EndpointVolume.exe`. Note that your platform is not Windows, you need set the environment variable `GOOS='windows'`.

## Usage

```console
./EndpointVolume --volume 0.1
```

Available flags are:

- `-v, --volume` sets the volume as scalar value
- `-g, --gain` sets the volume as level (dB) value
- `-m, --mute` sets mute state

The Windows Core Audio API was introduced Windows vista, so that the later than that version of Windows could run this example. However, I'm not sure because I've just tested this example on Windows 10 version 1607 at the moment. Operation verification including bug report are welcome.
