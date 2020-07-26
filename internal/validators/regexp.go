package validators

import "regexp"

var (
	reModSlug        = regexp.MustCompile(`^([\w]+)-([\w]+)$`)
	reModString      = regexp.MustCompile(`^([\w]+)\/([\w]+)$`)
	reGitURL         = regexp.MustCompile(`^(https?:\/\/[a-zA-Z0-9.]+\/|ssh:\/\/git@[a-zA-Z0-9.]+:|git@[a-zA-Z0-9.]+:)[a-z-A-Z0-9/]+.git$`)
	reGitURLNoSuffix = regexp.MustCompile(`^(https?:\/\/[a-zA-Z0-9.]+\/|ssh:\/\/git@[a-zA-Z0-9.]+:|git@[a-zA-Z0-9.]+:)[a-z-A-Z0-9/]+$`)
)
