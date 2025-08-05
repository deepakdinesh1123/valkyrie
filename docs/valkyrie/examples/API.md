## Cancel Execution Job

Cancel a running or queued execution job.

### Endpoint

```http
PUT {{BASE_URL}}/api/executions/jobs/:JobId
```

### Authentication

* **Bearer Token**

### Request Body

* **None**

### Path Parameters

| Name  | Type   | Description                            |
| ----- | ------ | -------------------------------------- |
| JobId | string | The ID of the execution job to cancel. |

### Responses

* **204 No Content**: Job cancelled successfully.
* **404 Not Found**: No job exists with the specified `JobId`.
* **4XX / 5XX**: Other error conditions.

---

### Example

#### Curl

```bash
export BASE_URL="http://localhost:3000"
export JOB_ID="123e4567-e89b-12d3-a456-426614174000"

curl -X PUT "${BASE_URL}/api/executions/jobs/${JOB_ID}"
```

Feel free to adjust `BASE_URL` and `JOB_ID` as needed.

## Create a Sandbox

Create a new sandbox environment with specified configurations.

### Endpoint

```http
POST {{BASE_URL}}/api/sandbox
```

### Authentication

* **None**

### Request Body

The request body should be a JSON object containing sandbox configuration details:

```json
{
  "nix_flake": "",
  "languages": [],
  "system_dependencies": [],
  "services": []
}
```

### Body Parameters

| Name                  | Type   | Description                             |
| --------------------- | ------ | --------------------------------------- |
| `nix_flake`           | string | Nix flake configuration (optional).     |
| `languages`           | array  | List of programming languages required. |
| `system_dependencies` | array  | List of system dependencies required.   |
| `services`            | array  | List of additional services to include. |

### Responses

* **201 Created**: Sandbox successfully created.
* **400 Bad Request**: Invalid request parameters.
* **4XX / 5XX**: Other error conditions.

---

### Example

#### Curl

```bash
export BASE_URL="http://localhost:3000"

curl -X POST "${BASE_URL}/api/sandbox" \
     -H "Content-Type: application/json" \
     -d '{
           "nix_flake": "",
           "languages": ["python"],
           "system_dependencies": ["git"],
           "services": ["postgresql"]
         }'
```

Adjust the request parameters (`languages`, \`system\_dependencies
