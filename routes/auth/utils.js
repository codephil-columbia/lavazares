const uidGen = require('uuid/v4');

const generateNewUid = () => {
    return uidGen();
}
  
const validatePassword = (storedHash, incomingPassword) => {
  return bcrypt.compare(incomingPassword, storedHash).then(res => {
    return res;
  })
}

module.exports = {
    generateNewUid,
    validatePassword
}