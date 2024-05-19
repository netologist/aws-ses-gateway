# Use an official Node.js runtime as the base image
FROM node:22

# Set the working directory inside the container
WORKDIR /usr/src/app

# Install app dependencies
RUN npm install -g aws-ses-v2-local

# Expose the port your app runs on
EXPOSE 8005

# Define the command to run your app
CMD ["aws-ses-v2-local", "--host", "0.0.0.0", "--port", "8005"]
