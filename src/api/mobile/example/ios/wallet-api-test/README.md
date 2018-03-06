# Build mobile.framework

```bash
gomobile bind -target=ios github.com/skycoin/skycoin-exchange/src/api/mobile
```

This will generate `Mobile.framework` file, then open the wallet-api-test.xcodeproj file with xcode, the `Mobile.framework` will be loaded automatically.
