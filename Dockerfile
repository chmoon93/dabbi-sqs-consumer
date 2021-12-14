FROM alpine:3.8

RUN apk add --no-cache libc6-compat

ADD ./dabbi-sqs-consumer /app/

CMD /app/dabbi-sqs-consumer
