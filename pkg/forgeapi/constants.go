// Package forgeapi provides a (incomplete) API for the Puppet Forge
package forgeapi

// ForgeURL is the default URL base for Puppet Forge
const ForgeURL = "https://forgeapi-cdn.puppet.com"

// ForgeURLIPv4Only is the secondary URL base for the Forge but does not accept IPv6 connections
const ForgeURLIPv4Only = "https://forgeapi.puppet.com"

// V3UsersEndpoint is the v3 Forge API endpoint for users
const V3UsersEndpoint = "/v3/users"

// V3ModulesEndpoint is the v3 Forge API endpoint for modules
const V3ModulesEndpoint = "/v3/modules"

// V3ReleasesEndpoint is the v3 Forge API endpoint for releases
const V3ReleasesEndpoint = "/v3/releases"
