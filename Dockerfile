FROM golang:alpine3.11 AS build

RUN apk add --no-cache git make

COPY . /go/src/project/
WORKDIR /go/src/project/
RUN make

FROM hashicorp/terraform:0.13.3
COPY --from=build /go/src/project/build/* /bin/
ENTRYPOINT ["/bin/sh"]