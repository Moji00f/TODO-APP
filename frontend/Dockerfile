FROM node:22.14.0-alpine AS build-stage
WORKDIR /app
COPY package*.json ./
RUN npm config set registry https://registry.npmmirror.com/ && \
    npm install --include=dev
COPY . .
ENV VITE_API_URL=http://localhost:5000/api
RUN npm run build

FROM nginx:1.25.1-alpine
COPY --from=build-stage /app/dist /usr/share/nginx/html
EXPOSE 80 
CMD ["nginx", "-g", "daemon off;"]

