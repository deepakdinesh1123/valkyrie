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

The Valkyrie API endpoints are grouped into logical sections for ease of navigation and use:

* **Users**: Manage user-related operations.
* **Payments**: Handle payment transactions and related actions.
* **Analytics**: Retrieve analytical data and insights.
* **Notifications**: Manage notifications sent to users.

Explore individual sections in the [Valkyrie API documentation](https://docs.riza.io/api-reference/introduction) for detailed endpoint descriptions and usage examples.

## Error Handling

The Valkyrie API returns standard HTTP status codes for indicating success or failure:

* `200`: Successful request.
* `400`: Bad request (validation failed or incorrect parameters).
* `401`: Unauthorized (invalid or missing authentication).
* `404`: Not found (endpoint or resource doesn't exist).
* `500`: Internal server error (server-side issue).

Always check the response status code and handle errors accordingly in your application.
