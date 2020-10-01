# Mock interfaces with mockery
# inpgk and testonly adds mock to the test files, so it could be ignored by coverage
# Other to ignore mocks is to use // +build !test and the run
# "go test -v -cover -tags test" but it seems it will have to be used in all packages
$GOPATH/bin/mockery -all -testonly -inpkg

# generate cover html
go tool cover -html=cover.out


# test all without mocks
# shows low coverage in case of -race
go test -race -coverpkg $(go list ./... | grep -v mocks | tr '\n' ',') ./... -coverprofile cover.out
go test -coverpkg $(go list ./... | grep -v mocks | tr '\n' ',') ./... -coverprofile cover.out
go test -coverprofile cover.out
