package utils

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

var AppPort = GetEnv("APP_PORT", "8080")
var LogfileName = GetEnv("APP_LOGFILE", "testServ.log")
var RequestLogs = GetEnv("REQ_LOGFILE", "requestlogs.log")
var RespLogs = GetEnv("RESP_LOGFILE", "resplogs.log")
var PGHostname = GetEnv("POSTGRES_HOSTNAME", "db")
var PGPort = GetEnv("POSTGRES_PORT", "5432")
var PGUsername = GetEnv("POSTGRES_USER", "docker")
var PGPassword = GetEnv("POSTGRES_PASSWORD", "docker")
var PGDB = GetEnv("POSTGRES_DB", "transaction_queue")
var CancelLog = GetEnv("CANCEL_LOGFILE", "cancellogs.log")
var CancelRows = GetEnv("CANCEL_ROWS", "10")
var Duration = GetEnv("DB_UPDATE_TIME", "10")

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func FormConnString(pgHostname string, pgPort string,
	pgUsername string, pgPassword string, pgDB string) string {
	return fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		pgHostname, pgPort, pgUsername, pgPassword, pgDB)
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandString(length int) string {
	return StringWithCharset(length, charset)
}

func CloseDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println("Unable to close DB")
	}
}

func CloseFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Println("Unable to close file")
	}
}

func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Println("Unable to close DB Rows")
	}
}

func RespLogger() gin.HandlerFunc {
	resplogs, err := os.Create(RespLogs)
	if err != nil {
		fmt.Println("Unable to create log file")
	}
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		statusCode := c.Writer.Status()
		timestamp := time.Now().UnixNano()
		_, err := resplogs.WriteString(
			fmt.Sprintf("Timestamp: %v Status Code: %v,  Response body: %v", timestamp, statusCode, blw.body.String()),
		)
		if err != nil {
			log.Println("Unable to log HTTP Responce")
		}

	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ReadBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		log.Println("Unable to log request info")
	}
	s := buf.String()
	return s
}

//Log every request header/body/IP/timestamp
func RequestLogger() gin.HandlerFunc {
	requestlogs, err := os.Create(RequestLogs)
	if err != nil {
		fmt.Println("Unable to create log file")
	}
	return func(c *gin.Context) {
		host := c.Request.Host

		buf, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Abort()
		}
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
		timestamp := time.Now().UnixNano()
		source := c.Request.Header.Get("Source-Type")
		_, err = requestlogs.WriteString(fmt.Sprintf("Time: %v, IP: %v, Source-Type: %v, Body: %v \n",
			timestamp, host, source, ReadBody(rdr1)))
		if err != nil {
			log.Println("Unable to log HTTP Request")
		}
		c.Request.Body = rdr2
		c.Next()
	}
}

func LogCancels(s string) {
	f, err := os.OpenFile(CancelLog, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer CloseFile(f)
	log.SetOutput(f)
	log.Println(s)
}
