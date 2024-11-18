# Redis-Clone: Building Redis from Scratch in Go

Welcome to **Go-Redis**, a personal project where I attempt to build a Redis-like database from scratch using Go and mostly standard libraries. This project is inspired by the incredible Redis project and is a testament to the brilliant work of the Redis developers. 🎉

**Note:** This project is meant for fun and educational purposes only. It is not intended to be used as a production-ready database.

## Features

- **RESP Protocol:** Implements the Redis Serialization Protocol (RESP) to communicate between the server and clients. RESP is a simple, efficient, and easy-to-implement protocol that supports various data types and commands.
  
- **TCP Server:** Handles client connections and requests over TCP. The server listens for incoming connections, processes commands, and sends back appropriate responses, just like Redis.

- **Replication:** Supports basic replication between instances. This feature allows one server to act as a primary and others as replicas, ensuring data consistency and redundancy across multiple nodes.

- **Key Expiry:** Implements key expiration functionality. Keys can be set to expire after a certain period, automatically removing them from the database when their time is up, helping manage memory usage effectively.

- **Transactions (TODO):** Planned for future implementation. Transactions will allow the execution of a series of commands atomically, ensuring that either all commands are executed or none are, maintaining data integrity.

- **RDB Persistence (TODO):** Planned for future implementation. This feature will provide a mechanism to save the in-memory dataset to disk and restore it upon restart, ensuring data durability.

## Project Status

🚧 **Under Construction:** This project is currently being built and is not fully optimized. I am following a "make it work, then optimize" mindset. Expect bugs, inefficiencies, and incomplete features.

## Huge Respects

A huge shoutout to the [Redis](https://redis.io/) project and its developers for their phenomenal work. This project wouldn't exist without their inspiration.

## About Me

I'm a software engineering student who loves programming. This project is a fun way for me to dive deep into the internals of databases and learn by doing.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details. You are essentially allowed to do whatever you want with this repo's code :)

---

Have fun exploring the internals of Redis, and happy coding! 😄
