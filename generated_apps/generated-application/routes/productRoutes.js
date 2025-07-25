const express = require('express');
const router = express.Router();
const productController = require('../controllers/productController');

// GET /api/products - Get all products
router.get('/', productController.getAll);

// GET /api/products/:id - Get product by ID
router.get('/:id', productController.getById);

// POST /api/products - Create new product
router.post('/', productController.create);

// PUT /api/products/:id - Update product
router.put('/:id', productController.update);

// DELETE /api/products/:id - Delete product
router.delete('/:id', productController.delete);

module.exports = router;