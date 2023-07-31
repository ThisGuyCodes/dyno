# Running locally
Assuming you have Docker installed, you can run `./run.sh` on mac or linux. Also feel free to `docker build .` and run yourself, default port inside the container is 8080.

If you don't have or you don't want to use docker, `go build` should also work (assuming a compatible version of Go, I tested 1.20), then just run `./dyno`, again default port is 8080.

While running in Docker works *great*, because of my insistence on running under privieleged containers, persistence is left unsolved, the database will be new on each run.

# Uploading
Do a multipart-form file upload to `/upload` as a POST, the file should be named `sequences` in the upload. Here's an example with curl:
`curl -F sequences=@sequence_file.txt.gz -X POST localhost:8080/upload`

# Querying
Do a POST or GET to `/hamming_matches` with a form or URL parameter `sequence` set to the sequence you want to query for. Results will include the descriptions, internal UUIDs, and full contents of any matching sequences (in JSON format). An example with curl:
`curl "localhost:8080/hamming_matches?sequence=${SEQUENCE}"`

This will return a 404 if nothing matches (and also an empty JSON array).
