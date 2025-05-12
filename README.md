# Notification Service

The Notification Service is a Go-based web application that provides real-time notifications via WebSocket and RESTful APIs.

## Features

- **Ping Endpoint**: Check if the service is running.
- **Create Notifications**: Add new notifications via a POST request.
- **Retrieve Notifications**: Fetch all notifications via a GET request.
- **Delete Notifications**: Remove a notification by ID via a DELETE request.
- **WebSocket Support**: Receive real-time notifications.

## Endpoints

- `GET /ping`: Health check endpoint.
- `POST /notify`: Create a new notification.
- `GET /notifications`: Retrieve all notifications.
- `DELETE /notifications/{id}`: Delete a notification by ID.
- `GET /ws`: Connect to the WebSocket for real-time updates.

## How to Run

1. Clone the repository.
2. Navigate to the project directory.
3. Run the application:
   ```bash
   go run main.go

## Demo

<video width="640" height="360" controls>
  <source src="demo.mp4" type="video/mp4">
  Your browser does not support the video tag.
</video>