

# Smart Falcon

## Overview

**Smart Falcon** is a blockchain-based project leveraging **Hyperledger Fabric** to effectively manage and track assets for financial institutions. This system provides a secure, transparent, and immutable solution for asset management. The project encompasses the development of chaincode for peer interactions and a REST API to facilitate smart contract invocations.

## Table of Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [Setup and Execution](#setup-and-execution)
- [Technologies Used](#technologies-used)
- [Testing](#testing)
- [License](#license)

## Project Structure

- **`chaincode/`**: Contains Go code for Hyperledger Fabric chaincode.
- **`api/`**: Contains the REST API code for blockchain interaction.
- **`docker/`**: Includes Docker configuration files.
- **Additional files**: Configuration and utility scripts from the Hyperledger open-source repository.

## Setup and Execution

### 1. Environment Setup

This project utilizes **WSL** (Windows Subsystem for Linux) and **Docker** for containerization. Ensure both are properly installed and configured on your system.

### 2. Cloning the Repository

Initialize the project by cloning the Hyperledger Fabric samples repository:

```bash
git clone https://github.com/hyperledger/fabric-samples.git
```

This provides essential resources for setting up a test network.

### 3. Chaincode Development and Deployment

- **Develop Chaincode**: Implement asset management operations using Go.
- **Package and Install Chaincode**: Package the chaincode and install it on the peer nodes within the Hyperledger Fabric network.
- **Configure Chaincode Path**: Set the appropriate chaincode path and verify communication between peers to ensure seamless interaction.
- **Query Chaincode**: Execute a query on the chaincode to retrieve the Package ID, confirming the correct path configuration.

### 4. REST API Development

1. **Develop REST API**: Create a REST API using Go to interact with the deployed chaincode.
2. **Implement API Endpoints**: Design and implement endpoints for creating, querying, and updating assets on the blockchain.

## Technologies Used

- **Hyperledger Fabric**: Blockchain framework.
- **Go**: Programming language for chaincode and REST API.
- **Docker**: Containerization platform.
- **WSL** (Windows Subsystem for Linux): Development environment.
- **Postman**: API testing tool.

## Testing

- **REST API**: Thoroughly tested using Postman to ensure proper communication with the Hyperledger Fabric network.
- **Chaincode**: Validated using peer commands in the terminal for operations such as asset creation and querying.

## License

This project is based on open-source resources provided by Hyperledger Fabric and adheres to its licensing terms. For detailed licensing information, please refer to the Hyperledger Fabric documentation.

---
