test:
	go clean -testcache && go test -cover ./...

build: test
	go build

install: test
	go install

release:
	docker build -t checkip-releases -f Releases.Dockerfile .
	docker create -ti --name checkip-releases checkip-releases sh
	test -d releases || mkdir releases
	docker cp checkip-releases:/releases/checkip_linux_amd64.tar.gz releases/
	docker cp checkip-releases:/releases/checkip_darwin_amd64.tar.gz releases/
	docker rm checkip-releases
	docker rmi checkip-releases