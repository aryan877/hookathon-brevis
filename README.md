# SuperFluidSwapper Brevis Integration

This repository contains the Brevis Network integration and ZK circuit implementations for the SuperFluidSwapper project.

## Key Components

- ZK circuit implementation for fee adjustment proofs
- Off-chain data processing and proof generation
- Integration with Brevis Network SDK

## Prerequisites

- Go (v1.20 or later)
- Node.js (v14 or later)
- Yarn package manager

## Quick Start

1. Clone the repository:

   ```
   git clone https://github.com/aryan877/hookathon-brevis.git
   cd hookathon-brevis
   ```

2. Install Go dependencies:

   ```
   go mod tidy
   ```

3. Install Node.js dependencies:

   ```
   yarn install
   ```

4. Set up environment variables:
   Create a `.env` file in the root directory and add the following:
   ```
   BREVIS_API_KEY=your_brevis_api_key
   PROVER_ENDPOINT=your_prover_endpoint
   ```

## Project Structure

- `app/`: Contains the TypeScript application for interacting with the Brevis SDK
  - `src/index.ts`: Main script for fetching data, generating proofs, and updating fees
- `prover/`: Contains the Go implementation of the ZK circuit and prover
  - `circuits/`: ZK circuit implementation
    - `circuit.go`: Main circuit logic for fee adjustment
    - `circuit_test.go`: Tests for the circuit

## Running the Prover

To start the local prover service:

```
go run cmd/prover/main.go
```

## Generating Proofs

To generate a proof:

```
yarn generate-proof
```

## Testing

Run the Go test suite:

```
go test ./...
```

This will run tests including those in `prover/circuits/circuit_test.go`.

## Scripts

- `yarn build`: Build the TypeScript files
- `yarn lint`: Run ESLint
- `yarn format`: Format the code using Prettier

## Brevis SDK Integration

The project uses the Brevis SDK for ZK proof generation and verification. The main integration can be found in `app/src/index.ts`.

## Key Files

- `prover/circuits/circuit.go`: Implements the ZK circuit for fee adjustment calculations
- `app/src/index.ts`: Off-chain script for data fetching, proof generation, and fee updates

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENS
