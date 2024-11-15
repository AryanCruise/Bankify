# Bankify - Banking Microservices API

### Description
**Bankify** is a RESTful banking API microservice application that provides essential banking functionalities such as **account management**, **money transfers**, **notifications**, and **authentication**. It is designed for ease of interaction and is accessible through a user-friendly **Swagger UI**.

### Features
- **Account Management**: 
  - Create, update, and delete accounts.
  - Retrieve account details and balance.

- **Transaction Management**: 
  - Deposit, withdraw, and transfer money between accounts.

- **Authentication**: 
  - Secure JWT-based user authentication to ensure data privacy and security.
 
- **Notification**: 
  - The producer sends data to the Kafka topic, consumed by email and SMS services to send transaction notifications.

- **Swagger Integration**: 
  - Intuitive UI for exploring and testing APIs.

- **PostgreSQL Support**: 
  - Persistent data storage for banking transactions, account details, and user information.

### Technologies Used
- **Programming Language**: Go (Golang)
- **Database**: PostgreSQL
- **API Documentation**: Swagger
- **Deployment**: AWS (EC2)

<br>**Deployed Application(link)**: http://51.20.98.52/swagger/ <br>
<br>Working perfectly on localhost:8080<br>
<br>Facing some issues with the deployed application (Status - working on it)<be>
<br><br>
### API Overview

![Swagger UI Interface](https://github.com/user-attachments/assets/e6a493d6-2a68-4c9d-aded-bdc7873a430c)
*Image: Swagger UI interface with API access, connected to PostgreSQL*

### Kafka-based Email and SMS notifications
![image](https://github.com/user-attachments/assets/669a86c6-d699-46c6-b505-d3517498cfc6)
*Image: Email notification sent to the user after deposit transactions*
<br><br>
![WhatsApp Image 2024-11-15 at 23 47 59_ff6ea253](https://github.com/user-attachments/assets/47598e84-b82a-46b7-9c0d-7eb5c4ac103e)
*Image: SMS notification sent to the user after deposit transaction*
