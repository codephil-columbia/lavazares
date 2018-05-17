FROM golang:1.10

WORKDIR /

COPY lavazares ./

CMD ["./lavazares"]
