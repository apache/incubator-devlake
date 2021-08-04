# NOTE: If you want this file to be re-compiled with docker-compose, you need to run `docker-compose up --build`

# Install node v14
FROM node:14

# Set the workdir /var/www/lake
WORKDIR /var/www/lake

# Copy the package.json to workdir
COPY package.json ./

# Copy application source
COPY . .  

# Run npm install - install the npm dependencies
RUN npm install

# Expose application port
EXPOSE 3000
EXPOSE 3001

ENV NODE_ENV=docker

# Start the application
CMD ["npm", "run", "all-docker"]