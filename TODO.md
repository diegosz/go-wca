# TODO(s)

- [ ] Test audio devices on the CI server
    I test the `example` on my Windows 10 laptop at the moment.
    I'm looking for CI which can test the audio device.
    For example, Appveyer CI supports only Windows Server and the server edition doesn't support physical audio device.

- [ ] Figure out how to handle bizarre default bit depth
    `IAudioClient::GetMixFormat` returns always 32 bit as a bit depth on my machine (Macbook Air / Windows 10 version 1607). I'm investigating this is my machine specific issue or not.
- [ ] Complete ObserveMicrophones example
  - [ ] Change the logic for tracking changes
  - [ ] Show the current state of the devices and the actual changes
- Check if stringToPropVariantis this correct
    Maybe we should use ole.SysAllocString instead of UTF16PtrFromString
    <https://github.com/go-ole/go-ole/blob/master/com.go>
    <https://www.roblocher.com/whitepapers/oletypes.html>
