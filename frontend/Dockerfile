FROM node:24-slim AS base
WORKDIR /app
COPY package*.json ./

FROM base AS development
RUN npm install
COPY . .
CMD ["npm", "run", "dev"]
