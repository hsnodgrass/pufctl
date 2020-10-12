// Package forgeapi provides a (incomplete) API for the Puppet Forge
package forgeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	qstring "github.com/google/go-querystring/query"
)

var (
	// Client implements the HTTPClient interface for making HTTP rquests
	Client HTTPClient

	transport = &http.Transport{
		MaxIdleConns:       3,
		IdleConnTimeout:    5 * time.Second,
		DisableCompression: true,
		Proxy:              http.ProxyFromEnvironment,
	}
)

// HTTPClient provides an interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func init() {
	Client = &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}
}

func userAgent(agent string) (string, string) {
	return "User-Agent", fmt.Sprintf("%s", agent)
}

func forgeURL(url string) (string, bool) {
	switch strings.ToLower(url) {
	case "ipv4", ForgeURLIPv4Only:
		return ForgeURLIPv4Only, false
	case "default", "ipv6", ForgeURL:
		return ForgeURL, false
	default:
		return url, true
	}
}

func requestBaseURL(url, endpoint string) (string, error) {
	baseURL, custom := forgeURL(url)
	if custom {
		return baseURL, nil
	}
	switch strings.ToLower(endpoint) {
	case "users":
		return fmt.Sprintf("%s%s", url, V3UsersEndpoint), nil
	case "modules":
		return fmt.Sprintf("%s%s", url, V3ModulesEndpoint), nil
	case "releases":
		return fmt.Sprintf("%s%s", url, V3ReleasesEndpoint), nil
	default:
		return "", fmt.Errorf("Specified endpoint %s is not a valid Forge API endpoint option", endpoint)
	}
}

func decodeModule(resp *http.Response) (Module, error) {
	var mod Module
	err := json.NewDecoder(resp.Body).Decode(&mod)
	if err != nil {
		return mod, &JSONDecodeError{Err: err}
	}
	return mod, nil
}

// GetRequest makes an HTTP Get request and returns a response
// on an HTTP response code of 200, or an error otherwise.
func GetRequest(url, agent string, client HTTPClient) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, &GetError{Err: err, URL: url, Headers: map[string]string{}}
	}
	header, value := userAgent(agent)
	req.Header.Add(header, value)
	resp, err := client.Do(req)
	if err != nil {
		hdrs := map[string]string{header: value}
		return nil, &GetError{Err: err, URL: url, Headers: hdrs}
	}
	if resp.StatusCode != 200 {
		return nil, &GetNon200Error{URL: url, StatusCode: resp.StatusCode}
	}
	return resp, nil
}

// FetchModule performs a Forge API get request for the named module and sends the result to the specified channel
func FetchModule(nameslug, url, agent string, wg *sync.WaitGroup) (Module, error) {
	defer wg.Done()
	baseURL, err := requestBaseURL(url, "modules")
	if err != nil {
		return Module{}, &FetchError{Err: err}
	}
	finalURL := fmt.Sprintf("%s/%s", baseURL, nameslug)
	resp, err := GetRequest(finalURL, agent, Client)
	if err != nil {
		return Module{}, &FetchError{Err: err}
	}
	mod, err := decodeModule(resp)
	if err != nil {
		return Module{}, &FetchError{Err: err}
	}
	return mod, nil
}

// FetchModuleDependencies returns a slice of Modules that are marked as dependencies of the given module
func FetchModuleDependencies(mod Module, url, agent string) ([]Module, []error) {
	var waitGroup sync.WaitGroup
	var deps []Module
	var errs []error
	waitGroup.Add(len(mod.CurrentRelease.Metadata.Dependencies))
	for _, d := range mod.CurrentRelease.Metadata.Dependencies {
		go func(name, url, agent string, wg *sync.WaitGroup) {
			_name := strings.Replace(name, "/", "-", 1)
			mod, err := FetchModule(_name, url, agent, &waitGroup)
			if err != nil {
				errs = append(errs, err)
			}
			deps = append(deps, mod)
		}(d.Name, url, agent, &waitGroup)
	}
	waitGroup.Wait()
	return deps, errs
}

// ListModules returns an array of Modules based on the search
// criterial outlined in the opts param.
func ListModules(url, agent string, opts ListModulesOpts) (ListResults, error) {
	var listRes ListResults
	baseURL, err := requestBaseURL(url, "modules")
	if err != nil {
		return listRes, fmt.Errorf("forgeapi: List failed: %w", err)
	}
	v, err := qstring.Values(opts)
	if err != nil {
		return listRes, fmt.Errorf("forgeapi: List failed: %w", err)
	}
	finalURL := fmt.Sprintf("%s?%s", baseURL, v.Encode())
	resp, err := GetRequest(finalURL, agent, Client)
	if err != nil {
		return listRes, fmt.Errorf("forgeapi: List failed: %w", err)
	}
	err = json.NewDecoder(resp.Body).Decode(&listRes)
	if err != nil {
		return listRes, fmt.Errorf("forgeapi: List failed: %w", err)
	}
	return listRes, nil
}
