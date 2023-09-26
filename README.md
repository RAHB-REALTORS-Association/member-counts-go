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

```sh
git clone https://github.com/RAHB-REALTORS-Association/member-counts-go.git
cd member-counts-go
go build
```

Before running, ensure that the environment variables are correctly set, either in your environment or in a `.env` file in the project directory.

## 🧑‍💻 Usage

Once configured and built, run the application using:

```sh
./member-counts-go
```

The application will refresh the specified Redash query, fetch the resulting data, process it to calculate the total member counts, and send this count to the specified Google Chat Webhook.

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
