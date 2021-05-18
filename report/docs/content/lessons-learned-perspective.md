# Lessons learned

## Evolution and refactoring

## Operation

### Monitoring

We discovered the importance of instrumenting production systems with metrics and visualizing them through graphs. Without having some sort of monitoring solution, the system is essentially a black box. As an example, after setting up Prometheus and Grafana, we discovered that one of our endpoints took on average about 10 seconds to respond to requests which was unacceptable. Through Grafana, we were able to pinpoint the responsible handler function, and we discovered that the slow execution time was due to an inefficient database query. We were able to fix the issue by creating an index one of our database tables.

## Maintenance

