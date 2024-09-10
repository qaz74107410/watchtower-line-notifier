# Watchtower LINE Notifier

This project provides a webhook service that receives notifications from Watchtower and forwards them to LINE Notify. It's designed to be easy to set up and configure, with support for both standalone and Docker deployments.

## Features

- Receives webhook notifications from Watchtower
- Forwards notifications to LINE Notify
- Easy configuration using YAML
- Supports both standalone and Docker deployments
- Automatic config file creation from example

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
     -v /path/to/your/config.yaml:/app/config.yaml \
     yourusername/watchtower-line-notifier:latest
   ```

### Standalone Installation

1. Clone the repository:

   ```
   git clone https://github.com/yourusername/watchtower-line-notifier.git
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
   -e WATCHTOWER_NOTIFICATION_URL=generic+http://your-webhook-url:8080/webhook
   ```

   Replace `your-webhook-url` with the actual URL or IP address where your Watchtower LINE Notifier is running.

3. Start or restart your Watchtower container with these new settings.

Now, when Watchtower performs updates, it will send notifications to your webhook, which will forward them to your LINE account.

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