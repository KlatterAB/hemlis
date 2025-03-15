// doc.go
// Package hemlis provides a caching wrapper around the Bitwarden Secrets Manager SDK.
//
// It provides thread-safe access to secrets with automatic caching and friendly name lookups.
// Secret names can be normalized using the optional NameNormalization option to handle case inconsistencies.
//
// Basic usage:
// cfg := hemlis.Config{
//  AccessToken: os.Getenv("BWS_ACCESS_TOKEN"),
//  OrganizationID: os.Getenv("BWS_ORGANIZATION_ID"),
//  IdentityURL: os.Getenv("BWS_IDENTITY_URL"),
//  APIURL: os.Getenv("BWS_API_URL"),
//  CacheDuration: 15 * time.Minute,
// }
//
// With name normalization (optional):
// lowerCaseNorm := hemlis.LowerCase
// cfg := hemlis.Config{
//  AccessToken: os.Getenv("BWS_ACCESS_TOKEN"),
//  // ... other fields ...
//  NameNormalization: &lowerCaseNorm, // Makes name lookups case-insensitive
// }
//
// manager, err := hemlis.New(cfg)
// if err != nil {
//  log.Fatal().Err(err).Msg("Failed to create secrets manager")
// }
//
// // With normalization, will find the secret regardless of case (e.g., "MY-SECRET", "My-Secret")
// secret, err := manager.GetSecretByName("my-secret")
package hemlis
