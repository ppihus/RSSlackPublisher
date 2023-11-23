# RSSlackPublisher

RSSlackPublisher is an automated tool written in Go, designed to fetch news items from various RSS feeds and send them directly to a Slack channel. This tool is especially useful for teams who want to stay updated with the latest news from specific sources without leaving their Slack workspace.

## Features

- **Multiple RSS Feeds**: Supports fetching news from multiple RSS feed URLs.
- **Slack Integration**: Sends news items directly to a configured Slack channel.
- **Configurable News Limit**: Set a maximum number of news items to be sent in a single run.
- **Persistence**: Keeps track of which news items have been sent to avoid duplicates.

## Configuration

The application requires a `config.yaml` file with the following structure:

```yaml
rss_feeds:
  - "http://feed1.com/rss"
  - "http://feed2.com/rss"
slack_webhook: "https://hooks.slack.com/services/XXXX/YYYY/ZZZZ"
max_news_count: 5
```

- rss_feeds: List of RSS feed URLs to monitor.
- slack_webhook: Your Slack webhook URL for sending messages.
- max_news_count: The maximum number of news items to send per run.

## Getting Started

### Prerequisites

- Go (version 1.x or later)
- A Slack webhook URL

### Installation

1. Clone the repository:
```shell
git clone https://github.com/ppihus/RSSlackPublisher.git
```
2. Navigate to the project directory:
```shell
cd RSSlackPublisher
```
3. Build the project
```shell
go build
```

### Usage
1. Edit config.yaml with your desired RSS feeds and Slack webhook URL.
2. Run the application:
```shell
./FeedToSlack
```

### Contributing
Contributions to FeedToSlack are welcome! Feel free to open issues or submit pull requests.