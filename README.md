# Location Tracking Application

This application allows users to track locations with support for login via username/password or Google OAuth.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing
purposes.

### Prerequisites

Before running the application, make sure you have the Protocol Buffer Compiler installed on your platform. The
application uses `.proto` files to define message types and service interfaces.

### Installing

Follow these steps to get your development environment running:

1. **Compile Protocol Buffers**

   Navigate to the `location` directory where the `location.proto` file is located and run the following command:

```bash
protoc --go_out=. --go_opt=paths=source_relative
--go-grpc_out=. --go-grpc_opt=paths=source_relative
location.proto
```

This command generates the Go code for the protocol buffer message types.

2. **Running the Application**

```bash 
go build && ./geolocation-tracking
 ```

#### but make sure you have the environment variables set up as per step 3

3. **Configuration for Google Maps and Google OAuth**

- **Configure Google Maps API Key**

  In the `maps.html` file, replace `YOUR_API_KEY_HERE` with a valid Google Maps API key to display the map correctly.
  Obtain an API key by following the instructions provided in
  the [Google Cloud documentation](https://cloud.google.com/maps-platform/).

  Example snippet from `maps.html`:

  ```html
  <!-- Replace "YOUR_API_KEY_HERE" with your actual Google Maps API key -->
  <img src="https://maps.googleapis.com/maps/api/staticmap?center=New+York,NY&zoom=13&size=600x300&maptype=roadmap&key=YOUR_API_KEY_HERE" alt="Google Map of New York" class="map-thumbnail">
  ```

- **Set Up Google OAuth**

  If you want to log in using Google OAuth, you need to set up OAuth credentials in the Google Cloud Platform and
  export `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` environment variables with your credentials. Follow the
  instructions in the [Google Identity Platform documentation](https://developers.google.com/identity) to obtain these
  credentials.

  Export the environment variables in your development environment:

  ```bash
  export GOOGLE_CLIENT_ID='your_google_client_id_here'
  export GOOGLE_CLIENT_SECRET='your_google_client_secret_here'
  ```

  Make sure these variables are available in the environment where your application is running.

## Using the Application

### Login Options

You can log in to the application using one of the following methods:

- **Username/Password**: Use the credentials provided within the code. This method is primarily for demonstration
  purposes.

- **Google OAuth**: Log in using your Google account to authenticate. Ensure you have set up Google OAuth credentials
  and configured them in the application.

## Built With

- [Protocol Buffers](https://developers.google.com/protocol-buffers) - Interface description language
- [Go](https://golang.org) - Programming language used
- [Gin](https://github.com/gin-gonic/gin) - Web framework used for building the HTTP server

## Authors

- **Vinod Halaharvi** - *Initial work*

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Want to hire me for your Golang needs?

- message me on Linkedin https://www.linkedin.com/in/vinod-halaharvi-289a1a13/
- or email me at vinod@smartify.software
