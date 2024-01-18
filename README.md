A repository for a simple Golang project for storing username and passwords, hashing the passwords after passing and fetch all the details.

Please ensure to create a local SQL database and necessary tables to implement all functionalities, the table structure:
![image](https://github.com/jaysri-eng/GoWeb/assets/72025056/7e964722-2add-4127-9b8b-c69f5cd20327)

Also, make sure to create a config.yaml file to put all your database configurations like username, password etc... this is a cleaner way to access details of the database. For reference, 
app:
  server:
    port: 80

database:
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: go
  sslmode: disable
  timezone: Asia/Kolkata
