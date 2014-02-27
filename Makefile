zabbix-aggregate-agent: main.go
	go build

test:
	go test zabbix_aggregate_agent/*_test.go

clean:
	rm -f zabbix-aggregate-agent index.html

binary:
	script/build.sh

index.html: README.md
	curl -s -H"Content-Type: text/x-markdown" -X POST --data-binary @README.md https://api.github.com/markdown/raw > index.html

release: index.html
	git checkout gh-pages
	git merge master
	script/build.sh
	git add bin && git commit -m "release binary" && git push origin gh-pages
	git checkout -
