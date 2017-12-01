const express = require('express');
const Pool = require('pg').Pool;
const bcrypt = require('bcrypt');

const router = express.Router();

const pool = new Pool({
  user: process.env.DBUSER,
  host: process.env.HOST,
  database: process.env.DB,
  password: process.env.PW,
  port: process.env.DBPORT
})

pool.on('error', (err, client) => {
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
          res.status(200).end();
        })
        .catch(err => {
          client.release();
          res.status(500).send({err: "Internal server error."});
        })
    })
    .catch(err => {
      res.status(500).send({err: "Internal server error."});
    })
})

router.post('/login', (req, res) => {
  pool.connect()
    .then(client => {
      return client.query('SELECT * FROM users WHERE username = $1 LIMIT 1', [req.username])
        .then(query => {
          const user = query.rows[0];
          _validatePassword(user.password, req.password).then(valid => {
            if(valid) {
              res.status(200).send(user);
            } else {
              res.status(400).send({err: "Username or password is incorrect."});
            }
          });
          client.release();
        })
        .catch(err => {
          client.release();
          res.status(400).send({err: "Username does not exist."});
        })
    })
    .catch(err => {
      res.status(500).send({err: "Internal server error"});
    })
})

module.exports = router;