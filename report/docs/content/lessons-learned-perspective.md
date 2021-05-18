# Lessons learned

## Evolution and refactoring



## Operation

### Monitoring

We discovered the importance of instrumenting production systems with metrics and visualizing them through graphs. Without having some sort of monitoring solution, the system is essentially a black box. As an example, after setting up Prometheus and Grafana, we discovered that one of our endpoints took on average about 10 seconds to respond to requests which was unacceptable. Through Grafana, we were able to pinpoint the responsible handler function, and we discovered that the slow execution time was due to an inefficient database query. We were able to fix the issue by creating an index one of our database tables.

## Maintenance



## DevOps style of work

Compared to previous development projects at ITU, our workflow has shifted significantly in this course. The workflow of DevOps has been much more structured and automated. The repository has functioned as the central component of our work by both hosting our codebase and serving as a communication hub through issues. Branching was used extensively to separate the main development code base from new feature implementations, thus allowing the development branch to always be in a functioning state. 













