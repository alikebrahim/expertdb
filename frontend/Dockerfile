FROM node:18-alpine

WORKDIR /app

# Copy package files and install dependencies
COPY package.json ./
RUN npm install

# Copy remaining files - these will be mounted in dev mode
COPY . .

# Set environment to development
ENV NODE_ENV=development

# Expose the listening port
EXPOSE 3000

# Run the application in development mode with host binding
CMD ["npm", "run", "dev"]
