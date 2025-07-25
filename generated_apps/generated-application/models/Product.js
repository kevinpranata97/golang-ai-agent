class Product {
  constructor(data = {}) {
    this.id = data.id || 0;
    this.name = data.name || '';
    this.description = data.description || '';
    this.price = data.price || 0;
    this.created_at = data.created_at || new Date();
  }

  // Validation method
  validate() {
    const errors = [];

    if (!this.id) {
      errors.push('id is required');
    }

    if (!this.name) {
      errors.push('name is required');
    }
    // Add validation for name: min=1,max=200


    if (!this.price) {
      errors.push('price is required');
    }
    // Add validation for price: min=0

    if (!this.created_at) {
      errors.push('created_at is required');
    }

    return errors;
  }

  // Convert to JSON
  toJSON() {
    return {
      id: this.id,
      name: this.name,
      description: this.description,
      price: this.price,
      created_at: this.created_at,
    };
  }

  // Create from database row
  static fromRow(row) {
    return new Product(row);
  }
}

module.exports = Product;