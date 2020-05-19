# NFTables Exporter

Prometheus exporter for NFTables byte and packet counters.

## Usage

This exporter only exports counters of rules with comments.
Ensure that every rule with a counter has a unique comment in its chain/table.

### Example

```
# HELP nftables_bytes_total Total number of bytes
# TYPE nftables_bytes_total counter
nftables_bytes_total{chain="input_syn",comment="ratelimit reached",table="filter"} 259
# HELP nftables_packets_total Total number of packets
# TYPE nftables_packets_total counter
nftables_packets_total{chain="input_syn",comment="ratelimit reached",table="filter"} 24966
```
