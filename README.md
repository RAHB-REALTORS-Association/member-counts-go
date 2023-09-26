[![Continuous Integration](https://github.com/RAHB-REALTORS-Association/member-counts-go/actions/workflows/go.yml/badge.svg)](https://github.com/RAHB-REALTORS-Association/member-counts-go/actions/workflows/go.yml)
[![Build](https://github.com/RAHB-REALTORS-Association/member-counts-go/actions/workflows/build.yml/badge.svg)](https://github.com/RAHB-REALTORS-Association/member-counts-go/actions/workflows/build.yml)
[![Docker](https://github.com/RAHB-REALTORS-Association/member-counts-go/actions/workflows/docker.yml/badge.svg)](https://github.com/RAHB-REALTORS-Association/member-counts-go/actions/workflows/docker.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Member Counts

**Member Counts** is a Go-based application designed to fetch and process member data from Redash and send total member counts to a Google Chat using a Webhook.

## 📖 Table of Contents
- [⚙️ Configuration](#️-configuration)
- [🚀 Deployment](#-deployment)
- [🧑‍💻 Usage](#-usage)
- [🛠️ Tech Stack](#️-tech-stack)
- [🌐 Community](#-community)
- [📄 License](#-license)

## ⚙️ Configuration

Environment variables are used to configure the application. They can be set in a `.env` file placed in the project directory or directly in the environment:

- `REDASH_BASE_URL`: The base URL of your Redash instance.
- `REDASH_API_KEY`: The API key for accessing the Redash API.
- `REDASH_QUERY_ID`: The ID of the Redash query to fetch.
- `GOOGLE_CHAT_WEBHOOK_URL`: The URL of the Google Chat Webhook.
- `SCHEDULE_HOUR`: The hour to trigger the action.
- `SCHEDULE_MINUTE`: The minute to trigger the action.
- `TIMEZONE`: Your timezone.

## 🚀 Deployment

### Deploy Using PaaS
[![Deploy to DO](https://www.deploytodo.com/do-btn-blue.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/RAHB-REALTORS-Association/member-counts-go/tree/main)

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/RAHB-REALTORS-Association/member-counts-go/tree/main)

### Build and Run Locally
```sh
git clone https://github.com/RAHB-REALTORS-Association/member-counts-go.git
cd member-counts-go
go build
```
Before running, ensure that the environment variables are correctly set, either in your environment or in a `.env` file in the project directory.

### Run Pre-built Docker Image
```sh
docker pull ghcr.io/rahb-realtors-association/member-counts-go:latest
docker run --env-file .env ghcr.io/rahb-realtors-association/member-counts-go:latest
```
Ensure that the `.env` file containing the environment variables is present in the directory from which you are running the Docker command.

### Building Docker Image Locally
```sh
git clone https://github.com/RAHB-REALTORS-Association/member-counts-go.git
cd member-counts-go
docker build -t member-counts-go .
docker run --env-file .env member-counts-go
```
Again, ensure the `.env` file is present in your project directory and contains the necessary environment variables before running the Docker image.

## 🧑‍💻 Usage

### Scheduled Execution
To run the application in a scheduled manner, use the following command:

```sh
./member-counts-go
```

Or, when using Docker:

```sh
docker run --env-file .env ghcr.io/rahb-realtors-association/member-counts-go:latest
```

In these scenarios, the application will remain idle until it reaches the specified hour and minute in your environment or `.env` file. Once the specified time is hit, the application will:

1. **Refresh** the specified Redash query.
2. **Fetch** the newly generated data.
3. **Process** the data to calculate the total member counts.
4. **Send** this count to the specified Google Chat Webhook.

### Immediate Execution

For FaaS environments or testing purposes, if you prefer to run the application immediately without waiting for the scheduled time, use the `--now` flag:

```sh
./member-counts-go --now
```

Or, when using Docker:

```sh
docker run --env-file .env ghcr.io/rahb-realtors-association/member-counts-go:latest --now
```

With the `--now` flag, the application bypasses the scheduler and executes the aforementioned steps immediately. This is especially useful for testing or when deploying in environments that are ephemeral, like certain FaaS platforms.

### Environment Configuration

Before running the application in any environment, ensure that the environment variables are correctly set, either in your environment or in a `.env` file in the project directory or the directory from which you are running the Docker command.

## 🛠️ Tech Stack

The technologies used in this project include:

- Go 1.17+ 🌿
- Redash API 📊
- Google Chat API 💬

## 🌐 Community

### Contributing 👥

Contributions to Member Counts are warmly welcomed. Please feel free to submit pull requests or open issues to discuss potential modifications or improvements. See our [Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/) for community guidelines.

[![Submit a PR](https://img.shields.io/badge/Submit_a_PR-GitHub-%23060606?style=for-the-badge&logo=github&logoColor=fff)](https://github.com/RAHB-REALTORS-Association/member-counts-go/compare)

### Reporting Bugs 🐛

Encountered a bug or have a feature suggestion? Please [raise an issue](https://github.com/RAHB-REALTORS-Association/member-counts-go/issues/new/choose) on GitHub with relevant details.

[![Raise an Issue](https://img.shields.io/badge/Raise_an_Issue-GitHub-%23060606?style=for-the-badge&logo=github&logoColor=fff)](https://github.com/RAHB-REALTORS-Association/member-counts-go/issues/new/choose)

## 📄 License
This project is open sourced under the MIT license. See the [LICENSE](LICENSE) file for more info.
