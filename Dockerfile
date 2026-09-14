# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# debug image is used to get busybox
FROM gcr.io/distroless/static-debian13:debug@sha256:07148a6899406df51906b183f581cf66e5b05fd51aca438bdf1d3998df566961
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/octl"]
COPY $TARGETPLATFORM/octl /usr/bin/
