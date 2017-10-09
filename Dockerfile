FROM golang:1.8
WORKDIR /app
ENV MYSQL_CONNECTION=user:password@tcp(sonar.paas.sbtech.com:3306)/build?charset=utf8
RUN  echo "{ \"type\": \"service_account\",\n \"project_id\": \"sbtech-pop-poc\",\n \"private_key_id\": \"4400cb812f33ef9f14648f48111cb461be71d5ad\",\n \"private_key\": \"-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCuJ5iGOWCZvMH+\n2LtMKIQ5zKLkwiT/m91pHCqfdZp7P1WpWau9pJUCaP6UNpmvWefcOCzCJi6u2j+b\nO9ya/f4hxlOueSoBygL6JXVhFaPIVR6NoDiClzxkUUkIISqlShc9yvuPG5tTK5cv\nQXO3S5Ux/Y9oh758xUnGwGIBd63tJy3LU3z3lO1K0nJ79xu0pM6Yynz+NGxiLeVL\n1q798nNl+wv7/4MxoMctJ0r8hFQK/BYR4JOFRTVh48Z0ySWivP9jDb3qzjKiMLl/\njLFBRn/7gaVDlDCVIIe7LbTRj1zoljd/d8MJcNE589/eA3Sq+f864nmFsV1NL4oU\nQphpp0l7AgMBAAECggEAEUWaCeMnjSJm9wMVqEWdv4G9EgTzLb/j4z3I4wxYfYAC\n2zXVI4zEE9DH7CPHdYLsIi9sDkvtIKhChf8XhO5LGvG+xjSmvDYFujtRezHDOZC/\njqTk6zFb6vMTftE97StTnMGZ1fqTRVoIYhu9VzYRc5DF0pnL2YyozCOzqe5C+0cC\nj6q8Or1/zvlMuFWKRvrfos9p23q2Pu+Aj680ix6AeS9eA6ck607tIZDs3mjk0y4a\naAY48LojxQ7JO0hV0YoEwNBVwKQeeu9L0/YuTaAmCEjL2hCPaa4rV1jpnZsYH53l\ne4KhpU+ptufugzVMFP+NOKVhHN7Uxey97jwopD9qIQKBgQDWT1PXz1kL1bZDLvg3\nmHxXpDQHKgoE9dOSV+yiPkvydgoyQ5PYrFEWX6ENKZQKIC5KgLG1+LXmnnTp27KS\nR4oLmFwBp9PlUtGdCUYLs8OsGKhK+HgHGuDMtWk2tZRIdFVk582IqA+mx4rbWI6d\nuqCyATdDLnPG0+pdEA46IYsqCwKBgQDQCIgVo8mPHuotAScsMXyOoEkGrxSU3PFg\n42hGZdw3gbQzSOkn2CXAuM4GtXlAyisXl2hbsXz5KyEZ4YLRuS6U/qeoB+3Z3Z5c\nwdG14S09qV1K5hEq33MduzmT0TJtxKP7K8EUS7fGqH0nbO0P/D6xo+kIU82StyBX\nh1UQZfZ0UQKBgAWK8bEwNK5graEZMuRFNloR7iKNTMsKDJnJWl5r3Az+To80Pjup\nYUOB36l2fNSyNmBI6c+6CgJX8NAnlXvBC+n229JTI+DXfoJgPgaJdLMbxCaEPJ56\nbHv+6SS4F4i0MY01jZROPHFk0cuQg8fCjutrqqE7L1ViH7qecq/ANFg3AoGBALyr\nkCf4FHwNsJpCWnGVK/76VWCPdt4Ph4a0l8SI6vEYXALLFFIDkTG5KXkiKqbc87oA\noi/Ox6X/PJUJVii4hwuv7QPStR+LA+3iBjyqzOoIaUjdYSJ95xxGydBKwS6WUZYh\nN4odpb2w31jbTCDcR0u6eUUJI/70wItBfJ9TqfIRAoGBAM9htQM4GCWsHQhEk3KW\nVYh6N25oC83lPhs+E/LogiFvgL9oAx+TnC76NQe8d/SjNz7v6gQy9m8RZRe0uTyC\n7J82NME/a00kR12Q2Wo76QvUBcLBxhb/I8EH8pHiYba5ZOtX8ICtHuy1Jfjlcq33\nco1SbdWYkS0ExUBmFBgmpZUf\n-----END PRIVATE KEY-----\",\n \"client_email\": \"popstatusservice@sbtech-pop-poc.iam.gserviceaccount.com\",\n \"client_id\": \"103842318454115281955\",\n \"auth_uri\": \"https://accounts.google.com/o/oauth2/auth\",\n \"token_uri\": \"https://accounts.google.com/o/oauth2/token\",\n \"auth_provider_x509_cert_url\": \"https://www.googleapis.com/oauth2/v1/certs\",\n \"client_x509_cert_url\": \"https://www.googleapis.com/robot/v1/metadata/x509/popstatusservice%40sbtech-pop-poc.iam.gserviceaccount.com\" }" >> ./gc_cloud.json
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/gc_cloud.json
EXPOSE 8080
COPY . /app
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/go-xorm/xorm
RUN go get github.com/gorilla/mux
RUN go get github.com/rs/cors
RUN go get golang.org/x/net/context
RUN go get golang.org/x/oauth2/google
RUN go get google.golang.org/api/compute/v1
RUN cd /app; go build -o myapp;
ENTRYPOINT ["./myapp"]