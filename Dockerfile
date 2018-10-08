FROM scratch
EXPOSE 8080
ENTRYPOINT ["/ext-spotbugs"]
COPY ./bin/ /