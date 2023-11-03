# gallery-service-uplode-get-deliting-photos


**CircleCI CI/CD status** 

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/RSO-project-Prepih/gallery-service-upload-get-deleting-photos/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/RSO-project-Prepih/gallery-service-upload-get-deleting-photos/tree/main)

--- 
## Microservice for uploading, getting and deleting photos.
In this microservice we can upload photos, get photos and delete photos. We can also get all photos from the database and get all photos from the database by user id. 

The end points are: 
- POST /uploadphotos
- GET /getphotos/:user_id
- DELETE /deletephoto/:user_id/:image_id
- GET /live
- GET /ready
- GET /metrics

## How to run the microservice locally
Add the .env file 
```bash
    cp .env.example .env
```
End fix the variables in the .env file

Then run the microservice with the command
```bash
    go run main.go
```
## The /uploadphotos endpoint

To upload photos we need to send a POST request to the /uploadphotos endpoint. The request body must be in JSON format and must contain the following fields:

- username  string
- user_id   uuid     
- image_name string
- data  base64 encoded string

The response will be in JSON format and will contain the following fields:
```JSON
{
    "message": "Photo uploaded successfully"
}
```
Postman example: 
![Postman example](/images/Capture199.JPG)

## The /getphotos/:user_id endpoint
To get photos we need to send a GET request to the /getphotos/:user_id endpoint. The user_id must be uuid.

example: curl -X GET http://localhost:8080/getphotos/123e4567-e89b-12d3-a456-426614174000

The response will be in JSON format and will contain the following fields:
```JSON
{
    "message": "Success - The action was successfully received, understood, and accepted",
    "data": [
        {
            "id": "b3d0a0a0-5b1a-4b1a-8b0a-0a0a0a0a0a0a",
            "username": "test",
            "user_id": "b3d0a0a0-5b1a-4b1a-8b0a-0a0a0a0a0a0a",
            "image_name": "test",
            "data": "data:image/jpeg;base64,/9j/4AAQSkZJR"
            }]
}
```
## The /deletephoto/:user_id/:image_id endpoint
To delete photos we need to send a DELETE request to the /deletephoto/:user_id/:image_id endpoint. The user_id and image_id must be uuid.

example: curl -X DELETE http://localhost:8080/deletephoto/123e4567-e89b-12d3-a456-426614174000/123e4567-e89b-12d3-a456-426614174000

The response will be in JSON format and will contain the following fields:
```JSON
{
    "message": "Photo deleted successfully"
}
```

## The /live endpoint, the /ready endpoint and the /metrics endpoint
This endpoints are for the Kubernetes. The /live endpoint and the /ready endpoint are for the Kubernetes liveness and readiness probes. The /metrics endpoint is for the Kubernetes metrics and Prometheus.

