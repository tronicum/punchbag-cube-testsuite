# Platform TODOs & Milestones

## 1. Expand Resource Coverage
- [ ] Add AWS resources: S3, CloudWatch, IAM, etc. (simulation + proxy)
- [ ] Add GCP resources: CloudSQL, Pub/Sub, IAM, etc. (simulation + proxy)
- [ ] Implement CRUD and validation for all new resources

## 2. Enhance Validation and Error Handling
- [ ] Add schema and required field checks for all resources
- [ ] Improve error messages and logging in simulation/proxy servers

## 3. Authentication & Authorization
- [ ] Add authentication (API keys, JWT, or OAuth) to all endpoints
- [ ] Implement role-based access control if needed

## 4. Dynamic Routing & Real Cloud Integration
- [ ] Allow dynamic switching between simulation and real cloud APIs
- [ ] Integrate real cloud SDKs/APIs for proxy/execute endpoints

## 5. CI/CD & Automated Testing
- [ ] Add end-to-end tests for all API flows (simulate, proxy, execute, validation)
- [ ] Integrate with CI/CD to run tests on PRs and releases

## 6. Documentation & Examples
- [ ] Expand documentation with usage examples for all endpoints and CLI
- [ ] Provide sample workflows for local dev, CI, and real cloud

## 7. User Experience
- [ ] Add interactive CLI features (resource discovery, auto-complete, wizards)
- [ ] Improve output formatting and error reporting for CLI and API
