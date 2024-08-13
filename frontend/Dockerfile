FROM node:alpine as build

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

EXPOSE 3000

RUN npm run build

FROM nginx:alpine

COPY --from=build /app/build /usr/share/nginx/html

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]