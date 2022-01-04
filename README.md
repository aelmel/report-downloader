#Report Downloader

Report downloader is a cli app that on configured cron will run endpoint
to generate a report and then will download to local temp directory

### Configuration
Can configure app using args or env variables
* API_URL            - report address to call report generation and downloading (default: http://localhost:12345)
* API_CLIENTS        - number of clients to call report generation tool (default: 2)
* GENERATE_FREQUENCY - cron expression to run generate api (default: */1 * * * *)

