FROM fedora
MAINTAINER Kent Gustavsson <kent@minoris.se>

RUN dnf update -y
RUN mkdir -p /opt/lightweb

ADD config.toml /opt/lightweb/
ADD lightweb /opt/lightweb/

RUN adduser lightweb

USER lightweb

EXPOSE 8080
ENTRYPOINT ["/opt/lightweb/lightweb"]

