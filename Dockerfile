FROM scratch

COPY guardian-server .
COPY config.json .

CMD ["./guardian-server"]
