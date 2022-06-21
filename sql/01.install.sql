-- login as postgres
CREATE ROLE filecoin_miner_id_peer_id WITH LOGIN PASSWORD '****secret****';
CREATE DATABASE filecoin_miner_id_peer_id;
\c filecoin_miner_id_peer_id
CREATE SCHEMA IF NOT EXISTS filecoin_miner_id_peer_id_scraper AUTHORIZATION filecoin_miner_id_peer_id;
CREATE SCHEMA IF NOT EXISTS filecoin_miner_id_peer_id_api AUTHORIZATION filecoin_miner_id_peer_id;
