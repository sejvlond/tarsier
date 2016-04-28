FROM scratch
MAINTAINER Ondřej Šejvl

WORKDIR /www/tarsier/

COPY conf/*.yaml conf/
COPY tarsier bin/

EXPOSE 8888 9999

VOLUME [ \
    "/www/tarsier/secrets", \
    "/www/tarsier/logs" \
    "/www/tarsier/tmp" \
]

USER 1000
ENTRYPOINT [ "/www/tarsier/bin/tarsier" ]
CMD [ "-c", "/www/tarsier/conf/tarsier.yaml" ]
