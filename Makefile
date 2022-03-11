build:
	docker build -t build_static .

run:
	docker run -ti --rm -v $(CURDIR):/mnt build_static 
