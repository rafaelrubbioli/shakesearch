run:
	go run cmd/api.go

deploy:
	docker build -t docker.pkg.github.com/rafaelrubbioli/shakesearch/shakesearch:01 .
	docker push docker.pkg.github.com/rafaelrubbioli/shakesearch/shakesearch:01
