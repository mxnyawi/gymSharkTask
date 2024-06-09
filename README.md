# GymShark Package Sale

GymShark Package Sale is a project that allows you to enter five packages of a set amount of orders and an order request. The project calculates and returns the least loss in terms of amount of packages and their worth. An order history is kept on the screen for comparison and validation.

## Prerequisites

Ensure you have the following installed on your local machine:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Installation & Usage

### Backend

The backend is a Go application. It's containerized using Docker, and the image is built automatically when you run `docker-compose up`.

### Frontend

The frontend is a React application. Like the backend, it's containerized using Docker, and the image is built automatically.

## Configuration

The application uses environment variables for configuration. These are read from a `config.env` file at the root of the project.

To set up your configuration:

1. Create a new file named `config.env` in the project root.
2. Add your environment variables in this file, one per line, in the format `VARIABLE_NAME=value`. For example:

    ```env
        BUCKET_NAME=myBucket
        SCOPE_NAME=myScope
        COLLECTION_NAME=myCollection
        DOCUMENT_ID=12345
        USERNAME=admin
        PASSWORD=password
    ```

Replace the values with your actual configuration.

Here's a list of the environment variables used by the application:

- `BUCKET_NAME`: The name of your bucket in the database.
- `SCOPE_NAME`: The name of your scope in the database.
- `COLLECTION_NAME`: The name of your collection in the database.
- `DOCUMENT_ID`: The ID of your order history in the database.
- `USERNAME`: The username to use for database authentication.
- `PASSWORD`: The password to use for database authentication.

Remember not to commit your `config.env` file to the Git repository. It's already listed in `.gitignore` to help prevent this.

### Running the Application

To start the application:

1. Clone the repository:
    ```bash
    git clone https://github.com/mxnyawi/gymSharkTask.git
    ```
2. Navigate to the project directory:
    ```bash
    cd gymSharkTask
    ```
3. Use Docker Compose to build and run the application:
    ```bash
    docker-compose up
    ```

This will start the backend, frontend, and the database. 

You can access the frontend at `http://localhost:3000`. This is where the React application will be running.

The database can be accessed at `http://localhost:8091`. This is useful for database management and viewing data.

Please ensure that these ports are available on your machine before running the application.

## API Endpoints

The application provides the following HTTP API endpoints:

- `POST /registerUser`: Registers a new user. The request body should include the user's details.

- `POST /loginUser`: Authenticates a user and logs them in. The request body should include the user's username and password.

- `POST /createAdminUser`: Creates a new admin user. The request body should include the admin user's details.

- `POST /order`: Creates a new order. The request body should include the order details.

- `POST /setDocument`: Creates a new document in the database. The request body should include the document details.

- `GET /getDocument`: Retrieves a document from the database. The request parameters should include the document ID.

- `PUT /updateDocument`: Updates a document in the database. The request body should include the new document details, and the request parameters should include the document ID.

All endpoints require authentication. This is handled by the `AuthMiddleware` function, which checks for a valid authentication token in the `Authorization` header of the request.

The application is configured to allow Cross-Origin Resource Sharing (CORS) from `http://localhost:3000`. This means that a frontend running on this URL can make requests to the API.

## Contributing

Instructions for contributing to your project.

## License

Information about the project's license.