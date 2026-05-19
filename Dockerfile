# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# debug image is used to get busybox
FROM gcr.io/distroless/static-debian13:debug@sha256:85554049bdfa232b64a328412dd8909cc3baad08474cf97b2e5cc0a74e23fc5e
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/octl"]
COPY $TARGETPLATFORM/octl /usr/bin/
