# Loopback audio capturing with shared event driven mode

This example shows that the loopback (a.k.a. what you hear) capturing with shared event driven mode.

## Prerequisites

- Go 1.13 or later

## Build

```console
go build
```

That's it. Then you'll get `LoopbackCaptureSharedEventDriven.exe`. Note that your platform is not Windows, you need set the environment variable `GOOS='windows'`.

## Usage

```console
./LoopbackCaptureSharedEventDriven -o music.wav -d 10
```

Please specify the flag `-o` or `--output` for saving audio file. The `-d` or `--duration` is optional and it indicates recording duration in second. If the recording duration was not specified, it keeps recording until receiving interruption by Ctrl-C.

## Known issues

According to [the documentation about loopback audio capturing](https://msdn.microsoft.com/en-us/library/windows/desktop/dd316551(v=vs.85).aspx), the audio capture client can set the audio ready event without any errors, but it will never fire.

> A pull-mode capture client does not receive any events when a stream is initialized with event-driven buffering and is loopback-enabled.
> 
> To work around this, initialize a render stream in event-driven mode. Each time the client receives an event for the render stream, it must signal the capture client to run the capture thread that reads the next set of samples from the capture endpoint buffer.

I don't know why this issue has not been fixed yet, but we need add the workaround for current implementation.

## Note

As with capturing from microphone, stability of the event driven mode is much lower than the timer driven mode. Because we cannot control the scheduling of goroutines, when the goroutine which observing audio ready event was stopped, stutterring occurs.

The Windows Core Audio API was introduced Windows vista, so that the later than that version of Windows could run this example. However, I'm not sure because I've just tested this example on Windows 10 version 1607 at the moment. Operation verification including bug report are welcome.
