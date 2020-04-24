FROM golang

ENV GOPATH=/go
ENV GOBIN=/go/covid/bin
ENV PATHVID=/go/covid

RUN echo $PATHVID

RUN mkdir -p $PATHVID/cmd/web
RUN mkdir -p $PATHVID/config
RUN mkdir -p $PATHVID/ui/html
run mkdir -p $PATHVID/bin

RUN mkdir -p /go/src/silverslanellc.com/covid/pkg/virusdata
RUN mkdir -p /go/src/github.com/Saied74/Lexer2

ADD cmd/web/*.* $PATHVID/cmd/web/
ADD ui/html/*.* $PATHVID/ui/html/
ADD config/*.* $PATHVID/config/
ADD config.csv $PATHVID

ADD pkg/virusdata/virusdata.go /go/src/silverslanellc.com/covid/pkg/virusdata
ADD pkg/lexer2/lexer.go /go/src/github.com/Saied74/Lexer2

RUN ls -l /go/covid/cmd/web
RUN ls -l $PATHVID

RUN go build -o=/go/covid/bin /go/covid/cmd/web/*.go

RUN ls -l /go/covid/bin

ENTRYPOINT ["/go/covid/bin/constants"]

EXPOSE 8080
