--
-- Import data
--

-- insert data from fs into json payloads
\set content `cat ../output/minerId-peerId.json`
insert into filecoin_miner_id_peer_id_scraper.payloads ("head", "json") values (1908194, :'content');

-- insert data from jsonb into relations table
insert into filecoin_miner_id_peer_id_api.relations ("head", "miner_id", "peer_id", "multiaddrs")
select p.head, j.key, j.value ->> 'peerId', array((select jsonb_array_elements_text(j.value -> 'multiaddrs')))
	from filecoin_miner_id_peer_id_scraper.payloads p, jsonb_each(p."json") j
	where p.processed is null or p.processed = false order by id asc;