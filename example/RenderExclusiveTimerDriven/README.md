# Rendering audio with exclusive timer driven mode

This example shows that the rendering audio with exclusive timer driven mode.

## Prerequisites

- Go 1.23 or later

## Build

```console
go build
```

That's it. Then you'll get `RenderExclusiveTimerDriven.exe`. Note that your platform is not Windows, you need set the environment variable `GOOS='windows'`.

## Usage

```console
./RenderExclusiveTimerDriven -i music.wav
```

Please specify the WAVE audio file with `-i` or `--input` flag.

## Note

Whether the wAV file can play or not is up to hardware settings. You need specify the WAV file which is encoded as same as output settings.

For example, when the output is configured to comibination of 96 kHz and 24 bit, you can only play WAV files which are encoded with 96 kHz and 24 bit.

The Windows Core Audio API was introduced Windows vista, so that the later than that version of Windows could run this example. However, I'm not sure because I've just tested this example on Windows 10 version 1607 at the moment. Operation verification including bug report are welcome.
