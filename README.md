# Tikkie Persons Service

## Decisions and rationale:

- I used AWS DocumentDB as the database since it is a managed mongo-compatible service, and Mongo was easiest to setup and use while I was developing locally.
- I used AWS Lambda as the managed serverless service (As I understand this is the service used by Tikkie)
- I used AWS SQS as the message queue since it is easy to setup and use, and it is a managed service.
- I created a dockerfile and docker-compose config to build the lambda function and run the solution components locally. This decision was made before I understood I was also expected to create the Infrastructure config & deploy to AWS.

## Notes on Structure

- In /external are packages that can be reused by other projects
- In /internal are packages only used by this project
