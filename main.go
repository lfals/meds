package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	testConnection()
	createTable()
	insertMeds()
}

func testConnection() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connection successful, test query result:", result)
	err = db.Close()

	if err != nil {
		log.Fatal(err)
	}
}

func createTable() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS meds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		TIPO_PRODUTO TEXT,
		NOME_PRODUTO TEXT,
		DATA_FINALIZACAO_PROCESSO TEXT,
		CATEGORIA_REGULATORIA TEXT,
		NUMERO_REGISTRO_PRODUTO TEXT,
		DATA_VENCIMENTO_REGISTRO TEXT,
		NUMERO_PROCESSO TEXT,
		CLASSE_TERAPEUTICA TEXT,
		EMPRESA_DETENTORA_REGISTRO TEXT,
		SITUACAO_REGISTRO TEXT,
		PRINCIPIO_ATIVO TEXT
	);`)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func insertMeds() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("./assets/DADOS_ABERTOS_MEDICAMENTOS(in).csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	// Skip header row
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	count, err := db.Query("SELECT COUNT(id) FROM meds")

	if err != nil {
		log.Fatal(err)
	}

	var rowCount int
	defer count.Close()
	for count.Next() {
		if err := count.Scan(&rowCount); err != nil {
			log.Fatal(err)
		}
	}
	if err := count.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Current row count: %d", rowCount)

	if condition := rowCount > 0; condition {
		log.Println("Data already inserted, skipping...")
		db.Close()
		return

	}

	stmt, err := db.Prepare(`INSERT INTO meds (
		TIPO_PRODUTO,
		NOME_PRODUTO,
		DATA_FINALIZACAO_PROCESSO,
		CATEGORIA_REGULATORIA,
		NUMERO_REGISTRO_PRODUTO,
		DATA_VENCIMENTO_REGISTRO,
		NUMERO_PROCESSO,
		CLASSE_TERAPEUTICA,
		EMPRESA_DETENTORA_REGISTRO,
		SITUACAO_REGISTRO,
		PRINCIPIO_ATIVO
	) VALUES (?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(
			record[0],  // TIPO_PRODUTO
			record[1],  // NOME_PRODUTO
			record[2],  // DATA_FINALIZACAO_PROCESSO
			record[3],  // CATEGORIA_REGULATORIA
			record[4],  // NUMERO_REGISTRO_PRODUTO
			record[5],  // DATA_VENCIMENTO_REGISTRO
			record[6],  // NUMERO_PROCESSO
			record[7],  // CLASSE_TERAPEUTICA
			record[8],  // EMPRESA_DETENTORA_REGISTRO
			record[9],  // SITUACAO_REGISTRO
			record[10], // PRINCIPIO_ATIVO
		)
		log.Println(record)
		if err != nil {
			log.Fatal(err)
		}
	}

	db.Close()
}

func main() {

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Missing search query parameter", http.StatusBadRequest)
			return
		}

		db, err := sql.Open("sqlite3", "database.db")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		rows, err := db.Query(`
			SELECT distinct NOME_PRODUTO, PRINCIPIO_ATIVO FROM meds 
			WHERE SITUACAO_REGISTRO = 'VÁLIDO' 
			and NOME_PRODUTO LIKE ? 
			OR PRINCIPIO_ATIVO LIKE ? 
			ORDER BY (
			CASE WHEN NOME_PRODUTO = ? THEN 1 
			WHEN NOME_PRODUTO LIKE ? THEN 2 
			ELSE 3 
			END),
			NOME_PRODUTO LIMIT 10`,
			"%"+query+"%", "%"+query+"%", "%"+query+"%", query, query+"%")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		columns, _ := rows.Columns()

		for rows.Next() {
			values := make([]interface{}, len(columns))
			pointers := make([]interface{}, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}

			if err := rows.Scan(pointers...); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			entry := make(map[string]interface{})
			for i, column := range columns {
				val := values[i]
				entry[column] = val
			}
			results = append(results, entry)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	fmt.Println("Server starting on port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
