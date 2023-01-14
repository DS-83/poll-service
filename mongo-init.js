db = db.getSiblingDB("admin")

db.createUser(
    {
        user: "admin",
        pwd: "admin", //passwordPrompt(), // or cleartext password
      roles: [
          { role: "userAdminAnyDatabase", db: "admin" },
          { role: "readWriteAnyDatabase", db: "admin" }
        ]
    }
    )

const database = 'polls';
const collection = 'polls';
    
db = db.getSiblingDB(database)
    // db.dropUser("dev")
db.createUser({
    user: 'dev',
    pwd: 'SuperSecretPassword',
    roles: [
        {
            role: 'dbAdmin',
            db: 'polls',
        },
    ],
});

db.createCollection(
    collection, {
    validator: {
        $jsonSchema: {}
    }
})


db.polls.insertOne(
    {
        "question": "Do?",
        "choice": [
          {
            "_id": {
              "$oid": "63c29534257f0d1b01f2a094"
            },
            "name": "you",
            "votes": 1
          },
          {
            "_id": {
              "$oid": "63c29534257f0d1b01f2a095"
            },
            "name": "like",
            "votes": 0
          },
          {
            "_id": {
              "$oid": "63c29534257f0d1b01f2a096"
            },
            "name": "it",
            "votes": 0
          }
        ]
      }
)

db.runCommand({
    collMod: "polls",
     validator: {
        $jsonSchema: {
           bsonType: "object",
           title: "Poll Object Validation",
           required: [ "choice", "question" ],
           properties: {
              question: {
                 bsonType: "string",
                 description: "'question' must be a string and is required"
              },
              choice: {
                 bsonType: "array",
                 description: "'choice' must be an array and is required",
                 required: ["name"],
                 properties: {
                  name: {
                    bsonType: "string",
                    description: "'name' must be a string and is required"
                  }
                 }
              },
           }
        }
     }
  } )
  

db.adminCommand( { shutdown: 1 } )
