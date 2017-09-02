/*
Package ingestion is a server that handles the forwarding of requests from
Beeswax log messages and device beacons to CPM

FIXME: workout final hierarchy and relationship of packages. I sense that some
of the functions and dependencies can be better separated (or combined)
for example the beeswax/auctionwon handler should probably be in the beeswax package

TODO: proper logger interface

TODO: add further tests for beacons and beeswax
*/
package ingestion
