package gitsource

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-billy"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/hsnodgrass/pufctl/internal/auth"
	"github.com/hsnodgrass/pufctl/internal/logging"
	"github.com/hsnodgrass/pufctl/pkg/forgeapi"
	"github.com/hsnodgrass/pufctl/pkg/puppetfileparser/ast"
)

// CloneFile retrives a specific file from a Git repository
func CloneFile(url, filepath, branch string, modauth auth.Auth) (billy.File, string, error) {
	fs := memfs.New()
	auth := modauth.Method
	ref := plumbing.NewBranchReferenceName(branch)
	clone, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:           url,
		Progress:      os.Stdout,
		Auth:          auth,
		ReferenceName: ref,
	})
	if err != nil {
		return nil, "", fmt.Errorf("Failed to clone repository with error: %w", err)
	}
	head, err := clone.Head()
	if err != nil {
		return nil, "", fmt.Errorf("Failed to get HEAD of repository with error: %w", err)
	}
	headRef := strings.Split(string(head.Name()), "/")
	defBranch := headRef[len(headRef)-1]
	logging.Debugln("Cloned git repository into in-memory file system. Cloned branch: ", defBranch)
	f, err := fs.Open(filepath)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to open %s with error: %w", filepath, err)
	}
	logging.Debugln("Successfully opened", filepath)
	return f, defBranch, nil
}

// GetModuleMeta parses a module's git repo for the metadata.json file,
// unmarshalls it into a forgeapi.ModuleMetadata struct, and returns the struct
func GetModuleMeta(url string, modauth auth.Auth) (string, *forgeapi.ModuleMetadata) {
	var meta forgeapi.ModuleMetadata
	m, defBranch, err := CloneFile(url, "metadata.json", "master", modauth)
	if err != nil {
		logging.Errorf("Failed to clone file metadata.json with error: %w", err)
	}
	err = json.NewDecoder(m).Decode(&meta)
	if err != nil {
		logging.Errorln("Failed to decode module metadata with error: ", err)
	}
	logging.Debugln("Successfully decoded metadata.json")
	return defBranch, &meta
}

// ResolveDeps resolves Git module dependencies
func ResolveDeps(meta *forgeapi.ModuleMetadata, puppetfile *ast.Puppetfile) bool {
	changes := false
	metaSlug := strings.Replace(meta.Name, "/", "-", 1)
	for _, d := range meta.Dependencies {
		modProp := []string{":latest"}
		slug := strings.Replace(d.Name, "/", "-", 1)
		err := puppetfile.AddModule(slug, modProp)
		if err != nil {
			logging.Warnln("Failed to add module to Puppetfile with error: ", err)
		} else {
			changes = true
			autoDep := fmt.Sprintf("Added as dependency of %s", metaSlug)
			err = puppetfile.AddModuleMetadata(slug, "autodep", autoDep)
			logging.Infof("Successfully added module %s to Puppetfile\n", slug)
		}
	}
	return changes
}
