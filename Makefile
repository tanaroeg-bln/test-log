build:
	docker buildx build --platform linux/amd64 -t test-log . 
	docker tag test-log:latest tanaroegbln/test-log:latest
	docker push tanaroegbln/test-log:latest