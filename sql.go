package main

import (
  // _ imports this package anonymously so we don't get an error
  // from not explicitly using it
  _ "github.com/lib/pq"
  "database/sql"
  "log"
  "time"
)

func main() {
  sql.Create("postgres", "postgres://postgres:postgres@localhost/godos_development?sslmode=disable")
  db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/godos_development?sslmode=disable")
  if err != nil {
    log.Fatal(err)
  }

  now := time.Now()
  res, err := db.Exec(
    "insert into todos (subject, description, created_at, updated_at) values ($1, $2, $3, $4)",
    "Mow the lawn", "", now, now)
  if err != nil {
    log.Fatal(err)
  }
  affected, _ := res.RowsAffected()
  log.Printf("Rows affected %d", affected)

  var subject string

  rows, err := db.Query("select subject from todos")
  for rows.Next() {
    if err := rows.Scan(&subject); err != nil {
      log.Fatal(err)
    }
    log.Printf("Subject is %s", subject)
  }

  if err := rows.Err(); err != nil {
    log.Fatal(err)
  }
}
