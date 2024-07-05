FROM node:19.5.0-alpine

COPY package.json /app/
COPY src /app/

WORKDIR /app/

RUN npm install

CMD ["node", "server.js"]