const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const morgan = require('morgan');
const db = require('./config/database');

// Import routes
const productRoutes = require('./routes/productRoutes');


const app = express();
const PORT = process.env.PORT || 8080;

// Middleware
app.use(helmet());
app.use(cors());
app.use(morgan('combined'));
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Routes
app.get('/', (req, res) => {
  res.json({
    message: 'Welcome to Generated Application API',
    version: '1.0.0',
    endpoints: [
      'GET /api/products',
      'GET /api/products/{id}',
      'POST /api/products',
      'PUT /api/products/{id}',
      'DELETE /api/products/{id}',
    ]
  });
});

app.use('/api/products', productRoutes);


// Error handling middleware
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({
    error: 'Something went wrong!',
    message: err.message
  });
});

// 404 handler
app.use('*', (req, res) => {
  res.status(404).json({
    error: 'Route not found'
  });
});

// Initialize database connection
db.connect().then(() => {
  console.log('Database connected successfully');
  
  app.listen(PORT, '0.0.0.0', () => {
    console.log('Server is running on port ' + PORT);
    console.log('API Documentation: http://localhost:' + PORT);
  });
}).catch(err => {
  console.error('Database connection failed:', err);
  process.exit(1);
});