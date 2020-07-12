FROM node as prod

WORKDIR /app
COPY package*.json ./

RUN npm install
WORKDIR /app/client

COPY ./client/package*.json ./

RUN npm install

WORKDIR /app

COPY . .
ENV NODE_ENV=production
CMD [ "npm", "start" ]