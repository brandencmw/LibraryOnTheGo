FROM golang:1.19.2-bullseye

RUN mkdir /app

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /golibrary

EXPOSE 8080

CMD [ "/golibrary" ]