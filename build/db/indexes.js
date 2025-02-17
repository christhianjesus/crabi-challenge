db = db.getSiblingDB('default'); 
db.createCollection("user");
db.user.createIndex({ "email": 1 }, { unique: true });