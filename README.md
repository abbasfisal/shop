
# 🛍️ Shop

A **scalable, session-based web application** built with **Golang**, designed for simplicity and efficiency.

---

## 🚀 Features
- Built with **Golang**, utilizing modern frameworks.
- Supports **job scheduling** and **task queues** using the Asynq package.
- Database migrations and management with Makefile commands.
- Fully configurable via `.env` file.

---

## 📋 Prerequisites
Before starting, ensure you have:
- **Golang** installed (version >= 1.19 recommended).
- **Make** utility installed on your system.
- Access to a **PostgreSQL/MySQL database**.
- A `.env` configuration file. To set it up:
  ```bash
  cp .env.example .env
  ```

---

## 🛠️ Setup and Run

### 1️⃣ Database Configuration
Generate the SQL migrator and database configuration:
```bash
make generate-sql-migrator-dbconfig
```

Apply database migrations:
```bash
make migration-up
```

### 2️⃣ Start the Application
To run the application, execute:
```bash
make run
```

---

## ⚙️ Background Processes

### 🎯 Job Worker
Start the job worker to handle background tasks:
```bash
make start-worker
```

### ⏰ Scheduler
Start the task scheduler for periodic jobs:
```bash
make start-schedule
```
### 🔎 monitoring & administering Asynq task queue
    go to admin route => /monitoring

---

## 🔄 Database Management

### Apply Migrations
Run the following command to apply migrations:
```bash
make migration-up
```

### Roll Back Migrations
To undo the last migration, use:
```bash
make migration-down
```

---

## 🧑‍💻 Contributing
We welcome contributions! To get started:
1. Fork the repository.
2. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```  
3. Commit your changes and open a pull request.

---

## 📜 License
This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.

---

## 🤝 Support
If you have any questions or run into issues, feel free to open an issue or contact us directly.

---

🎉 Happy coding!