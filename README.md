# Watchtower LINE Notifier

Watchtower LINE Notifier is a service that receives notifications from [Watchtower](https://github.com/containrrr/watchtower) and forwards them to LINE Notify. It's designed to be easy to set up and configure, with support for both standalone and Docker deployments.

## Table of Contents

1. [Features](#features)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
   - [Docker Installation](#docker-installation)
   - [Standalone Installation](#standalone-installation)
4. [Configuration](#configuration)
5. [Usage](#usage)
6. [Docker Compose Setup](#docker-compose-setup)
7. [Automated Builds](#automated-builds)
8. [Development](#development)
9. [License](#license)
10. [Acknowledgments](#acknowledgments)

## Features

- Receives webhook notifications from Watchtower
- Forwards notifications to LINE Notify
- Easy configuration using YAML
- Supports both standalone and Docker deployments
- Automatic config file creation from example
- Accepts raw body content in webhook requests

## Prerequisites

- Go 1.16+ (for standalone installation)
- Docker (for Docker installation)
- LINE Notify account and token

## Installation

### Docker Installation

1. Pull the Docker image from Docker Hub:

   ```
   docker pull qaz74107410/watchtower-line-notifier:latest
   ```

2. Create a `config.yaml` file on your host machine with your LINE Notify token:

   ```yaml
   server:
     port: 8080

   line:
     api_endpoint: "https://notify-api.line.me/api/notify"
     token: "YOUR_LINE_NOTIFY_TOKEN"
   ```

3. Run the Docker container, mounting your config file:

   ```
   docker run -d \
     --name watchtower-line-notifier \
     -p 8080:8080 \
     -v /path/to/your/config.yaml:/root/config.yaml \
     qaz74107410/watchtower-line-notifier:latest
   ```

### Standalone Installation

1. Clone the repository:

   ```
   git clone https://github.com/qaz74107410/watchtower-line-notifier.git
   cd watchtower-line-notifier
   ```

2. Install dependencies:

   ```
   go mod tidy
   ```

3. Build the application:

   ```
   go build -o watchtower-line-notifier
   ```

4. Run the application:

   ```
   ./watchtower-line-notifier
   ```

   On first run, it will create a `config.yaml` file from the example if it doesn't exist. Edit this file to include your LINE Notify token.

## Configuration

The `config.yaml` file contains the following settings:

```yaml
server:
  port: 8080  # The port on which the webhook server will listen

line:
  api_endpoint: "https://notify-api.line.me/api/notify"  # LINE Notify API endpoint
  token: "YOUR_LINE_NOTIFY_TOKEN"  # Your LINE Notify token
```

## Usage

1. Ensure the Watchtower LINE Notifier is running and accessible.

2. Configure Watchtower to send notifications to your webhook URL. Add the following environment variables to your Watchtower container:

   ```
   -e WATCHTOWER_NOTIFICATIONS=http
   -e WATCHTOWER_NOTIFICATION_URL="http://your-webhook-url:8080/webhook"
   ```

   Replace `your-webhook-url` with the actual URL or IP address where your Watchtower LINE Notifier is running.

3. Start or restart your Watchtower container with these new settings.

Now, when Watchtower performs updates, it will send notifications to your webhook, which will forward them to your LINE account.

## Docker Compose Setup

For easy setup of both the Watchtower LINE Notifier and Watchtower itself, you can use Docker Compose:

1. Create a `docker-compose.yml` file with the following content:

   ```yaml
   version: '3'

   services:
     watchtower-line-notifier:
       image: qaz74107410/watchtower-line-notifier:latest
       container_name: watchtower-line-notifier
       ports:
         - "8080:8080"
       volumes:
         - ./config.yaml:/root/config.yaml
       restart: unless-stopped

     watchtower:
       image: containrrr/watchtower
       container_name: watchtower
       volumes:
         - /var/run/docker.sock:/var/run/docker.sock
       environment:
         - WATCHTOWER_NOTIFICATIONS=http
         - WATCHTOWER_NOTIFICATION_URL=http://watchtower-line-notifier:8080/webhook
       restart: unless-stopped

   networks:
     default:
       name: watchtower-network
   ```

2. Ensure you have a `config.yaml` file in the same directory with your LINE Notify token.

3. Run the services:

   ```bash
   docker-compose up -d
   ```

This will start both the Watchtower LINE Notifier and Watchtower services.

## Automated Builds

This project uses GitHub Actions to automatically build and push the Docker image to Docker Hub whenever changes are pushed to the `main` branch.

The workflow does the following:
1. Builds the Docker image
2. Logs in to Docker Hub using secure credentials
3. Pushes the new image to Docker Hub with the `latest` tag

You can find the workflow configuration in `.github/workflows/docker-build-push.yml`.

### For Contributors

If you're contributing to this project, you don't need to manually build and push Docker images. Just push your changes to the `main` branch (or merge a pull request), and GitHub Actions will handle the build and deploy process automatically.

### For Users

Thanks to this automated process, you can always be sure that the `latest` tag of our Docker image contains the most up-to-date version of the application. You don't need to do anything special to get the latest version - just pull the `latest` tag as usual:

```bash
docker pull qaz74107410/watchtower-line-notifier:latest
```

## Development

To contribute to this project:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Watchtower](https://github.com/containrrr/watchtower) for container updating
- [LINE Notify](https://notify-bot.line.me/) for notification services
- [Viper](https://github.com/spf13/viper) for configuration management