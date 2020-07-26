package auth

func TestGitAuth() {
	testCasesValid := []struct {
		url        string
		sshKeyPath string
		user       string
		pass       string
		token      string
	}{
		{
			url:  "https://gitserver.fake.com/fakeorg/fakemod",
			user: "fakedude",
			pass: "fakepass",
		},
		{
			url:   "https://gitserver.fake.com/fakeorg/fakemod",
			user:  "fakedude",
			token: "faketoken",
		},
		{
			url:  "http://gitserver.fake.com/fakeorg/fakemod",
			user: "fakedude",
			pass: "fakepass",
		},
		{
			url:   "http://gitserver.fake.com/fakeorg/fakemod",
			user:  "fakedude",
			token: "faketoken",
		},
		{
			url:        "ssh://git@gitserver.fake.com/fakeorg/fakemod.git",
			sshKeyPath: "/home/user/.ssh/id_rsa",
		},
		{
			url:        "git@gitserver.fake.com/fakeorg/fakemod.git",
			sshKeyPath: "/home/user/.ssh/id_rsa",
		},
		{
			url:        "ssh://git@gitserver.fake.com/fakeorg/fakemod",
			sshKeyPath: "/home/user/.ssh/id_rsa",
		},
		{
			url:        "git@gitserver.fake.com/fakeorg/fakemod",
			sshKeyPath: "/home/user/.ssh/id_rsa",
		},
	}
}
