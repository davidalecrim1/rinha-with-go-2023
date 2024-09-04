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
