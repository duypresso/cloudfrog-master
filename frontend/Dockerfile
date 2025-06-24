FROM --platform=$BUILDPLATFORM node:18-alpine as build

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# Development server
FROM --platform=$TARGETPLATFORM node:18-alpine

WORKDIR /app
COPY --from=build /app /app

EXPOSE 3000
CMD ["npm", "run", "dev", "--", "--host"]
