--
-- Create filecoin_miner_id_peer_id_scraper schema
--

CREATE SCHEMA IF NOT EXISTS filecoin_miner_id_peer_id_scraper AUTHORIZATION filecoin_miner_id_peer_id;

--DROP TABLE IF EXISTS filecoin_miner_id_peer_id_scraper.payloads;
CREATE TABLE IF NOT EXISTS filecoin_miner_id_peer_id_scraper.payloads (
	"id" SERIAL PRIMARY KEY,
	"head" INTEGER NOT NULL,
	"json" JSONB DEFAULT NULL,
	"processed" BOOLEAN DEFAULT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS payloads_id_idx ON filecoin_miner_id_peer_id_scraper.payloads ("id");
CREATE INDEX IF NOT EXISTS payloads_head_idx ON filecoin_miner_id_peer_id_scraper.payloads ("head");
CREATE INDEX IF NOT EXISTS payloads_processed_idx ON filecoin_miner_id_peer_id_scraper.payloads ("processed");
