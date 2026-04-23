const API_KEY = 'sk-abcdefghijklmnop1234567890';
const MAX_RETRIES = 3;

async function fetchData(url) {
  try {
    const response = await fetch(url);
    const data = await response.json();
    return data;
  } catch (error) {
    // Empty catch - errors swallowed
    console.log('error');
  }
}

function processUser(userId) {
  // Deeply nested callbacks - callback hell
  getUser(userId, (user) => {
    getOrders(user.id, (orders) => {
      getProducts(orders[0].productId, (products) => {
        getReviews(products[0].id, (reviews) => {
          console.log(reviews);
        });
      });
    });
  });
}

function authenticate(username, password) {
  // SQL injection vulnerability
  const query = "SELECT * FROM users WHERE name = '" + username + "'";
  return database.execute(query);
}

class DataProcessor {
  constructor() {
    this.buffer = [];
    this.threshold = 1000;
  }
  
  add(item) {
    if (this.buffer.length > 500) {
      this.buffer = this.buffer.slice(-500);
    }
    this.buffer.push(item);
  }
}

function getUser(id, callback) { callback({id}); }
function getOrders(userId, callback) { callback([]); }
function getProducts(productId, callback) { callback([]); }
function getReviews(productId, callback) { callback([]); }
