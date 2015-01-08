package main

import (
  // _ imports this package anonymously so we don't get an error
  // from not explicitly using it
  _ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"
  "log"
  "time"
)

type Todo struct {
  Id          int
  Subject     string
  Description string
  Completed   bool
  CreatedAt   time.Time `db:"created_at"`
  UpdatedAt   time.Time `db:"updated_at"`
}

func main() {
  db, err := sqlx.Open("postgres", "postgres://postgres:postgres@localhost/godos_development?sslmode=disable")
  if err != nil {
    log.Fatal(err)
  }

  // Start db transaction
  tx := db.MustBegin()
  now := time.Now()
  t := Todo{
    Subject:      "Mow Lawn",
    Description:  "Yuck!",
    CreatedAt:    now,
    UpdatedAt:    now,
  }
  tx.Exec("insert into todos (subject, description, created_at, updated_at) values $1, $2, $3, $4", t.Subject, t.Description, t.CreatedAt, t.UpdatedAt)
  // This is an invalid record that will rollback the transaction
  tx.Exec("insert into todos (subject, description, created_at, updated_at) values $1, $2, $3, $4", nil, t.Description, t.CreatedAt, t.UpdatedAt)
  // End db transaction
  tx.Commit()

  todos := []Todo{}
  err = db.Select(&todos, "select * from todos")
  if err != nil {
    log.Fatal(err)
  }

  for _, todo := range todos {
    log.Printf("Subject is %s", todo.Subject)
  }
}
