const express = require('express');
const Pool = require('pg').Pool;

const router = express.Router();

const pool = new Pool({
  user: 'cesaribarra',
  host: 'localhost',
  database: 'typephil',
  password: 'cesaribarra',
  port: 5432
})

const insertQuery = (tableName, tableFields, values) => {
  
}

router.post('/', (req, res) => {
  const { username, password } = req.body;
  if(username && password) {
    pool.connect()
      .then(client => {
        return client.query('INSERT INTO users(username, password) VALUES($1, $2)', [username, password])
          .then(res => {
            client.release();
          })
          .catch(err => {
            client.release();
            console.log(err.stack);
          })
      })
      .catch(err => {
        console.log(err);
      })
    res.end('good');
  } else {
    res.end('bad');
  }
})

module.exports = router;