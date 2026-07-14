# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# debug image is used to get busybox
FROM gcr.io/distroless/static-debian13:debug@sha256:e741251ccc55dd6cec4a99ff21c0766df31891fabb4a50727104619a7e6ff4f2
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/octl"]
COPY $TARGETPLATFORM/octl /usr/bin/
