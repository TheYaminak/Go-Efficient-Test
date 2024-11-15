# Real Estate Data Processing Project

This project demonstrates the performance comparison between sequential and concurrent data processing using Go routines. The dataset used is related to real estate sales in the United States and can be downloaded from the following link: [Real Estate Sales Data (2001-2018)](https://catalog.data.gov/dataset/real-estate-sales-2001-2018).

## Overview

The project is structured to show the benefits of using Go routines for concurrent processing versus traditional sequential processing. It does this by reading a CSV file containing real estate sales data, validating the data, and then sending it to two different servers: one that processes the data sequentially and another that processes it concurrently. Finally, the results and average response times are compared.

## Project Structure

The project is divided into several packages to maintain modularity and clarity:

### 1. `models`
- Contains the data model (`RealEstate`) representing the real estate sales records.

### 2. `csvreader`
- Responsible for reading the CSV file and parsing it into Go structs (`RealEstate` records).
- **Function**: `ReadCSV(filename string) ([]models.RealEstate, error)` reads the CSV file and returns a list of real estate records.

### 3. `handlers`
- Contains the HTTP handlers that receive and process the data.
- **Functions**:
  - `SequentialHandler`: Processes incoming HTTP requests by handling records sequentially.
  - `ConcurrentHandler`: Processes incoming HTTP requests by handling records concurrently using Go routines.
  - `StartServers()`: Starts both the sequential and concurrent HTTP servers on ports `8080` and `8081` respectively.

### 4. `services`
- Contains the core logic for validating and processing real estate records.
- **Functions**:
  - `ValidateData(record models.RealEstate) bool`: Validates the real estate data to ensure it meets required criteria.
  - `AdjustValues(record *models.RealEstate)`: Applies adjustments to the data, such as market fluctuations.

### 5. `dispatcher`
- Contains the logic to send requests to the servers for comparison.
- **Function**: `SendRequestsToServers(records []models.RealEstate)` is responsible for sending the records to both sequential and concurrent servers and calculating the average response time.

### 6. `main`
- Entry point of the application.
- **Steps**:
  1. Starts the sequential and concurrent HTTP servers using `handlers.StartServers()`.
  2. Reads the real estate sales data using `csvreader.ReadCSV(filename)`, which loads data from a CSV file.
  3. Calls `dispatcher.SendRequestsToServers(records)` to send the records to the servers and compare the response times.

## How to Run

1. **Download the Dataset**:
   - Download the real estate dataset from [this link](https://catalog.data.gov/dataset/real-estate-sales-2001-2018) and save it as `Real_Estate_Sales_2001-2022_GL.csv` in the root directory.

2. **Build and Run the Project**:
   - Make sure Go is installed on your system.
   - Run the following commands:
     ```sh
     go mod tidy
     go run main.go
     ```

3. **Endpoints**:
   - The project starts two servers:
     - Sequential Server: `http://localhost:8081/process_sequential`
     - Concurrent Server: `http://localhost:8082/process_concurrent`

4. **Results**:
   - After processing the data, metrics are printed to the console showing the total time and average response time for both the sequential and concurrent processing approaches.

## Purpose

This project demonstrates:
- The performance improvement that can be achieved by using Go routines for concurrent data processing.
- How to design an application in Go that utilizes both sequential and concurrent paradigms.

## Folder Structure
```
real_estate_project/
├── models/
│   └── realestate.go          # Contains the RealEstate struct
├── csvreader/
│   └── reader.go              # Reads and parses CSV data
├── handlers/
│   └── handler.go             # Contains HTTP handlers for sequential and concurrent processing
│   └── goroutines.go
│   └── sequential.go
├── services/
│   └── service.go             # Core business logic for validating and processing records
├── dispatcher/
│   └── dispatcher.go          # Sends requests to both servers and collects metrics
├── main.go                    # Entry point for the project
└── Real_Estate_Sales_2001-2022_GL.csv  # CSV file containing real estate data
```

## Future Improvements
- **Error Handling**: Improve error handling to manage edge cases more robustly.
- **Logging**: Add structured logging for better insights and debugging.
- **Configuration**: Use configuration files for server ports, URLs, and other parameters.
- **Testing**: Add unit and integration tests to ensure reliability.

## License
This project is licensed under the MIT License, which means you are free to use, modify, and distribute the code, as long as proper attribution is given to this repository. Please include a link back to the [GitHub repository](https://github.com/TheYaminak/Go-Efficient-Test) when using this code.

