package integrTest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/handler"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func decodeBody[T any](body io.Reader) (T, error) {
	var res T

	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&res)
	if err != nil {
		return res, errors.New("bad response body")
	}

	return res, nil
}

func TestPostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	dbName := "postgres"
	dbUser := "postgres"
	dbPassword := "12345"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:14",
		postgres.WithInitScripts(filepath.Join("migrations", "init.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	defer postgresContainer.Terminate(ctx)
	if err != nil {
		t.Errorf("failed to start container: %s", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Errorf("get conn string error: %s", err)
	}

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		t.Errorf("sql open error: %s", err)
	}

	if err = db.Ping(); err != nil {
		t.Errorf("ping error: %s", err)
	}

	server := handler.NewServer(db)
	r := http.NewServeMux()
	h := api.HandlerFromMux(server, r)

	srv := httptest.NewServer(h)

	resp, err := srv.Client().Post(srv.URL+"/dummyLogin", "application/json", bytes.NewBufferString(`{"role": "moderator"}`))
	if err != nil {
		t.Errorf("dummy login error: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("bad dummy login response status: %v", resp.StatusCode)
	}

	moderatorToken, err := decodeBody[api.Token](resp.Body)
	if err != nil {
		t.Errorf("decode token response error: %s", err)
	}

	req, err := http.NewRequest("POST", srv.URL+"/pvz", bytes.NewBufferString(`{"city": "Санкт-Петербург", "id": "53aa35c8-e659-44b2-882f-f6056e443c99"}`))
	if err != nil {
		t.Errorf("http create request error: %s", err)
	}
	req.Header.Set("Authorization", moderatorToken)

	resp, err = srv.Client().Do(req)
	if err != nil {
		t.Errorf("post pvz error: %s", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("bad post pvz response status: %v", resp.StatusCode)
	}

	pvz, err := decodeBody[api.PVZ](resp.Body)
	if err != nil {
		t.Errorf("decode pvz response error: %s", err)
	}

	resp, err = srv.Client().Post(srv.URL+"/dummyLogin", "application/json", bytes.NewBufferString(`{"role": "employee"}`))
	if err != nil {
		t.Errorf("dummy login error: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("bad dummy login response status: %v", resp.StatusCode)
	}

	employeeToken, err := decodeBody[api.Token](resp.Body)
	if err != nil {
		t.Errorf("decode token response error: %s", err)
	}

	req, err = http.NewRequest("POST", srv.URL+"/receptions", bytes.NewBufferString(fmt.Sprintf(`{"pvzId": "%s"}`, pvz.Id)))
	if err != nil {
		t.Errorf("http create request error: %s", err)
	}
	req.Header.Set("Authorization", employeeToken)

	resp, err = srv.Client().Do(req)
	if err != nil {
		t.Errorf("post receptions error: %s", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("bad receptions response status: %v", resp.StatusCode)
	}

	_, err = decodeBody[api.Reception](resp.Body)
	if err != nil {
		t.Errorf("decode reception response error: %s", err)
	}

	for i := 0; i < 50; i++ {
		req, err = http.NewRequest("POST", srv.URL+"/products", bytes.NewBufferString(fmt.Sprintf(`{"pvzId": "%s", "type": "одежда"}`, pvz.Id)))
		if err != nil {
			t.Errorf("http create request error: %s", err)
		}
		req.Header.Set("Authorization", employeeToken)

		resp, err = srv.Client().Do(req)
		if err != nil {
			t.Errorf("post products error: %s", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("bad products response status: %v", resp.StatusCode)
		}

		_, err = decodeBody[api.Product](resp.Body)
		if err != nil {
			t.Errorf("decode product response error: %s", err)
		}
	}
}
