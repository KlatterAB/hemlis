# Bitwarden Secrets Manager Go Wrapper - hemlis

A Go package that provides a caching wrapper around the Bitwarden Secrets Manager SDK.

## Features

- Thread-safe access to secrets
- Automatic caching with configurable duration
- Friendly name lookups for secrets
- Simple, idiomatic Go API

## Installation

```bash
go get github.com/klatterab/hemlis
```

## Usage

```go
import "github.com/klatterab/hemlis"

// Configuration
cfg := hemlis.Config{
    AccessToken:    os.Getenv("BITWARDEN_ACCESS_TOKEN"),
    OrganizationID: os.Getenv("BITWARDEN_ORGANIZATION_ID"),
    IdentityURL:    "https://identity.bitwarden.com",
    APIURL:         "https://api.bitwarden.com",
    CacheDuration:  15 * time.Minute,
}

// Create manager
manager, err := hemlis.New(cfg)
if err != nil {
    log.Fatal(err)
}

// Get a secret by its friendly name
secret, err := manager.GetSecretByName("my-secret")
if err != nil {
    log.Fatal(err)
}
```
