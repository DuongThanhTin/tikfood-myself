FROM node:22-alpine

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY package.json tsconfig.json ./
RUN npm install

COPY src ./src
COPY prompts ./prompts

RUN npm run build

EXPOSE 8080

CMD ["npm", "start"]
