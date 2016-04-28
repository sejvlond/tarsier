tarsier:
	docker run \
    	--rm \
	    -v `pwd`:/src \
	    sejvlond/tarsier_build \
	    github.com/sejvlond/tarsier

build: tarsier
	docker build -t sejvlond/tarsier .

push: build
	docker push sejvlond/tarsier

run:
	docker run \
		--rm \
		-it \
		-v `pwd`/secrets:/www/tarsier/secrets \
		-v `pwd`/logs:/www/tarsier/logs \
		-v `pwd`/tmp:/www/tarsier/tmp \
		-P \
		sejvlond/tarsier

clean:
	sudo rm tarsier
