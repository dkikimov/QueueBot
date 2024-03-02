## QueueBot: Manage Lab Submission Queues with Telegram

**QueueBot** is a GoLang-based Telegram bot designed to help students in my university group efficiently manage queues for submitting laboratory works. With this bot, students can:

* **Create new queues** for specific lab assignments.
* **Join or leave existing queues** seamlessly.
* **Choose between shuffling** the queue for fairness or **advancing in straight order**.
* See who is **currently passing** a lab work.

**Benefits:**

* **Reduces waiting time:** Queues ensure organized submissions, preventing chaotic rushes.
* **Fairness and transparency:** Shuffling and order options provide equal opportunity.
* **Convenience and accessibility:** Students can manage queues directly through Telegram.
* **Notifications:** Keeps everyone informed about their position and queue updates.

## Getting Started:

1. **Clone the repository:**

   ```bash
   git clone https://github.com/dkikimov/QueueBot.git
   ```

2. **Install dependencies:**

   ```bash
   cd QueueBot
   go mod download
   ```

3. **Build and run the bot:**

   ```bash
   go build cmd/main.go
   BOT_TOKEN={your_token} DEBUG={true or false} DATABASE_PATH={path} ./main
   ```

### Docker way

1. **Build the image:**

   ```bash
   docker build . -t queue-bot
   ```
2. **Create .env file:**

   ```
   BOT_TOKEN=your_token
   DEBUG=true     (false is default)
   DATABASE_PATH=path
   ```

3. **Create and run container:**

   ```bash
   docker run --env-file .env --rm queue-bot 
   ```
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.
* If you have suggestions for adding or removing projects, feel free to [open an issue](https://github.com/dkikimov/QueueBot/issues/new) to discuss it, or directly create a pull request after you edit the *README.md* file with necessary changes.
* Create individual PR for each suggestion.

### Creating A Pull Request

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
