FROM iron/go:dev
WORKDIR /app
ENV MYSQL_CONNECTION
EXPOSE 8080
COPY . /app
RUN cd /app; go build -o myapp;
ENTRYPOINT ["./myapp"]