FROM scratch

COPY guardian-server .
COPY config.json .

EXPOSE 8000

CMD ["./guardian-server"]
