# Change Log
All notable changes to this project will be documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).
 
## [v0.0.1] - 2024-09-03
 
Here is my first working version of the application.

Results:
    - **/contagem-pessoas** -> 5257
    - **Gatling output**: rinhabackendsimulation-20240903120253660
 
## [v0.0.2] - 2024-09-03
 
Results:
    - **/contagem-pessoas** -> 5612
    - **Gatling output**: rinhabackendsimulation-20240903175250292
    - commit: 61030342dee8e5a89431a85d548b658c0d8d603b
 
## [v0.0.3] - 2024-09-03
 
Results:
    - **/contagem-pessoas** -> 7418
    - **Gatling output**: rinhabackendsimulation-20240903183250538
    - commit: 67991054a2294c1f95c8588086aab7dfc29d1d60
  
### Added

### Changed
- More cpu to the database, less to nginx and apps

### Fixed

## [v0.0.4] - 2024-09-03
 
Results:
    - **/contagem-pessoas** -> 10610
    - **Gatling output**: rinhabackendsimulation-20240903194230195
    - commit: 9a308d346d5e8ebcc97ba569966d40c82c83cc0f
  
### Added

### Changed
- Added Go context to the queries in the database

### Fixed

## [v0.0.5] - 2024-09-04
 
Results:
    - **/contagem-pessoas** -> 31798
    - **Gatling output**: rinhabackendsimulation-20240904152625151
    - commit: 8055928b1ad7ce07f1aa13af759a4da149cea3c3
  
### Added
- Connection pool in the Postgre Driver
- Created a search column with auto generate data in Postgres
  - This increased the performance by from 8k/10k to 29/30k
- Created an index for search column
  - This was about 10% gain in inserted users.

### Changed
- Improved Nginx configuration
  - Removed access logging

### Fixed
- Fixed the UNIQUE logic to be on the database side, and not to query the database first then insert the new person.

## [v0.0.6] - 2024-09-04
 
Results:
    - **/contagem-pessoas** -> 24804
    - **Gatling output**: rinhabackendsimulation-20240904172439247
    - commit: 1985849ff84e3fe073a6cc4c41f5bd65dbe5bd12
  
### Added

### Changed

### Fixed
- There was a validation bug adding more results then it should in the backend. I discovered that removing CPU and Memory limitations

## [v0.0.7] - 2024-09-04
 
Results:
    - **/contagem-pessoas** -> 36874
    - **Gatling output**: rinhabackendsimulation-20240904175808154
    - commit: 0224100a244624597c4430beb0480cf1bf242f7f
  
### Added

### Changed

### Fixed
- I forgot to remove the ILIKE in the search query instead of using just LIKE given I had changed the index to lower case.
- I also did some clean up and other settings that didn't affected much the results.


## [v0.0.8] - 2024-09-04
 
Results:
    - **/contagem-pessoas** -> 39957
    - **Gatling output**: rinhabackendsimulation-20240904203628807
    - commit: 81129dd3689806030d61de3913aa9dded460bf65
  
### Added

### Changed
- I've changed the Postgres cache size that will help with the CPU bound time.

### Fixed

## [v0.0.9] - 2024-09-04
 
Results:
    - **/contagem-pessoas** -> 46500
    - **Gatling output**: rinhabackendsimulation-20240905114121892
    - commit: 81129dd3689806030d61de3913aa9dded460bf65

Restrictions:
    - Using Linux
    - Removing the CPU/Memory limits because I haven't implemented the GiSP feature in Postgres.
  
### Added

### Changed
- network_mode=host that will remove the IO Exception.
- new nginx configuration to also help with IO Exception.

### Fixed

## [v0.0.10] - 2024-09-05
 
Results:
    - **/contagem-pessoas** -> 46240
    - **Gatling output**: rinhabackendsimulation-20240905214422519
    - commit: 7f638c96313dfb2d8128eca3730cdc346eb2d225
  
### Added

### Changed
- changed from `pg` to `pgx and pgxpool`, this was the bottleneck that was somehow making Postgres have more latency or use more resources to respond.

### Fixed

## [v0.1.0] - 2024-09-06
 
Results:
    - **/contagem-pessoas** -> 45718
    - **Gatling output**: rinhabackendsimulation-20240906193400539
    - commit:754d56e49f373d2bcf4d59bdbd606990b7836ec0
  
### Added

### Changed
- use `[fiber](https://docs.gofiber.io)` instead of `net/http` library for REST API. This seems to give a slightly faster response time in miliseconds.

### Fixed

## [v0.2.0] - 2024-09-06
 
Results:
    - **/contagem-pessoas** -> 45919
    - **Gatling output**: rinhabackendsimulation-20240907150913707
    - commit: 84fda5fd4f7de161b31fd97aef8fad0053ef1c13

It's hard to measure real results given the network mode bridge, but this seems a superior solution for high load on production.
  
### Added
- adding `rueidis` for Redis driver for caching the nickname and create person.
- worker to insert in bulk in a async way, ensure consistency using the cache.
- logger to all layers
- new e2e tests

### Changed
- ajust docker compose memory and cpu distribution.

### Fixed