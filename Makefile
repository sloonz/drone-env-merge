build:
	docker build -t sloonz/drone-env-merge .

publish:
	docker push sloonz/drone-env-merge
