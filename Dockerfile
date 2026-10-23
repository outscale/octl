# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# debug image is used to get busybox
FROM gcr.io/distroless/static-debian13:debug@sha256:6a5e9a4921c51d613ab8128402badf19d27b16627f92f97146f4597d2c593104
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/octl"]
COPY $TARGETPLATFORM/octl /usr/bin/
