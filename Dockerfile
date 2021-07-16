# Install node v14
FROM node:14

# Set the workdir /var/www/lake
WORKDIR /var/www/lake

# Copy the package.json to workdir
COPY package.json ./

# Run npm install - install the npm dependencies
RUN npm install

# Copy application source
COPY . .  

# Copy .env.docker to workdir/.env - use the docker env
COPY .env.docker ./.env

# Expose application port
EXPOSE 3000

# Generate build
# RUN npm run build

# Start the application
# CMD node index.js