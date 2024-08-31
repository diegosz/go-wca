# Observe Microphones

This example shows how to observe the capturing devices.

> [!WARNING]
> This example is a work in progress.
> It may not work as expected yet...

## Prerequisites

- Go 1.23 or later

## Build

```console
go build
```

That's it. Then you'll get `ObserveMicrophones.exe`. Note that your platform is not Windows, you need set the environment variable `GOOS='windows'`.

## Usage

```console
./ObserveMicrophones --interval 1
```

Available flags are:

- `-i, --interval` sets the observe interval in seconds

The Windows Core Audio API was introduced Windows vista, so that the later than that version of Windows could run this example. However, I'm not sure because I've just tested this example on Windows 10 version 1607 at the moment. Operation verification including bug report are welcome.
