FROM scratch
EXPOSE 8800
COPY ca-certificates.crt /etc/ssl/certs/
COPY ./bin/server/lols /
CMD ["/lols", "-p", "8800"]
