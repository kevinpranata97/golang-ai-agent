# Generated Application

Create a Node.js API for selling satay with endpoints for managing satay products, orders, and customers.

## Features

- product_management


## Prerequisites

- Node.js 18+ 
- npm or yarn
- sqlite database

## Installation

1. Clone the repository
2. Install dependencies:
   `bash
   npm install
   `

3. Copy environment configuration:
   `bash
   cp .env.example .env
   `

4. Update the `.env` file with your configuration

5. Set up your sqlite database

6. Start the application:
   `bash
   npm run dev
   `

## API Endpoints

- `GET /api/products` - Get all products
- `GET /api/products/{id}` - Get product by ID
- `POST /api/products` - Create new product
- `PUT /api/products/{id}` - Update product
- `DELETE /api/products/{id}` - Delete product


## Project Structure

`
Generated Application/
├── app.js              # Main application file
├── package.json        # Dependencies and scripts
├── .env.example        # Environment configuration template
├── Dockerfile          # Docker configuration
├── controllers/        # Request handlers
├── models/            # Data models
├── routes/            # API routes
├── middleware/        # Custom middleware
└── config/            # Configuration files
`

## Development

- `npm run dev` - Start development server with auto-reload
- `npm start` - Start production server
- `npm test` - Run tests

## Docker

Build and run with Docker:

`bash
docker build -t Generated Application .
docker run -p 8080:8080 Generated Application
`

## License

MIT