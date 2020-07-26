package responses

// ModuleJSONBody200 provides a string of a valid JSON response for
// a module fetch operation (status code 200)
const ModuleJSONBody200 = `{
	"uri": "/v3/modules/puppetlabs-apache",
	"slug": "puppetlabs-apache",
	"name": "apache",
	"downloads": 6618694,
	"created_at": "2010-05-19 05:46:26 -0700",
	"updated_at": "2010-05-19 05:46:26 -0700",
	"deprecated_at": null,
	"deprecated_for": null,
	"superseded_by": null,
	"supported": true,
	"endorsement": "supported",
	"module_group": "base",
	"owner": {
	  "uri": "/v3/users/puppetlabs",
	  "slug": "puppetlabs",
	  "username": "puppetlabs",
	  "gravatar_id": "fdd009b7c1ec96e088b389f773e87aec"
	},
	"current_release": {
	  "uri": "/v3/releases/puppetlabs-apache-4.0.0",
	  "slug": "puppetlabs-apache-4.0.0",
	  "module": {
		"uri": "/v3/modules/puppetlabs-apache",
		"slug": "puppetlabs-apache",
		"name": "apache",
		"deprecated_at": null,
		"owner": {
		  "uri": "/v3/users/puppetlabs",
		  "slug": "puppetlabs",
		  "username": "puppetlabs",
		  "gravatar_id": "fdd009b7c1ec96e088b389f773e87aec"
		}
	  },
	  "version": "4.0.0",
	  "metadata": {},
	  "tags": [
		"string"
	  ],
	  "supported": true,
	  "pdk": true,
	  "validation_score": 100,
	  "file_uri": "/v3/files/puppetlabs-apache-4.0.0.tar.gz",
	  "file_size": 317119,
	  "file_md5": "e579a02bd16f9dec56018d407be7ab32",
	  "file_sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	  "downloads": 12782,
	  "readme": "string",
	  "changelog": "string",
	  "license": "string",
	  "reference": "string",
	  "docs": {},
	  "pe_compatibility": [
		"string"
	  ],
	  "tasks": [
		{
		  "name": "init",
		  "executable": "init.sh",
		  "executables": [
			"init.sh",
			"init.ps1"
		  ],
		  "description": "Allows you to perform apache service functions",
		  "metadata": {}
		}
	  ],
	  "plans": [
		{
		  "uri": "/v3/releases/puppetlabs-lvm-1.4.0/plans/lvm::expand",
		  "name": "lvm::expand"
		}
	  ],
	  "created_at": "2019-01-10 05:18:00 -0800",
	  "updated_at": "2019-01-10 05:20:17 -0800",
	  "deleted_at": null,
	  "deleted_for": null
	},
	"releases": [
	  {
		"uri": "/v3/releases/puppetlabs-apache-4.0.0",
		"slug": "puppetlabs-apache-4.0.0",
		"version": "4.0.0",
		"supported": true,
		"created_at": "2019-01-10 05:18:00 -0800",
		"deleted_at": null,
		"file_uri": "/v3/files/puppetlabs-apache-4.0.0.tar.gz",
		"file_size": 317119
	  }
	],
	"feedback_score": 84,
	"homepage_url": "https://github.com/puppetlabs/puppetlabs-apache",
	"issues_url": "https://tickets.puppetlabs.com/browse/MODULES"
  }`
