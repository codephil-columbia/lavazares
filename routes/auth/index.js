const express = require('express');
const Pool = require('pg').Pool;
const bcrypt = require('bcrypt');

const router = express.Router();

const pool = new Pool({
  user: process.env.USER,
  host: process.env.HOST,
  database: process.env.DB,
  password: process.env.PW,
  port: process.env.DBPORT
})

pool.on('error', (err, client) => {
  console.log(err);
  process.exit(-1);
})

const _recievedMetadata = (req, res, next) => {
  const { username, password } = req.body;
  if (!username || !password) {
    res.end('username or password was null');
  } else {
    req.username = username;
    req.password = password;
    next();
  }
}

const _hashPassword = (req, res, next) => {
  bcrypt.hash(req.password, 10).then(hash => {
    req.password = hash;
    next();
  });
}

const _validatePassword = (storedHash, incomingPassword) => {
  return bcrypt.compare(incomingPassword, storedHash).then(res => {
    return res;
  })
}

router.use(_recievedMetadata);

router.post('/signup', [_hashPassword], (req, res) => {
  pool.connect()
    .then(client => {
      return client.query('INSERT INTO users(username, password) VALUES($1, $2)', [req.username, req.password])
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
})

router.post('/login', (req, res) => {
  pool.connect()
    .then(client => {
      return client.query('SELECT * FROM users WHERE username = $1 LIMIT 1', [req.username])
        .then(query => {
          const user = query.rows[0];
          _validatePassword(user.password, req.password).then(valid => {
            if(valid) 
              res.end('matched');
            else 
              res.end('didnt match');
          });
          client.release();
        })
        .catch(err => {
          client.release();
          console.log(err);
        })
    })
    .catch(err => {
      console.log(err);
      res.end('err');
    })
})


module.exports = router;