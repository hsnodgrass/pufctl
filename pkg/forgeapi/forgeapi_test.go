package forgeapi

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hsnodgrass/pufctl/pkg/forgeapi/utils/mocks"
	"github.com/hsnodgrass/pufctl/pkg/forgeapi/utils/mocks/responses"
)

var (
	modJSONBody200 io.ReadCloser
	modResponse200 *http.Response
	fakeUAStr      = "test/0.0.0"
	fakeURLStr     = "https://fakeforge.com"
)

func init() {
	Client = &mocks.MockClient{}
	modJSONBody200 = ioutil.NopCloser(bytes.NewReader([]byte(responses.ModuleJSONBody200)))
	modResponse200 = &http.Response{
		StatusCode: 200,
		Body:       modJSONBody200,
	}
}

func TestGetRequest(t *testing.T) {
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return modResponse200, nil
	}
	response, err := GetRequest(fakeURLStr, fakeUAStr, Client)
	if err != nil {
		t.Errorf("Failed to get response with error: %w", err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Status code of response is not 200")
	}
}

func TestDecodeModule(t *testing.T) {
	mod, err := decodeModule(modResponse200)
	if err != nil {
		t.Errorf("Decoding module failed with error: %w", err)
	}
	if mod.Name != "apache" {
		t.Errorf("Failed to parse Module.Name")
	}
	if mod.Owner.Username != "puppetlabs" {
		t.Errorf("Failed to parse Module.Owner.Username")
	}
	if mod.CurrentRelease.Version != "4.0.0" {
		t.Errorf("Failed to parse Module.CurrentRelease.Version")
	}
	if mod.CurrentRelease.Module.Name != "apache" {
		t.Errorf("Failed to parse Module.CurrentRelease.Module.Name")
	}
	if len(mod.CurrentRelease.Tasks) != 1 {
		t.Errorf("Number of parsed tasks (%d) is not what was expected (%d)", len(mod.CurrentRelease.Tasks), 1)
	}
	if len(mod.Releases) != 1 {
		t.Errorf("Number of releases (%d) is not what was expected (%d)", len(mod.Releases), 1)
	}
}

func TestUserAgent(t *testing.T) {
	k, v := userAgent(fakeUAStr)
	if k != "User-Agent" || v != fakeUAStr {
		t.Errorf("Failed to retrieve User-Agent header pair")
	}
}

func TestForgeURL(t *testing.T) {
	ipv4URL, custom4 := forgeURL("ipv4")
	if ipv4URL != ForgeURLIPv4Only {
		t.Errorf("Failed to return proper URL for IPv4 Forge URL")
	}
	if custom4 {
		t.Errorf("Failed to return false for custom bool IPv4 Forge URL")
	}
	ipv6URL, custom6 := forgeURL("ipv6")
	if ipv6URL != ForgeURL {
		t.Errorf("Failed to return proper URL for IPv6 Forge URL")
	}
	if custom6 {
		t.Errorf("Failed to return false for custom bool IPv6 Forge URL")
	}
	customURL, customCustom := forgeURL(fakeURLStr)
	if customURL != fakeURLStr {
		t.Errorf("Failed to return proper URL for Custom Forge URL")
	}
	if !customCustom {
		t.Errorf("Failed to return true for custom bool Custom Forge URL")
	}
}

func TestRequestBaseURL(t *testing.T) {
	endpoints := []string{"users", "modules", "releases"}
	definedURLs := []string{ForgeURLIPv4Only, ForgeURL}
	for _, e := range endpoints {
		for _, d := range definedURLs {
			u, err := requestBaseURL(d, e)
			if err != nil {
				t.Errorf("Failed to get API URL for endpoint %s and URL input %s with error: %w", e, d, err)
			}
			switch e {
			case "users":
				retURL := fmt.Sprintf("%s%s", d, V3UsersEndpoint)
				if u != retURL {
					t.Errorf("Expected %s, got %s", retURL, u)
				}
			case "modules":
				retURL := fmt.Sprintf("%s%s", d, V3ModulesEndpoint)
				if u != retURL {
					t.Errorf("Expected %s, got %s", retURL, u)
				}
			case "releases":
				retURL := fmt.Sprintf("%s%s", d, V3ReleasesEndpoint)
				if u != retURL {
					t.Errorf("Expected %s, got %s", retURL, u)
				}
			}
		}
	}
	u, err := requestBaseURL(fakeURLStr, "")
	if err != nil {
		t.Errorf("Custom base URL failed with error: %w", err)
	}
	if u != fakeURLStr {
		t.Errorf("Expected %s, got %s", fakeURLStr, u)
	}
	_, err = requestBaseURL(ForgeURL, "fake")
	if err == nil {
		t.Errorf("Invalid URL input did not return an error")
	}
}
