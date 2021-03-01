FROM node:12.18-alpine
COPY package*.json ./
RUN npm ci
COPY . .

EXPOSE 8080
CMD ["node", "index.js"]