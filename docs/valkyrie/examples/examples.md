# Valkyrie API Reference

Welcome to the Valkyrie API reference documentation. This document provides a brief overview and introduction to using the Valkyrie API effectively.

## Introduction

Valkyrie provides a set of APIs designed to interact with its robust backend services, enabling seamless integration and functionality for your applications. This guide introduces key concepts, authentication methods, and endpoints available in the Valkyrie API.

## Base URL

All Valkyrie API endpoints are accessed via the following base URL:

```
https://valkyrie.evnix.cloud
```

Ensure all your API requests start with this URL.

## Authentication

Valkyrie API uses API tokens for authentication. Include your API token in request headers as shown:

```bash
curl -X GET https://valkyrie.evnix.cloud \
     -H "Authorization: Bearer YOUR_API_TOKEN"
```

Replace `YOUR_API_TOKEN` with your actual API token.

## API Endpoints

Explore individual endpoints in the [Valkyrie API documentation](API.md) for detailed endpoint descriptions and usage examples.

* [**Get all execution jobs**](API.md)
* [**Get all executionss**](API.md)
* [**Get all language versions**](API.md)
* [**Get all languages**](API.md)
* [**Get execution config**](API.md)
* [**Get execution job**](API.md)
* [**Get execution result by id**](API.md)
* [**Get executions of given job**](API.md)
* [**Get language by ID**](API.md)
* [**Get language version by ID**](API.md)
* [**Get version**](API.md)
* [**Health Check**](API.md)
* [**Fetch Flake**](API.md)
* [**Execute a script**](API.md)
* [**Delete execution job**](API.md)
* [**Create a sandbox**](API.md#create-a-sandbox)
* [**Cancel Execution Job**](API.md#cancel-execution-job)


## Error Handling

The Valkyrie API returns standard HTTP status codes for indicating success or failure:

* `200`: Successful request.
* `400`: Bad request (validation failed or incorrect parameters).
* `401`: Unauthorized (invalid or missing authentication).
* `404`: Not found (endpoint or resource doesn't exist).
* `500`: Internal server error (server-side issue).

Always check the response status code and handle errors accordingly in your application.
