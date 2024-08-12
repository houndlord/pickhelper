package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestNewDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPing()

	testDB, err := NewDB("sqlmock_db")
	assert.NoError(t, err)
	assert.NotNil(t, testDB)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSavePatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testDB := &DB{db}

	mock.ExpectExec("INSERT INTO patches").WithArgs("13.10").WillReturnResult(sqlmock.NewResult(1, 1))

	err = testDB.SavePatch(PatchInfo{Version: "13.10"})
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateScrapingStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testDB := &DB{db}

	mock.ExpectExec("INSERT INTO scraping_status").WithArgs("13.10", "13.9", true).WillReturnResult(sqlmock.NewResult(1, 1))

	err = testDB.UpdateScrapingStatus(ScrapingStatus{CurrentPatch: "13.10", LastScrapedPatch: "13.9", IsUpdating: true})
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCurrentPatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testDB := &DB{db}

	rows := sqlmock.NewRows([]string{"version"}).AddRow("13.10")
	mock.ExpectQuery("SELECT version FROM patches").WillReturnRows(rows)

	patch, err := testDB.GetCurrentPatch()
	assert.NoError(t, err)
	assert.Equal(t, "13.10", patch.Version)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveChampion(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testDB := &DB{db}

	mock.ExpectExec("INSERT INTO champions").WithArgs("Ahri", "http://example.com/ahri.png").WillReturnResult(sqlmock.NewResult(1, 1))

	err = testDB.SaveChampion(Champion{Name: "Ahri", AvatarURL: "http://example.com/ahri.png"})
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveMatchups(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testDB := &DB{db}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO matchups").WithArgs("Ahri", "Zed", "mid", 48.5, 1000, "13.10").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	matchups := []Matchup{
		{Champion: "Zed", WinRate: "48.5", SampleSize: "1000"},
	}

	err = testDB.SaveMatchups("Ahri", "mid", matchups, "13.10")
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTopMatchups(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testDB := &DB{db}

	rows := sqlmock.NewRows([]string{"name", "win_rate", "sample_size"}).
		AddRow("Zed", 48.5, 1000).
		AddRow("Yasuo", 51.2, 800)

	mock.ExpectQuery("SELECT c.name, m.win_rate, m.sample_size FROM matchups").WithArgs("Ahri", "mid", "13.10", 2).WillReturnRows(rows)

	matchups, err := testDB.GetTopMatchups("Ahri", "mid", 2, "13.10")
	assert.NoError(t, err)
	assert.Len(t, matchups, 2)
	assert.Equal(t, "Zed", matchups[0].Champion)
	assert.Equal(t, "48.50", matchups[0].WinRate)
	assert.Equal(t, "1000", matchups[0].SampleSize)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllChampions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testDB := &DB{db}

	rows := sqlmock.NewRows([]string{"name", "avatar_url"}).
		AddRow("Ahri", "http://example.com/ahri.png").
		AddRow("Zed", "http://example.com/zed.png")

	mock.ExpectQuery("SELECT name, avatar_url FROM champions").WillReturnRows(rows)

	champions, err := testDB.GetAllChampions()
	assert.NoError(t, err)
	assert.Len(t, champions, 2)
	assert.Equal(t, "Ahri", champions[0].Name)
	assert.Equal(t, "http://example.com/ahri.png", champions[0].AvatarURL)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMatchupsEndpoint(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testDB := &DB{db}

	r := gin.Default()
	r.GET("/matchups/:champion/:role", func(c *gin.Context) {
		champion := c.Param("champion")
		role := c.Param("role")

		status := ScrapingStatus{CurrentPatch: "13.10", LastScrapedPatch: "13.10", IsUpdating: false}
		rows := sqlmock.NewRows([]string{"name", "win_rate", "sample_size"}).
			AddRow("Zed", 48.5, 1000).
			AddRow("Yasuo", 51.2, 800)

		mock.ExpectQuery("SELECT current_patch, last_scraped_patch, is_updating FROM scraping_status").WillReturnRows(sqlmock.NewRows([]string{"current_patch", "last_scraped_patch", "is_updating"}).AddRow(status.CurrentPatch, status.LastScrapedPatch, status.IsUpdating))
		mock.ExpectQuery("SELECT c.name, m.win_rate, m.sample_size FROM matchups").WithArgs(champion, role, status.LastScrapedPatch, 8).WillReturnRows(rows)

		matchups, err := testDB.GetTopMatchups(champion, role, 8, status.LastScrapedPatch)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"patch": status.LastScrapedPatch, "matchups": matchups})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/matchups/Ahri/mid", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "13.10")
	assert.Contains(t, w.Body.String(), "Zed")
	assert.Contains(t, w.Body.String(), "Yasuo")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChampionsEndpoint(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testDB := &DB{db}

	r := gin.Default()
	r.GET("/champions", func(c *gin.Context) {
		rows := sqlmock.NewRows([]string{"name", "avatar_url"}).
			AddRow("Ahri", "http://example.com/ahri.png").
			AddRow("Zed", "http://example.com/zed.png")

		mock.ExpectQuery("SELECT name, avatar_url FROM champions").WillReturnRows(rows)

		champions, err := testDB.GetAllChampions()
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if len(champions) == 0 {
			c.JSON(404, gin.H{"error": "No champions found"})
			return
		}

		c.JSON(200, gin.H{"champions": champions})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/champions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Ahri")
	assert.Contains(t, w.Body.String(), "Zed")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
