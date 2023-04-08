#
# MailHog Dockerfile
#

FROM golang:1.20-alpine as builder

# Install MailHog:
RUN apk --no-cache add --virtual build-dependencies git
RUN mkdir -p /root/gocode \
  && cd /root/gocode \
  && git clone https://github.com/jerry-yuan/Mail-Blackhole.git \
  && cd Mail-Blackhole \
  && go build -o MailBlackhole

FROM alpine:3
# Add mailhog user/group with uid/gid 1000.
# This is a workaround for boot2docker issue #581, see
# https://github.com/boot2docker/boot2docker/issues/581
RUN adduser -D -u 1000 mailblackhole

COPY --from=builder /root/gocode/Mail-Blackhole/MailBlackhole /usr/local/bin/MailBlackhole

USER mailblackhole

WORKDIR /home/mailblackhole

ENTRYPOINT ["MailBlackhole"]

# Expose the SMTP and HTTP ports:
EXPOSE 1025 8025
