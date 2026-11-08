# chirpy
Server is set to use port 8080
Root path tied to index.html
Logo added under assets
Readiness endpoint  created, files are served under 
\n Metric counter endoint /metrics and /reset. \n /healthz and /metrics are GET restricted, /reset is POST restricted
\n Admin namespace contains metrics and reset, API for healthz\n
Chirp validation at /api/validate_chirp, chirps must be less than of equal to 140 characters long \n
