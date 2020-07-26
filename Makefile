MAIN := main.go
PUFCTLVER := 0.0.0
DATE := $(shell date --iso)
BINARY := pufctl
LINUX := linux
DARWIN := darwin
WINDOWS := windows
ARCH := amd64
GOARCH := GOARCH=${ARCH}
GOOS := GOOS=
BUILDCMD := go build
PUFCTLPKG := github.com/hsnodgrass/pufctl/internal

define quickbuildcmd
	echo "Performing $(1) quick build..."
	${GOOS}$(1) ${GOARCH} ${BUILDCMD} -o bin/$(1)/$(2) ${MAIN}
endef

define cleancmd
	echo "Deleting $(1) binaries..."
	rm -rf bin/$(1)/*
endef

define releasecmd
	echo "Compiling pufctl version ${PUFCTLVER} release binary for $(1)..."
	mkdir -p release/$(1)/
	${GOOS}$(1) ${GOARCH} ${BUILDCMD} \
	-ldflags="-X '${PUFCTLPKG}/version.Ver=${PUFCTLVER}' -X '${PUFCTLPKG}/version.OS=$(1)' -X '${PUFCTLPKG}/version.Arch=${ARCH}' -X '${PUFCTLPKG}/version.Time=$(3)'" \
	-o release/$(1)/$(2) ${MAIN}
endef

define checksumcmd
	echo "Creating sha256 checksum for release binary $(2)..."
	$(shell sha256sum release/$(1)/$(2) > release/$(1)/sha256checksum.txt)
endef


build: nixbuild darbuild winbuild

nixbuild:
	@$(call quickbuildcmd,${LINUX},${BINARY})

darbuild:
	@$(call quickbuildcmd,${DARWIN},${BINARY})

winbuild:
	@$(call quickbuildcmd,${WINDOWS},${BINARY}.exe)

clean: nixclean darclean winclean

nixclean:
	@$(call cleancmd,${LINUX})

darclean:
	@$(call cleancmd,${DARWIN})

winclean:
	@$(call cleancmd,${WINDOWS})

run:
	@go run ${MAIN}

checksum: nixchecksum darchecksum winchecksum

nixchecksum:
	@$(call checksumcmd,${LINUX},${BINARY})

darchecksum:
	@$(call checksumcmd,${DARWIN},${BINARY})

winchecksum:
	@$(call checksumcmd,${WINDOWS},${BINARY}.exe)

sign:
	gon -log-level=info ./gon.json

release: nixrelease darrelease winrelease checksum

nixrelease:
	@$(call releasecmd,${LINUX},${BINARY},$(DATE))

darrelease:
	@$(call releasecmd,${DARWIN},${BINARY},$(DATE))

winrelease:
	@$(call releasecmd,${WINDOWS},${BINARY}.exe,$(DATE))

