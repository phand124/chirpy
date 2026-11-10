# chirpy
Server is set to use port 8080
Root path tied to index.html
Logo added under assets
Readiness endpoint  created, files are served under 
\n Metric counter endoint /metrics and /reset. \n /healthz and /metrics are GET restricted, /reset is POST restricted
\n Admin namespace contains metrics and reset, API for healthz\n
Chirp validation at /api/validate_chirp, chirps must be less than of equal to 140 characters long \n
Censorship for foul language in place \n
[main c3843af] Moving to Postgre for database, created users table with goose for migrations\n
 1 file changed, 10 insertions(+)
 create mode 100644 sql/schema/001_users.sql
Reset endpoint is locked to the 'dev' platform, deletes all Users from the server \n
