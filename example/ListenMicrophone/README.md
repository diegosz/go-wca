# Listen Microphone

This example shows how to listen a capturing device by name.

## Prerequisites

- Go 1.23 or later

## Build

```console
go build
```

That's it. Then you'll get `ListenMicrophone.exe`. Note that your platform is not Windows, you need set the environment variable `GOOS='windows'`.

## Usage

```console
./ListenMicrophone --listen true --in "Microphone (Realtek(R) Audio)" --out "Speakers (Realtek(R) Audio)"
```

> [!NOTE]
> It has to run with the administrator privilege.

Available flags are:

- `--in` sets the input device name to listen
- `--out` sets the output device name
