const express = require('express');
const Pool = require('pg').Pool;
const bcrypt = require('bcrypt');
const utils = require('./utils');

const router = express.Router();

const pool = new Pool({
  user: process.env.DBUSER,
  host: process.env.HOST,
  database: process.env.DB,
  password: process.env.PW,
  port: process.env.DBPORT
})

pool.on('error', (err, client) => {
  console.log(err);
  process.exit(-1);
})

const hashPassword = (req, res, next) => {
  bcrypt.hash(req.password, 10)
    .then(hash => {
      req.password = hash;
      next();
    }).catch(err => {
      console.log(err);
    })
}

router.post('/signup', (req, res) => {
  const {username, password, email} = req.body; 
  const uid = utils.generateNewUid();
  pool.query(
    'INSERT INTO users(username, password, email, uid) VALUES($1, $2, $3, $4)', 
    [username, password, email, uid]
  ).then(query => {  
      res.status(200).json({uid});
  }).catch(err => {
      console.log(err);
      res.status(500).send({err: "Error making new user", err});
  })
})

router.post('/login', (req, res) => {
  pool.connect()
    .then(client => {
      return client.query('SELECT * FROM users WHERE username = $1 LIMIT 1', [req.username])
        .then(query => {
          const user = query.rows[0];
          utils.validatePassword(user.password, req.password).then(valid => {
            if(valid) {
              console.log(req.session.id);
              res.status(200).send(user.username);
            } else {
              res.status(400).send({err: "Username or password is incorrect."});
            }
          });
          client.release();
        })
        .catch(err => {
          client.release();
          res.status(400).send({err: "Username does not exist. " + err});
        })
    })
    .catch(err => {
      res.status(500).send({err: "Internal server error " + err});
    })
})

module.exports = router;