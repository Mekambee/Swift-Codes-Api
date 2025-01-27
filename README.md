## Swift-Codes-API 

## Remitly 2025 Internship Recruitment Task



## 1. Used Technologies

- **Go (Golang) 1.20+**
    - [**Gin**](https://github.com/gin-gonic/gin) Framework for REST API
    - [**excelize**](https://github.com/xuri/excelize) for parsing `.xlsx` files
- **PostgreSQL** as the database
- **Docker** + **Docker Compose** for containerization
- **Tests**
    - Built-in Go testing (`testing`) plus [**testify**](https://github.com/stretchr/testify) for assertions

---

## 2. Application Set-up and Running
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/Mekambee/Swift-Codes-Api
2. **Install Docker** and **Docker Compose**. (If you don't have them already) 
   - For Docker, follow the instructions at: [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/)
   - For Docker Compose, follow the instructions at: [https://docs.docker.com/compose/install/](https://docs.docker.com/compose/install/)

3. **In the project root directory, run:**
   ```bash
   docker-compose up --build
   ```
   After the process, the application should be available and running at:
   ```
   http://localhost:8080
   ```
-   **Docker Compose** will pull (or fetch from its cache) the base images (`golang:1.20-alpine`, `alpine:3.17`, `postgres:14`, etc.), compile the application in the builder stage, then start both the PostgreSQL container and your application container. No local installation of Go or PostgreSQL is needed—only Docker and Docker Compose.
    
-   On first startup, the app automatically creates and migrates the `swift_codes` table. 
    After such setup, both the API and the database should be running, and you can interact with the API
    at the localhost:8080 address. The instruction below will guide you through the API endpoints and 
    you will find the guidelines on how to pass the data to the API.

#### ⚠️ "For the application to function properly, ensure that your port :5432 is not occupied, as this is the default port on which the PostgreSQL database listens. (And ofcourse port :8080, where the application is running)"

---

## 3. Endpoints


### 3.1. GET `/v1/swift-codes/{swiftCode}`

Returns details for a single SWIFT code.
If the code ends with `XXX` (headquarter), includes a `branches` array in the response.

Example Response Structure, if Swift Code ends with "XXX" (Bank headquarter)

  ```json
  {
    "address": "string",
    "bankName": "string",
    "countryISO2": "string",
    "countryName": "string",
    "isHeadquarter": true,
    "swiftCode": "ABCDEF12XXX",
    "branches": [
      {
        "address": "string",
        "bankName": "string",
        "countryISO2": "string",
        "isHeadquarter": false,
        "swiftCode": "ABCDEF12XYZ"
      }
    ]
  }
  ```
Example Branch Response Structure, if Swift Code suffix is different than "XXX" (Branch)

  ```json
  {
    "address": "string",
    "bankName": "string",
    "countryISO2": "string",
    "countryName": "string",
    "isHeadquarter": false,
    "swiftCode": "ABCDEF12XYZ"
  }
  ```

### 3.2. GET `/v1/swift-codes/country/{countryISO2}`

Returns all SWIFT codes data (both HQ and branches) for a specific ISO2 country code.

Example Response Structure

  ```json
  {
    "countryISO2": "PL",
    "countryName": "POLAND",
    "swiftCodes": [
      {
        "address": "string",
        "bankName": "string",
        "countryISO2": "PL",
        "isHeadquarter": true,
        "swiftCode": "ABCDEFPLXXX"
      }
    ]
  }
  ```

### 3.3. POST `/v1/swift-codes/`

Adds a new SWIFT code to the database.
In order to add a new record, you must send a JSON body 
in an HTTP `POST` request. For example:

  ```json
{    
  "address": "123 Sample Street, City",
  "bankName": "Example Bank",
  "countryISO2": "us",
  "countryName": "united states",
  "isHeadquarter": true,
  "swiftCode": "EXAMUSNYXXX"
}
  ```

Response structure, when the record is successfully added:
  ```json
  {
    "message": "Swift code created"
  }
  ```

### 3.4. DELETE `/v1/swift-codes/{swiftCode}`

Deletes a record if `swiftCode`, `bankName`, and `countryISO2` match the database entry.

Example request for deleting a record, with `bankName` and `countryISO2` data in RequestBody:
```bash
DELETE /v1/swift-codes/TESTUSNYXXX
Content-Type: application/json

{
  "bankName": "TEST BANK",
  "countryISO2": "US"
}
```
### ⚠️ Important info about my implementation of the DELETE request endpoint:

The original requirement from exercise instruction states:

"Deletes swift-code data if swiftCode, bankName, and countryISO2 match the one in the database."

However, it does not specify the request format for bankName and countryISO2.
In this project, I interpret the requirement by putting bankName and countryISO2 in the request body (JSON) as you can see
above (The example request for deleting a record).

Response Structure, when successfully deleted:
  ```json
  {
    "message": "Swift code deleted"
  }
```
If all three fields (`swiftCode` in path, `bankName`, and `countryISO2` in body) match a record in the database, the record is deleted. Otherwise, a 404 Not Found is returned.

Rationale: This satisfies the specification that `bankName` and `countryISO2` must match, while `swiftCode` is already given in the path parameter.

### 🟢 Fifth additional endpoint (not specified in exercise instruction), which allows user to pass the .xlsx via POST request
#### It was added in order to facilitate testing and experimentation with various .xlsx files.
#### The .xlsx file must follow the same column structure (e.g., COUNTRY ISO2 CODE, SWIFT CODE, etc.). Any row with fewer than the required columns is skipped.
#### After sending a successful request, the user will be able to retrieve the data from sent .xlsx file in the API, because it will be immediately parsed and stored in the database.
### 3.5. POST `/v1/swift-codes/import`

Allows uploading an `.xlsx` file via multipart/form-data for bulk import.

Usage Example:

```bash
POST /v1/swift-codes/import
(form-data: key="file", attach an .xlsx)
```

Response:
```json
{
  "message": "Import successful"
}
```

---
## 4. Example requests structures for testing - How to interact with the API using cURL

You can obviously send requests in any way you like (e.g. Api Testing Apps like: Postman, Talend API Tester, ....), but here are a few examples using cURL, you can just execute by pasting them into your terminal 😊

#### Those commands below only demonstrate how to interact with each endpoint using cURL, with example data

- ### POST request into `/v1/swift-codes/import` - Data load into your API
```bash
curl -X POST http://localhost:8080/v1/swift-codes/import \
  -F "file=@/path/to/your/testfile.xlsx"
```

- ### GET request into `/v1/swift-codes/{swift-code}` - Get an information response for single SWIFT code ("ABCDPLPWXXX" example)
```bash
curl -X GET http://localhost:8080/v1/swift-codes/ABCDPLPWXXX 
```
- ### GET request into `/v1/swift-codes/country/{countryISO2}` - Get an information response for all SWIFT codes for a specific ISO2 country code ("PL" example)
```bash
curl -X GET http://localhost:8080/v1/swift-codes/country/PL
```
- ### POST request into `/v1/swift-codes/` - Add a new SWIFT code record to the database
```bash
curl -X POST http://localhost:8080/v1/swift-codes/ \
  -H "Content-Type: application/json" \
  -d '{
    "address": "Some Address",
    "bankName": "New Bank",
    "countryISO2": "us",
    "countryName": "united states",
    "isHeadquarter": true,
    "swiftCode": "NEWBUSNYXXX"
  }'
```
- ### DELETE request into `/v1/swift-codes/{swiftCode}` - Delete a record from the database
```bash
curl -X DELETE http://localhost:8080/v1/swift-codes/NEWBUSNYXXX \
  -H "Content-Type: application/json" \
  -d '{
    "bankName": "New Bank",
    "countryISO2": "US"
  }'
```


## 5. How the Parser Works

- Uses `excelize` to read `.xlsx` files.
- Required columns: `COUNTRY ISO2 CODE`, `SWIFT CODE`, `CODE TYPE`, `NAME`, `ADDRESS`, `TOWN NAME`, `COUNTRY NAME`, `TIME ZONE`.
- Only significant, selected columns are stored; others (e.g., `CODE TYPE`, `TIME ZONE`) are ignored.
- Merges `address` and `town` into one `Address` field.
- If a row has fewer than 8 columns, it is skipped.
- If any column is empty (e.g., `address`), it becomes an empty string—rows are not discarded if other key fields exist.

---

## 6. Test Overview

Below is a summary of all test categories in this project, along with their purposes and typical file locations. All tests can be executed simply by running:

```bash
go test ./... -v
```

Passing the `-v` flag will show verbose output, including the names of each test as they run.
19 tests should be run and passed while executing the command above. 

⚠️ Warning: Running tests locally via `go test ./... -v` requires `Go` to be installed on your machine. It can be installed from [here](https://golang.org/doc/install).


⚠️ Important Testing Information

During the testing of final version of this project on other computer with Windows system, some issues were observed when running all tests using the `go test ./...` command:

Behavior on Windows:

Running `go test ./...` caused 2 tests to fail:

- TestSwiftService_Basic
- TestIntegration_GetCountry

To resolve this, it was necessary to navigate to the specific directory containing the tests and execute go test directly within that directory. When running tests this way, all tests passed successfully.

On macOS (Where the project was developed), running `go test ./...` executed all tests successfully without any failures.

Recommendation:
If you encounter issues with failing tests on a Windows (or possibly other operating
system), try the following steps:

Navigate to the directory containing the specific tests, and then run `go test -v` directly in that directory.

This difference in behavior may be related to system-specific configurations or dependencies.

## 6.1) Parser Unit Tests 
### (`internal/services/parser_test.go`)
#### There are multiple .xlsx samples in the project (`internal/services/testdata/`) to verify every condition mentioned below.

These tests validate the `.xlsx` parsing logic in isolation, ensuring corner cases are handled:

#### **TestParseSwiftXLSX_Basic**
- Verifies that a row without `XXX` suffix is interpreted as a branch and parsed correctly.

#### **TestParseSwiftXLSX_HQ**
- Checks if a SWIFT code ending with `XXX` is marked as a headquarter.

#### **TestParseSwiftXLSX_CompletelyEmptyFile**
- Ensures that an entirely empty file yields no records.

#### **TestParseSwiftXLSX_OnlyHeaders**
- Ensures that a file with only headers (no data rows) yields no records.

#### **TestParseSwiftXLSX_MissingColumns**
- Tests skipping rows that have fewer than the required columns.

#### **TestParseSwiftXLSX_Uppercase**
- Ensures `countryISO2` and `countryName` get converted to uppercase.

#### **TestParseSwiftXLSX_InvalidFile**
- Confirms that a non-Excel file raises an error.

#### **TestParseSwiftXLSX_SpecialCharacters**
- Checks that special diacritics are preserved in the address.

#### **TestParseSwiftXLSX_MixedCase**
- Verifies mixed-case `countryISO2` / `countryName` are forced to uppercase.

---

## 6.2) Integration Tests

These tests load data into PostgreSQL , then call the actual REST endpoints to confirm everything works end-to-end.
#### There is .xlsx sample file in the project (`test/testdata/integration_test.xlsx/`) which holds about 150 records excerpted from the original file. Integration tests uses it in order to test the API.
### `test/integration_test.go` (Main Integration Tests)

#### **TestIntegration_GetSwiftCode**
- Connects to DB and truncates the table.
- Parses `integration_test.xlsx`, saves records.
- Calls `GET /v1/swift-codes/{swiftCode}`.
- Checks if HQ or branch details are correct.

#### **TestIntegration_GetCountry**
- Connects/truncates DB.
- Loads `integration_test.xlsx`, saves records.
- Calls `GET /v1/swift-codes/country/CL` (as an example).
- Expects a `200` response with an array of SWIFT codes in `CL`.

#### **TestIntegration_PostNewSwiftCode**
- Truncates the DB.
- Calls `POST /v1/swift-codes/` with JSON body.
- Verifies the newly created record in the DB.

#### **TestIntegration_DeleteSwiftCode**
- Inserts a record into `swift_codes`.
- Calls `DELETE /v1/swift-codes/{swiftCode}` with JSON containing `bankName` and `countryISO2`.
- Confirms the record was actually removed.

#### **TestImportSwiftCodesHandler**

Located in a file such as `internal/api/import_handler_test.go` checks the `POST /v1/swift-codes/import` endpoint by:

- Uploading a sample `.xlsx` stored in (`internal/api/testdata`) file using `multipart/form-data`.
- Ensuring records are saved to the database.
- Verifying that the endpoint returns a success message.

---

## 6.3) Router Test (`internal/api/router_test.go`)

### **TestSetupRouter**
- Instantiates the router via `SetupRouter()` to verify routes are registered.
- Confirms that `GET /v1/swift-codes/:swiftCode` (and others) exist.

---

## 6.4) Swift Service Test (`internal/services/swift_service_test.go`)

### **TestSwiftService_Basic**
- Directly tests methods like `SaveSwiftCodes`, `GetSwiftCode`, `GetBranchesByHQ`, and `DeleteSwiftCode` in `swift_service.go`.
- It connects to the DB, inserts sample data, checks retrieval, and finally confirms deletion.
- This ensures the lower-level service layer is covered, outside of the full REST flow.

---

### ⚠️ Integration test usage will clear the data currently stored in DB!
### ⚠️ No test require any data stored in DB beforehand; they are self-contained and work with data samples stored int the mentioned directions, or they create own data

---

## 7. Additional Notes

- **JSON Key Order**: Although the response structures follow the exact fields required by the specification, JSON objects do not guarantee a specific ordering of keys. For example, in the requirement, a headquarter SWIFT code response might be described as:
  ```json
  {
    "address": "string",
    "bankName": "string",
    "countryISO2": "string",
    "countryName": "string",
    "isHeadquarter": true,
    "swiftCode": "ABCDEF12XXX",
    "branches": [...]
  }
  ```
  But in actual serialization, the "branches" array could appear in a different position (e.g., right after "bankName"), especially if your JSON library sorts keys alphabetically. The data remains the same; only the key order may differ. Forcing the order would require custom data structures, which in my opinion is unnecessary complexity for this task.

## Happy Testing and Exploring the Swift-Codes-API! 🚀
# Author: Kamil Piechota

