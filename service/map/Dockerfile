FROM node:latest AS builder

WORKDIR /
RUN git clone https://github.com/chronark/atlas.git
WORKDIR /atlas

RUN yarn install
RUN yarn build


FROM nginx:1.12-alpine
COPY --from=builder /atlas/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
