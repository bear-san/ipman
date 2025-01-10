FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY ./bin/ipman /ipman
USER 65532:65532

ENTRYPOINT ["/ipman"]
