package forgesource

import (
	"fmt"
	"sync"

	"github.com/hsnodgrass/pufctl/internal/logging"
	"github.com/hsnodgrass/pufctl/pkg/forgeapi"
	"github.com/hsnodgrass/pufctl/pkg/puppetfileparser/ast"
)

// GetModule returns a Puppet Forge Module object
func GetModule(nameslug, url, agent string) (forgeapi.Module, error) {
	logging.Debugln("Fetching module", nameslug)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	outMod, outErr := forgeapi.FetchModule(nameslug, url, agent, &waitGroup)
	waitGroup.Wait()
	return outMod, outErr
}

// GetModuleDependencies returns Module objects for specified dependencies in the given module
func GetModuleDependencies(mod forgeapi.Module, forge, agent string) ([]forgeapi.Module, error) {
	mods, errs := forgeapi.FetchModuleDependencies(mod, forge, agent)
	if len(errs) > 0 {
		for _, e := range errs {
			logging.Warnf("Failed to get a dependency: %#v", e)
		}
	}
	logging.Debugln("Successfully fetched module dependencies")
	return mods, nil
}

// ResolveDeps resolves Forge module dependencies
func ResolveDeps(mod forgeapi.Module, url, agent string, puppetfile *ast.Puppetfile) bool {
	changes := false
	deps, err := GetModuleDependencies(mod, url, agent)
	if err != nil {
		logging.Errorln("Failed to get module dependencies with error: ", err)
	}
	for _, d := range deps {
		modProp := []string{d.CurrentRelease.Version}
		err = puppetfile.AddModule(d.Slug, modProp)
		if err != nil {
			logging.Warnln("Failed to add module to Puppetfile with error: ", err)
		} else {
			changes = true
			autoDep := fmt.Sprintf("Added as dependency of %s", mod.Slug)
			err = puppetfile.AddModuleMetadata(d.Slug, "autodep", autoDep)
			logging.Infoln("Added to Puppetfile: ", d.Slug)
		}
	}
	return changes
}

// Search implements searching for modules via the Puppet Forge
func Search(url, agent string, opts forgeapi.ListModulesOpts) ([]forgeapi.Module, error) {
	result, err := forgeapi.ListModules(url, agent, opts)
	if err != nil {
		return []forgeapi.Module{}, fmt.Errorf("Search failed: %w", err)
	}
	return result.Results, nil
}
