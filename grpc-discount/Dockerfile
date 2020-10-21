FROM node:current-alpine

ADD index.js /app/
ADD package.json /app/
ADD package-lock.json /app/
ADD product.proto /app/

WORKDIR /app

RUN npm install

ENTRYPOINT [ "node", "index.js"]
