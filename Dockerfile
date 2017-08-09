FROM iron/go:dev
WORKDIR /app
ENV MYSQL_CONNECTION=user:password@tcp(sonar.paas.sbtech.com:3306)/build?charset=utf8
EXPOSE 8080
COPY . /app
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/go-xorm/xorm
RUN go get github.com/gorilla/mux
RUN go get github.com/rs/cors
RUN cd /app; go build -o myapp;
ENTRYPOINT ["./myapp"]