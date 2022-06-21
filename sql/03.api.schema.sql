--
-- Create filecoin_miner_id_peer_id_api schema
--

CREATE SCHEMA IF NOT EXISTS filecoin_miner_id_peer_id_api AUTHORIZATION filecoin_miner_id_peer_id;

--DROP TABLE IF EXISTS filecoin_miner_id_peer_id_api.relations;
CREATE TABLE IF NOT EXISTS filecoin_miner_id_peer_id_api.relations (
	"id" SERIAL PRIMARY KEY,
	"head" INTEGER NOT NULL,
	"miner_id" VARCHAR(25) NOT NULL UNIQUE,
	"peer_id" VARCHAR(128) NOT NULL,
	"multiaddrs" VARCHAR(255)[] DEFAULT NULL,
	"timestamp" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS relations_id_idx ON filecoin_miner_id_peer_id_api.relations ("id");
CREATE INDEX IF NOT EXISTS relations_miner_id_idx ON filecoin_miner_id_peer_id_api.relations ("miner_id");
CREATE INDEX IF NOT EXISTS relations_peer_id_idx ON filecoin_miner_id_peer_id_api.relations ("peer_id");
