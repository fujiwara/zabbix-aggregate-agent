zabbix-aggregate-agent:
	go build

test:
	go test zabbix_aggregate_agent/*_test.go

clean:
	rm -f zabbix-aggregate-agent

binary:
	script/build.sh

pages:
	cat html/index_head.html > index.html
	curl -s -H"Content-Type: text/x-markdown" -X POST --data-binary @README.md https://api.github.com/markdown/raw >> index.html
	cat html/index_foot.html >> index.html

release:
	git checkout gh-pages
	git merge master
	make pages
	script/build.sh
	git add bin index.html && git commit -m "release binary" && git push origin gh-pages
	git checkout -
