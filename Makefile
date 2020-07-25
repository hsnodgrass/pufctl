MAIN := main.go
BINARY := pufctl
EXE := pufctl.exe
ARCH := amd64
NIXBIN := bin/linux/${BINARY}
DARBIN := bin/darwin/${BINARY}
WINBIN := bin/windows/${EXE}
NIXBUILD := GOOS=linux GOARCH=${ARCH} go build -o ${NIXBIN} ${MAIN}
DARBUILD := GOOS=darwin GOARCH=${ARCH} go build -o ${DARBIN} ${MAIN}
WINBUILD := GOOS=windows GOARCH=${ARCH} go build -o ${WINBIN} ${MAIN}
NIXCLEAN := [ -f ${NIXBIN} ] && rm -f ${NIXBIN} || echo "No file to delete."
DARCLEAN := [ -f ${DARBIN} ] && rm -f ${DARBIN} || echo "No file to delete."
WINCLEAN := [ -f ${WINBIN} ] && rm -f ${WINBIN} || echo "No file to delete."

build: nixbuild darbuild winbuild

nixbuild:
	@echo "Performing linux quick build..." 
	@$(NIXBUILD)

darbuild:
	@echo "Performing darwin quick build..."
	@$(DARBUILD)

winbuild:
	@echo "Performing windows quick build..."
	@$(WINBUILD)

clean: nixclean darclean winclean

nixclean:
	@echo "Deleting linux binary..."
	@$(NIXCLEAN)

darclean:
	@echo "Deleting darwin binary..."
	@$(DARCLEAN)

winclean:
	@echo "Deleting windows exe..."
	@$(WINCLEAN)

run:
	@go run ${MAIN}
	

