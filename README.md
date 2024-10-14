# Tikkie Persons Service

## Decisions and rationale:

- I used AWS DocumentDB as the database since it is a managed mongo-compatible service, and Mongo was easiest to setup and use while I was developing locally.
- I used AWS Lambda as the managed serverless service (As I understand this is the service used by Tikkie)
- I used AWS SQS as the message queue since it is easy to setup and use, and it is a managed service.
- I created a dockerfile and docker-compose config to build the lambda function and run the solution components locally. This decision was made before I understood I was also expected to create the Infrastructure config & deploy to AWS.

## Notes on Structure

```
├── app
│   ├── cmd
│   ├── external
│   └── internal
└── infra
    ├── bin
    ├── lib
    └── test
```

- /app contains the service code
- /infra contains the cdk IAC config and scripts to deploy to AWS
- app/external contains packages that could be reused by other projects
- app/internal contains core packages, specific to this project
- app/cmd contains main application entry points (http server, lambda function, cli interface, etc)

Gotchas

- I had to use the `aws-sdk` package to connect to DocumentDB, since the `mongodb` package does not support DocumentDB.
- To build and deploy from my local machine (m1 mac), I set build arguments `GOARCH=amd64` and `GOOS=linux` in the Dockerfile. This prepares the function for the lambda x86_64 architectures.
