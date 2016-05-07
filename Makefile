IMAGE=sejvlond/tarsier

tarsier:
	docker run \
    	--rm \
	    -v `pwd`:/src \
	    sejvlond/tarsier_build \
	    github.com/sejvlond/tarsier

build: tarsier
	docker build -t ${IMAGE} .

push: build
	docker push ${IMAGE}

run:
	docker run \
		--rm \
		-it \
		-v `pwd`/secrets:/www/tarsier/secrets \
		-v `pwd`/logs:/www/tarsier/logs \
		-v `pwd`/tmp:/www/tarsier/tmp \
		-P \
		${IMAGE}

clean:
	sudo rm tarsier
