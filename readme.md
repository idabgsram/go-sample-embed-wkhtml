# Sample Embed wkhtmlpdf

This is a sample of wkhtmltopdf implementation in go, with air autoreload integrated and dockerized wkhtmltopdf

# Run
1. Install docker
2. Execute the docker compose file, `docker compose -p sample-embed-wkhtml up -d --build -V`
3. Send empty POST request to localhost:4000/build-pdf
4. Check output folder
5. To stop, execute `docker compose -p sample-embed-wkhtml down -v`