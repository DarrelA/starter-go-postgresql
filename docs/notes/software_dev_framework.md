<a name="readme-top"></a>

- [Middlewares](#middlewares)
  - [Logging and Monitoring](#logging-and-monitoring)
    - [Custom Headers: `Correlation-ID` and `Request-ID`](#custom-headers-correlation-id-and-request-id)
- [References](#references)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Middlewares

| **Middleware**                                        | **Use Case(s)**                                                                                                                                                    |
| ----------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Handling CORS<br/>(Cross-Origin Resource Sharing)** | Allowing or restricting resources on a web server to be requested from another domain.<br/>Configuring CORS policies to secure and control access to resources.    |
| **Authentication and Authorization**                  | Ensuring users are authenticated before accessing certain resources.<br/>Validating user permissions and roles to control access to specific endpoints or actions. |
| **Logging and Monitoring**                            | Tracking and recording requests and responses for auditing and debugging.<br/>Monitoring application performance and detecting issues or anomalies.                |
| **Data Validation and Sanitization**                  | Ensuring incoming request data adheres to expected formats and rules.<br/>Removing or escaping potentially harmful data to prevent injection attacks.              |
| **Session Management**                                | Maintaining user session state across multiple requests.<br/>Handling session creation, expiration, and validation.                                                |

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Logging and Monitoring

| **Information**             | **Use Case (Software Engineers)**                                                                                                                   | **Use Case (DevOps)**                                                                                                                                             |
| --------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Host/</br>Server Info**   | Identifying issues related to specific servers or instances in a distributed environment.                                                           | Monitoring server health and performance to identify potential issues and balance load across instances.                                                          |
| **Request Method**          | Identifying the type of HTTP request (GET, POST, etc.) to understand user actions and flow through the application.                                 | Monitoring the distribution of different request types to identify unusual patterns or spikes that may indicate issues.                                           |
| **HTTP Referrer**           | Understanding the sources of traffic and user navigation patterns to optimize user experience and marketing efforts.                                | Tracking referral sources to identify potential security threats, such as phishing or spam attacks, and to optimize resource allocation based on traffic sources. |
| **Request URL**             | Debugging issues related to specific endpoints and understanding which resources are being accessed.                                                | Tracking usage patterns across different endpoints to optimize performance and resource allocation.                                                               |
| **Status Code**             | Identifying success or failure of requests to pinpoint potential issues in the application logic.                                                   | Monitoring overall application health and identifying trends in error rates or unusual spikes in specific status codes.                                           |
| **Response Time**           | Measuring the performance of individual requests to identify slow responses and potential bottlenecks in the application.                           | Ensuring the application meets performance SLAs and identifying areas for optimization to improve response times.                                                 |
| **Latency**                 | Measuring the time taken for different parts of the request processing to identify and optimize slow components.                                    | Monitoring end-to-end latency to ensure the application meets performance SLAs and identifying areas for optimization.                                            |
| **IP Address**              | Tracking user activity and identifying potential malicious activities or misuse of the application.                                                 | Monitoring the geographic distribution of traffic and identifying potential security threats such as DDoS attacks.                                                |
| **User Agent**              | Understanding the types of clients (browsers, mobile devices, etc.) accessing the application to ensure compatibility and optimize user experience. | Tracking the mix of client types to plan for capacity and performance testing, as well as identifying potentially malicious clients.                              |
| **Correlation ID**          | Tracing requests through the entire application stack to diagnose issues and ensure consistency across services.                                    | Linking logs from different services and components to get a holistic view of the application's behavior and identify root causes of issues.                      |
| **Request ID**              | Tracking and debugging specific requests within a service to identify issues and optimize request handling.                                         | Monitoring individual request logs to identify patterns, performance issues, and potential security threats at a granular level.                                  |
| **User ID/</br>Session ID** | Tracking user activity and debugging issues related to specific users or sessions.                                                                  | Monitoring user behavior and session patterns to identify potential misuse or performance issues affecting specific user groups.                                  |

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Custom Headers: `Correlation-ID` and `Request-ID`

- **`Correlation-ID`**
  - **Purpose:** Used to trace and correlate a series of related requests across multiple services. It helps in understanding the journey of a particular transaction or user action through the system.
  - **Scope:** Typically spans multiple services and systems. It's used to connect different parts of a transaction that might pass through various microservices.
  - **Usage:** The `Correlation-ID` header is included in HTTP requests and responses to maintain the correlation ID across different services.
- **`Request-ID`**
  - **Purpose:** Used to uniquely identify a single request within a service. It helps in logging and debugging individual requests.
  - **Scope:** Usually limited to a single request-response cycle within a service. It can be useful for tracking the processing of a request through different components of the same service.
  - **Usage:** The `Request-ID` header is included in HTTP requests and responses to identify the specific request.
- **Key Differences**
  - **Scope:** `Correlation-ID` has a broader scope, typically spanning multiple services, while `Request-ID` is more narrowly focused on a single request within a service.
  - **Usage Scenario:** Use `Correlation-ID` when you need to trace an end-to-end transaction through multiple microservices. Use `Request-ID` for tracking and debugging individual requests within a specific service.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# References

- [Hypertext Transfer Protocol (HTTP) Field Name Registry](https://www.iana.org/assignments/http-fields/http-fields.xhtml)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
<p align="right">(<a href="../README.md">back to main</a>)</p>

---
