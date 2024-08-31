# Set Default Device By ID

This example shows how to set the default device by ID.

## Prerequisites

- Go 1.23 or later

## Build

```console
go build
```

That's it. Then you'll get `SetDefaultDeviceByID.exe`. Note that your platform is not Windows, you need set the environment variable `GOOS='windows'`.

## Usage

```console
./SetDefaultDeviceByID --device "5eb93697-2840-43ab-ad5d-a8cb60cf92b1"
```

Available flags are:

- `-d, --device` define device ID

The Windows Core Audio API was introduced Windows vista, so that the later than that version of Windows could run this example. However, I'm not sure because I've just tested this example on Windows 10 version 1607 at the moment. Operation verification including bug report are welcome.
