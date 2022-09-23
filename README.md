# whosdriving-be

... some explanations ...

```bash
docker build -t whosdriving-be:latest .
docker run -it --rm -p 9000:9000 -v /Users/carl/Projects/data:/app/data --name whosdriving-app whosdriving-be
```

## Mutations
. findOrCreate
```graphql
mutation FindOrCreateUser($newUser: NewUser!) {
  findOrCreateUser(input:$newUser) {
    email,
    lastName,
    firstName,
    profile,
    role
  }
}
```

variables
```grapql
{
  "newUser" : {
    "email": "john.smith@sample.com",
    "firstName": "John",
    "lastName": "Smith"
  }
}
```

## Query
Query
```grapql
query user($email: String!) {
 user(email:$email){
  email,
  firstName,
  lastName,
  profile,
  role
} 
}
```

variables
```grapql
{
  "email" : "john.smith@sample.com"
}
```
