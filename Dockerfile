FROM scratch
MAINTAINER Stevenson Jean-Pierre <stevenson.jean-pierre@sage.com>
ADD bin/linux_service_now_connector service_now_connector
EXPOSE 8080
ENTRYPOINT ["/service_now_connector"]