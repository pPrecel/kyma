FROM eu.gcr.io/kyma-project/test-infra/buildpack-golang-toolbox:v20200117-d3885041 as builder

ENV BASE_APP_DIR /workspace/go/src/github.com/kyma-project/kyma/components/kyma-operator
WORKDIR ${BASE_APP_DIR}

#
# Copy files
#

COPY ./vendor/ ${BASE_APP_DIR}/vendor/
COPY ./pkg/ ${BASE_APP_DIR}/pkg/
COPY ./cmd/ ${BASE_APP_DIR}/cmd/
COPY ./licenses/ ${BASE_APP_DIR}/licenses/

#
# Build app
#

RUN CGO_ENABLED=0 GOOS=linux go build -v -o kyma-operator ${BASE_APP_DIR}/cmd/operator/main.go
RUN mkdir /app && mv ./kyma-operator /app/kyma-operator && mv ${BASE_APP_DIR}/licenses /app/licenses

FROM alpine:3.14.0
LABEL source = git@github.com:kyma-project/kyma.git
WORKDIR /app

#
# Install certificates
#

RUN apk update && apk add ca-certificates git jq && rm -rf /var/cache/apk/*

#
# Copy binary
#

COPY --from=builder --chown=nobody:nogroup /app /app

#
# Set permissions
# 

RUN mkdir /kyma &&\
	chown -R nobody:nogroup /app &&\
	chown -R nobody:nogroup /kyma
USER nobody

#
# Run app
#
ENTRYPOINT ["/app/kyma-operator"]
