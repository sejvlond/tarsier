FROM scratch
MAINTAINER Ondřej Šejvl

VOLUME [ \
    "/www/tarsier/conf", \
    "/www/tarsier/secrets", \
    "/www/tarsier/logs" \
    "/www/tarsier/tmp" \
]

WORKDIR /www/tarsier/

# default configuration
COPY conf/tarsier.yaml conf/
COPY conf/kafkafeeder.yaml logs/
COPY tarsier bin/

EXPOSE 8888 9999

ENTRYPOINT [ "/www/tarsier/bin/tarsier" ]
CMD [ "-c", "/www/tarsier/conf/tarsier.yaml" ]
