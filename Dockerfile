#Variables that indicate where the application's main package is located
ENV APPLICATION_PACKAGE=./cmd/api
ENV CONF_DIR=/app/kit/config

ADD ./scripts/ /commands/
RUN chmod a+x /commands/*