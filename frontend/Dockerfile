FROM node:latest

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

EXPOSE 3000

RUN npm run build

COPY /app/build /usr/share/nginx/html