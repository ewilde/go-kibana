FROM node:9.7.1-alpine
RUN npm install -g makelogs
RUN apk update && apk add curl
ADD entrypoint.sh /entrypoint.sh
ADD scripts /scripts

ENTRYPOINT ["/entrypoint.sh"]
