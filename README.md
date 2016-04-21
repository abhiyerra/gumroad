# gumroad

Golang API to verify that Gumroad License Keys are valid.


## Usage

```
import (
    "github.com/abhiyerra/gumroad"
)
```

```
err := gumroad.VerifyLicense(productPermalink, licenseKey, incrementUsesCount)
if err != nil {
     // Verification failed
}

// All good
```
