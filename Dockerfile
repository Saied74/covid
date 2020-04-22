FROM golang

ENV GOPATH=/go
ENV GOBIN=/go/covid/bin
ENV CONFIG=/go/covid

RUN echo $GOBIN

RUN mkdir -p /go/covid/cmd/web
RUN mkdir -p /go/src/covid/config
RUN mkdir -p /go/src/silverslanellc.com/covid/pkg/virusdata
RUN mkdir -p /go/src/github.com/Saied74/Lexer2
RUN mkdir -p /go/src/ui/html
run mkdir -p /go/covid/bin

ADD cmd/web/*.* /go/covid/cmd/web/
ADD ui/html/*.* /go/src/covid/ui/html/
ADD config/*.* /go/src/covid/config/
ADD pkg/virusdata/virusdata.go /go/src/silverslanellc.com/covid/pkg/virusdata
ADD pkg/lexer2/lexer.go /go/src/github.com/Saied74/Lexer2
ADD config.csv /go/covid

RUN ls -l /go/covid/cmd/web



RUN go build -o=/go/covid/bin /go/covid/cmd/web/*.go

RUN ls -l /go/covid/bin

RUN ls -l /go/covid/cmd/web

ENTRYPOINT ["/go/covid/bin/constants", "-e=CONFIG"]

EXPOSE 8080
