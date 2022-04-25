FROM agent/agent as build-agent
FROM alpine

# no-cache flag: https://github.com/gliderlabs/docker-alpine/blob/master/docs/usage.md#disabling-cache

RUN apk add --no-cache python3
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools

# Add non root user
RUN addgroup -S app && adduser app -S -G app

WORKDIR /home/app/

COPY .. .

RUN chown -R app /home/app

USER app

COPY --from=build-agent /build/agent .

ENV fprocess="python3 index.py"

EXPOSE 80

CMD [ "./agent" ]