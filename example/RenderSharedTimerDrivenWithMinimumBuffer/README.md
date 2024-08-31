# Rendering audio with shared timer driven mode (using `IAudioClient3` interface)

This example shows that the rendering audio with shared timer driven mode.

FYI, `IAudioClient3` interface was introduced from Windows 10, which interface allows us to render audio with minimum buffer size. It's ideal for professional music production such as real-time monitoring.

## Prerequisites

- Go 1.23 or later
- Windows 10 v1607 or later

## Build

```console
go build
```

That's it. Then you'll get `RenderSharedTimerDrivenWithMinimumBuffer.exe`. Note that your platform is non-Windows, you need set the environment variable `GOOS='windows'`.

## Usage

```console
./RenderSharedTimerDrivenWithMinimumBuffer -i music.wav
```

Please specify the WAVE audio file with `-i` or `--input` flag.

## Note

This example doesn't convert the sample rate of the input. The rendering step will be failed when the sample rate and bit depth of input audio doesn't match system default sample rate and bit depth.

To avoid this error, please set the system default settings for shared mode playback according to the audio file you want to play.

1. Open control panel and select sound.
1. Select the main playback device and open property.
1. Select the advanced tab and check the default sample rate and bit depth.
1. Apply changes.

For example, if you want to play the wave audio file which is extracted from DVD, in other words, which is recorded with 48000 Hz / 16 bit, you need choose 4800 Hz / 16 bit for shared mode audio playback.

The Windows Core Audio API was introduced Windows vista, so that the later than that version of Windows could run this example. However, I'm not sure because I've just tested this example on Windows 10 version 1607 at the moment. Operation verification including bug report are welcome.
