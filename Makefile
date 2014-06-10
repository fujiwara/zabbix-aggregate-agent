GIT_REV := $(shell git rev-parse --short HEAD)

all:
	sed -i.bak -e "s/HEAD/$(GIT_REV)/" revision.go
	gox -os="linux darwin" -arch="amd64 386" -output "pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"
	git checkout revision.go
	rm revision.go.bak

test:
	go test zabbix_aggregate_agent/config_test.go

pages:
	cat html/index_head.html > index.html
	curl -s -H"Content-Type: text/x-markdown" -X POST --data-binary @README.md https://api.github.com/markdown/raw >> index.html
	cat html/index_foot.html >> index.html

release_doc:
	git checkout gh-pages
	git merge master
	make pages
	git add index.html && git commit -m "release pages" && git push origin gh-pages
	git checkout -
