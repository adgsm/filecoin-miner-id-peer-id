package api

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/adgsm/filecoin-miner-id-peer-id/internal"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/jackc/pgx/v4/pgxpool"
)

// declare global vars
type Api struct {
	Router *mux.Router
}

var config internal.Config
var rcerr error
var confsPath = "configs/configs"
var db *pgxpool.Pool

func New(dtb *pgxpool.Pool) http.Handler {
	// read configs
	config, rcerr = internal.ReadConfigs(confsPath)
	if rcerr != nil {
		panic(rcerr)
	}

	// set db pointer
	db = dtb

	// set api struct
	a := &Api{
		Router: mux.NewRouter(),
	}
	a.Router.Host(config["api_host"])

	// set api v1 subroute
	v1 := a.Router.PathPrefix("/minerid-peerid/api/v1").Subrouter()

	// inti routes
	initRoutes(v1)

	// allow cros-origine requests
	cr := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://localhost:3000", fmt.Sprintf("https://%s", config["api_host"])},
		AllowCredentials: true,
	})
	hndl := cr.Handler(v1)

	return hndl
}

func initRoutes(r *mux.Router) {
	// search peer ids on provided miner ids
	r.HandleFunc("/peer-id", searchPeerIds).Methods(http.MethodGet)
	r.HandleFunc("/peer-id?miner_id={miner_id}", searchPeerIds).Methods(http.MethodGet)

	// search miner ids on provided peer ids
	r.HandleFunc("/miner-id", searchMinerIds).Methods(http.MethodGet)
	r.HandleFunc("/miner-id?peer_id={peer_id}", searchMinerIds).Methods(http.MethodGet)
}

func searchPeerIds(w http.ResponseWriter, r *http.Request) {
	// declare types
	type Record struct {
		Head       int
		MinerId    string
		PeerId     string
		Multiaddrs []string
	}

	// set defalt response content type
	w.Header().Set("Content-Type", "application/json")

	// collect query parameters
	queryParams := r.URL.Query()

	// check for provided miner ids
	miners := queryParams.Get("miner_id")

	internal.WriteLog("info", fmt.Sprintf("Search peer ids for provided miners: '%s'.", miners), "api")

	var minersArr interface {
		sql.Scanner
		driver.Valuer
	}
	// split miners comma delimited string into sql array
	if len(miners) == 0 {
		minersArr = pq.Array([]sql.NullString{})
	} else {
		minerIds := strings.Split(internal.SqlNullableString(miners).String, ",")
		for i := range minerIds {
			minerIds[i] = strings.TrimSpace(minerIds[i])
		}
		minersArr = pq.Array(minerIds)
	}

	// search sites with provided parameters
	rows, rowsErr := db.Query(context.Background(), "select \"head\", \"miner_id\", \"peer_id\", \"multiaddrs\" from filecoin_miner_id_peer_id_api.relations where miner_id = any($1);", minersArr)

	if rowsErr != nil {
		fmt.Print(rowsErr.Error())
		message := "Error occured whilst searching for peer ids."
		jsonMessage := fmt.Sprintf("{\"message\":\"%s\"}", message)
		internal.WriteLog("error", message, "api")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(jsonMessage))
		return
	}

	defer rows.Close()

	records := []Record{}

	for rows.Next() {
		var record Record
		if recordErr := rows.Scan(&record.Head, &record.MinerId,
			&record.PeerId, &record.Multiaddrs); recordErr != nil {
			message := fmt.Sprintf("Error occured whilsrt reading peer ids from the database. (%s)", recordErr.Error())
			jsonMessage := fmt.Sprintf("{\"message\":\"%s\"}", message)
			internal.WriteLog("error", message, "api")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(jsonMessage))
			return
		}
		records = append(records, record)
	}

	// send response
	sitesJson, errJson := json.Marshal(records)
	if errJson != nil {
		message := "Cannot marshal the database response for searched miner ids."
		jsonMessage := fmt.Sprintf("{\"message\":\"%s\"}", message)
		internal.WriteLog("error", message, "api")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(jsonMessage))
		return
	}

	// response writter
	w.WriteHeader(http.StatusOK)
	w.Write(sitesJson)
}

func searchMinerIds(w http.ResponseWriter, r *http.Request) {
	// declare types
	type Record struct {
		Head       int
		MinerId    string
		PeerId     string
		Multiaddrs []string
	}

	// set defalt response content type
	w.Header().Set("Content-Type", "application/json")

	// collect query parameters
	queryParams := r.URL.Query()

	// check for provided peer ids
	peers := queryParams.Get("peer_id")

	internal.WriteLog("info", fmt.Sprintf("Search miner ids for provided peers: '%s'.", peers), "api")

	var peersArr interface {
		sql.Scanner
		driver.Valuer
	}
	// split peers comma delimited string into sql array
	if len(peers) == 0 {
		peersArr = pq.Array([]sql.NullString{})
	} else {
		peerIds := strings.Split(internal.SqlNullableString(peers).String, ",")
		for i := range peerIds {
			peerIds[i] = strings.TrimSpace(peerIds[i])
		}
		peersArr = pq.Array(peerIds)
	}

	// search sites with provided parameters
	rows, rowsErr := db.Query(context.Background(), "select \"head\", \"miner_id\", \"peer_id\", \"multiaddrs\" from filecoin_miner_id_peer_id_api.relations where peer_id = any($1);", peersArr)

	if rowsErr != nil {
		fmt.Print(rowsErr.Error())
		message := "Error occured whilst searching for miner ids."
		jsonMessage := fmt.Sprintf("{\"message\":\"%s\"}", message)
		internal.WriteLog("error", message, "api")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(jsonMessage))
		return
	}

	defer rows.Close()

	records := []Record{}

	for rows.Next() {
		var record Record
		if recordErr := rows.Scan(&record.Head, &record.MinerId,
			&record.PeerId, &record.Multiaddrs); recordErr != nil {
			message := fmt.Sprintf("Error occured whilsrt reading miner ids from the database. (%s)", recordErr.Error())
			jsonMessage := fmt.Sprintf("{\"message\":\"%s\"}", message)
			internal.WriteLog("error", message, "api")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(jsonMessage))
			return
		}
		records = append(records, record)
	}

	// send response
	sitesJson, errJson := json.Marshal(records)
	if errJson != nil {
		message := "Cannot marshal the database response for searched peer ids."
		jsonMessage := fmt.Sprintf("{\"message\":\"%s\"}", message)
		internal.WriteLog("error", message, "api")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(jsonMessage))
		return
	}

	// response writter
	w.WriteHeader(http.StatusOK)
	w.Write(sitesJson)
}
