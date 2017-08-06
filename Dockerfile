FROM iron/go:dev
WORKDIR /app
COPY . /app
RUN cd /app; go build -o myapp;
ENTRYPOINT ["./myapp"]