// doc.go
// Package hemlis provides a caching wrapper around the Bitwarden Secrets Manager SDK.
//
// It provides thread-safe access to secrets with automatic caching and friendly name lookups.
//
// cfg := hemlis.Config{
//  AccessToken: os.Getenv("BWS_ACCESS_TOKEN"),
//  OrganizationID: os.Getenv("BWS_ORGANIZATION_ID"),
//  IdentityURL: os.Getenv("BWS_IDENTITY_URL"),
//  APIURL: os.Getenv("BWS_API_URL"),
//  CacheDuration: 15 * time.Minute,
// }
//
// manager, err := hemlis.New(cfg)
// if err != nil {
//  log.Fatal().Err(err).Msg("Failed to create secrets manager")
// }
//
// secret, err := manager.GetSecretByName("my-secret")
package hemlis
