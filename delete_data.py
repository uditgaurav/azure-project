import mysql.connector
from mysql.connector import errorcode

# Obtain connection string information from the portal

config = {
  'host':'udit-mysql.mysql.database.azure.com',
  'user':'udit@udit-mysql',
  'password':'mayadata@1658',
  'database':'litmus',
  'client_flags': [mysql.connector.ClientFlag.SSL],
  'ssl_ca': '/var/wwww/html/DigiCertGlobalRootG2.crt.pem'
}

# Construct connection string

try:
   conn = mysql.connector.connect(**config)
   print("Connection established")
except mysql.connector.Error as err:
  if err.errno == errorcode.ER_ACCESS_DENIED_ERROR:
    print("Something is wrong with the user name or password")
  elif err.errno == errorcode.ER_BAD_DB_ERROR:
    print("Database does not exist")
  else:
    print(err)
else:
  cursor = conn.cursor()

  # Drop previous table of same name if one exists
  cursor.execute("DROP TABLE IF EXISTS inventory;")
  print("Finished dropping table (if existed).")
  
  # Delete a data row in the table
  cursor.execute("DELETE FROM inventory WHERE name=%(param1)s;", {'param1':"orange"})
  print("Deleted",cursor.rowcount,"row(s) of data.")
  # Cleanup
  conn.commit()
  cursor.close()
  conn.close()
  print("Done.")
