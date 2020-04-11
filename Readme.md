A) Build commands:

auctioneer -> 
1. cd auctioneer
2. go build
3. docker build -t auctioneer:1.0.0 .

bidder ->
1. cd bidder
2. go build
3. docker build -t bidder:1.0.0 .


B) Bring up auction-system: docker-compose -f docker-compose.yml up -d

C) Brind down system:  docker-compose -f docker-compose.yml down