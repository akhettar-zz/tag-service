FROM projectwave/golang-alpine

ENV ENVIRONMENT default
ENV SRC_FOLDER /go/src/github.com/tag-service
ENV CONFIG_FOLDER /go/config
ENV PKG_FOLDER /go/pkg

RUN mkdir -p $SRC_FOLDER $PKG_FOLDER

WORKDIR $SRC_FOLDER

COPY . $SRC_FOLDER
COPY config $CONFIG_FOLDER

RUN go install && rm -rf $PKG_FOLDER

HEALTHCHECK --interval=5s --retries=10 CMD curl -fs http://localhost:8080/health || exit 1

EXPOSE 8080

CMD /go/bin/tag-service
