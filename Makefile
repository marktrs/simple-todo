deps-cleancache:
	go clean -modcache

test-unit:
	go test ./...

test-coverage:
	go test $(go list ./...) -race -covermode atomic -coverprofile=coverage.out ./...

test-report:
	go tool cover -html=coverage.out

mock:
	mockgen -source=repository/user.go \
		-package testutil \
		-destination=testutil/mocks/repository/user.go
	mockgen -source=repository/task.go \
		-package testutil \
		-destination=testutil/mocks/repository/task.go