// Package forgeapi provides a (incomplete) API for the Puppet Forge
package forgeapi

// Module is a struct representation of the Puppet Forge OpenAPI Module object
type Module struct {
	URI            string               `json:"uri" alias:"uri"`
	Slug           string               `json:"slug" alias:"slug"`
	Name           string               `json:"name" alias:"name"`
	Downloads      int                  `json:"downloads" alias:"downloads"`
	CreatedAt      string               `json:"created_at" alias:"createdat,created_at"`
	UpdatedAt      string               `json:"updated_at" alias:"updatedat,updated_at"`
	DeprecatedAt   string               `json:"deprecated_at,omitempty" alias:"deprecatedat,deprecated_at"`
	DeprecatedFor  string               `json:"deprecated_for,omitempty" alias:"deprecatedfor,deprecated_for"`
	SupersededBy   ModuleMinimal        `json:"superseded_by,omitempty" alias:"supersededby,superseded_by"`
	Supported      bool                 `json:"supported" alias:"supported"`
	Endorsement    string               `json:"endorsement,omitempty" alias:"endorsement"`
	ModuleGroup    string               `json:"module_group" alias:"modulegroup,module_group"`
	Owner          UserAbbreviated      `json:"owner" alias:"owner"`
	CurrentRelease Release              `json:"current_release" alias:"currentrelease,current_release"`
	Releases       []ReleaseAbbreviated `json:"releases,omitempty" alias:"releases"`
	FeedbackScore  int                  `json:"feedback_score" alias:"feedbackscore,feedback_score"`
	HomepageURL    string               `json:"homepage_url" alias:"homepageurl,homepage_url"`
	IssuesURL      string               `json:"issues_url" alias:"issuesurl,issues_url"`
}

// ModuleAbbreviated is a struct representation of the Puppet Forge OpenAPI ModuleAbbreviated object
type ModuleAbbreviated struct {
	URI          string          `json:"uri" alias:"uri"`
	Slug         string          `json:"slug" alias:"slug"`
	Name         string          `json:"name" alias:"name"`
	DeprecatedAt string          `json:"deprecated_at,omitempty" alias:"deprecatedat,deprecated_at"`
	Owner        UserAbbreviated `json:"owner" alias:"owner"`
}

// ModuleMinimal is a struct representation of the Puppet Forge OpenAPI ModuleMinimal object
type ModuleMinimal struct {
	URI  string `json:"uri" alias:"uri"`
	Slug string `json:"slug" alias:"slug"`
}

// Release is a struct representation of the Puppet Forge OpenAPI Release object
type Release struct {
	URI             string                   `json:"uri" alias:"uri"`
	Slug            string                   `json:"slug" alias:"slug"`
	Module          ModuleAbbreviated        `json:"module" alias:"module"`
	Version         string                   `json:"version" alias:"version"`
	Metadata        ModuleMetadata           `json:"metadata" alias:"metadata"`
	Tags            []string                 `json:"tags" alias:"tags"`
	Supported       bool                     `json:"supported" alias:"supported"`
	PDK             bool                     `json:"pdk" alias:"pdk"`
	ValidationScore int                      `json:"validation_score" alias:"validationscore,validation_score"`
	FileURI         string                   `json:"file_uri" alias:"fileuri,file_uri"`
	FileSize        int                      `json:"file_size" alias:"filesize,file_size"`
	FileMD5         string                   `json:"file_md5" alias:"filemd5,file_md5"`
	FileSHA256      string                   `json:"file_sha256" alias:"filesha256,file_sha256"`
	Downloads       int                      `json:"downloads" alias:"downloads"`
	ReadMe          string                   `json:"readme,omitempty" alias:"readme"`
	ChangeLog       string                   `json:"changelog,omitempty" alias:"changelog"`
	License         string                   `json:"license,omitempty" alias:"license"`
	Reference       string                   `json:"reference,omitempty" alias:"reference"`
	Tasks           []ReleaseTask            `json:"tasks,omitempty" alias:"tasks"`
	Plans           []ReleasePlanAbbreviated `json:"plans,omitempty" alias:"plans"`
	CreatedAt       string                   `json:"created_at" alias:"createdat,created_at"`
	UpdatedAt       string                   `json:"updated_at" alias:"updatedat,updated_at"`
	DeletedAt       string                   `json:"deleted_at,omitempty" alias:"deletedat,deleted_at"`
	DeletedFor      string                   `json:"deleted_for,omitempty" alias:"deletedfor,deleted_for"`
}

// ReleaseAbbreviated is a struct representation of the Puppet Forge OpenAPI ReleaseAbbreviated object
type ReleaseAbbreviated struct {
	URI       string            `json:"uri" alias:"uri"`
	Slug      string            `json:"slug" alias:"slug"`
	Module    ModuleAbbreviated `json:"module" alias:"module"`
	Version   string            `json:"version" alias:"version"`
	Supported bool              `json:"supported" alias:"supported"`
	CreatedAt string            `json:"created_at" alias:"createdat,created_at"`
	DeletedAt string            `json:"deleted_at,omitempty" alias:"deletedat,deleted_at"`
	FileURI   string            `json:"file_uri" alias:"fileuri,file_uri"`
	FileSize  int               `json:"file_size" alias:"filesize,file_size"`
}

// ReleaseMinimal is a struct representation of the Puppet Forge OpenAPI ReleaseMinimal object
type ReleaseMinimal struct {
	URI     string `json:"uri" alias:"uri"`
	Slug    string `json:"slug" alias:"slug"`
	FileURI string `json:"file_uri" alias:"fileuri,file_uri"`
}

// ReleaseTask is a struct representation of the Puppet Forge OpenAPI ReleaseTask object
type ReleaseTask struct {
	Name        string       `json:"name" alias:"name"`
	Executables []string     `json:"executables" alias:"executables"`
	Description string       `json:"description" alias:"description"`
	Metadata    TaskMetadata `json:"metadata" alias:"metadata"`
}

// ReleasePlanAbbreviated is a struct representation of the Puppet Forge OpenAPI ReleasePlanAbbreviated object
type ReleasePlanAbbreviated struct {
	URI  string `json:"uri" alias:"uri"`
	Name string `json:"name" alias:"name"`
}

// UserAbbreviated is a struct representation of the Puppet Forge OpenAPI UserAbbreviated object
type UserAbbreviated struct {
	URI        string `json:"uri" alias:"uri"`
	Slug       string `json:"slug" alias:"slug"`
	Username   string `json:"username" alias:"username"`
	GravatarID string `json:"gravatar_id,omitempty" alias:"gravatarid,gravatar_id"`
}

// ModuleMetadata is a struct representation of a metadata.json file
type ModuleMetadata struct {
	Name                   string                     `json:"name" alias:"name"`
	Version                string                     `json:"version" alias:"version"`
	Author                 string                     `json:"author" alias:"author"`
	License                string                     `json:"license" alias:"license"`
	Summary                string                     `json:"summary" alias:"summary"`
	Source                 string                     `json:"source" alias:"source"`
	Dependencies           []ModuleMetadataDependency `json:"dependencies,omitempty" alias:"dependencies"`
	Requirements           []ModuleMetadataDependency `json:"requirements,omitempty" alias:"requirements"`
	ProjectPage            string                     `json:"project_page,omitempty" alias:"projectpage,project_page"`
	IssuesURL              string                     `json:"issues_url,omitempty" alias:"issues_url,issuesurl"`
	OperatingSystemSupport []ModuleMetadataOS         `json:"operatingsystem_support,omitempty" alias:"operatingsystemsupport,operatingsystem_support,os"`
	Tags                   []string                   `json:"tags,omitempty" alias:"tags"`
}

// ModuleMetadataDependency is a struct representation of a dependency / requirement
// in a metadata.json file
type ModuleMetadataDependency struct {
	Name               string `json:"name" alias:"name"`
	VersionRequirement string `json:"version_requirement" alias:"versionrequirement,version_requirement,version"`
}

// ModuleMetadataOS is a struct representation of a supported OS in a metadata.json file
type ModuleMetadataOS struct {
	OperatingSystem        string   `json:"operatingsystem" alias:"operatingsystem,os"`
	OperatingSystemRelease []string `json:"operatingsystemrelease" alias:"operatingsystemrelease,release"`
}

// TaskMetadata is a struct representation of a Bolt task's metadata file
type TaskMetadata struct {
	Description       string                       `json:"description,omitempty" alias:"description"`
	InputMethod       string                       `json:"input_method" alias:"inputmethod,input_method"`
	PuppetTaskVersion int                          `json:"puppet_task_version" alias:"puppet_task_version,puppettaskversion,taskversion,task_version"`
	SupportsNoop      bool                         `json:"supports_noop" alias:"supportsnoop,supports_noop"`
	Implementations   []TaskMetadataImplementation `json:"implementations,omitempty" alias:"implementations"`
	Files             []string                     `json:"files,omitempty" alias:"files"`
	Private           bool                         `json:"private" alias:"private"`
	Remote            bool                         `json:"remote" alias:"remote"`
}

// TaskMetadataImplementation is a struct representation of implementations in a
// Bolt task's metadata file
type TaskMetadataImplementation struct {
	Name         string   `json:"name" alias:"name"`
	Requirements []string `json:"requirements" alias:"requirements"`
}

// ListModulesOpts is used by the List* functions to construct the
// search parameters for the list modules API endpoint
type ListModulesOpts struct {
	Limit           int      `url:"limit,omitempty"`
	Offset          int      `url:"offset,omitempty"`
	SortBy          string   `url:"sort_by,omitempty"`
	Query           string   `url:"query,omitempty"`
	Tag             string   `url:"tag,omitempty"`
	Owner           string   `url:"owner,omitempty"`
	WithTasks       bool     `url:"with_tasks,omitempty"`
	WithPlans       bool     `url:"with_plans,omitempty"`
	WithPDK         bool     `url:"with_pdk,omitempty"`
	Endorsements    []string `url:"endorsements,omitempty"`
	OperatingSystem string   `url:"operatingsystem,omitempty"`
	PERequirement   string   `url:"pe_requirement,omitempty"`
	PupRequirement  string   `url:"puppet_requirement,omitempty"`
	MinScore        int      `url:"with_minimum_score,omitempty"`
	ModuleGroups    []string `url:"module_groups,omitempty"`
	ShowDeleted     bool     `url:"show_deleted,omitempty"`
	HideDeprecated  bool     `url:"hide_deprecated,omitempty"`
	OnlyLatest      bool     `url:"only_latest,omitempty"`
	Slugs           []string `url:"slugs,omitempty"`
	WithHTML        bool     `url:"with_html,omitempty"`
	IncludeFields   []string `url:"include_fields,omitempty"`
	ExcludeFields   []string `url:"exclude_fields,omitempty"`
}

// ListResults is a struct representation of the JSON return
// of the ForgeAPI listmodules endpoint
type ListResults struct {
	Pagination ListPagination `json:"pagination"`
	Results    []Module       `json:"results"`
}

// ListPagination is a struct representing the pagination field
// of the JSON return from the ForgeAPI listmodules endpoint
type ListPagination struct {
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	First    string `json:"first"`
	Previous string `json:"previous"`
	Current  string `json:"current"`
	Next     string `json:"next"`
	Total    int    `json:"total"`
}
