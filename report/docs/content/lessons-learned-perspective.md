# Lessons learned

## Evolution and refactoring

Over the time that we spent developing this project, we saw the potential in using tools to assist our work. We added static code analysis, benchmarking and error logging. All of these three things had a heavy impact on the evolution of the application. We benchmarked the application both internally and externally to figure out weak points with regards to efficiency and find potential bottlenecks. We also applied logging to detect any mistakes that we could have missed. To ensure that our codebase was bug-free, consistent and understandable we also used the SonarQube static analysis tool.

All of the elements mentioned here played a big role in the evolution and refactoring of the project. The changes we implemented reflects the results of these tools and had an overall positive impact on how we matured the application.

## Operation

The application was deployed and active while the simulator was pushing out messages. It was an interesting experience to observe the application and how it responded to various requests. With the monitoring tools, we had available we were able to spot issues that occurred and pinpoint which endpoints were affected. It was also helpful to see the data provided by the teachers benchmarking the performance.

Over time the application went from a single one node service to a more sophisticated manager-worker architecture that utilized load balancing and replication. By using these features we were able to scale our services horizontally to meet increasing traffic.  

### Monitoring

We discovered the importance of instrumenting production systems with metrics and visualizing them through graphs. Without having some sort of monitoring solution, the system is essentially a black box. As an example, after setting up Prometheus and Grafana, we discovered that one of our endpoints took on average about 10 seconds to respond to requests which were unacceptable. Through Grafana, we were able to pinpoint the responsible handler function, and we discovered that the slow execution time was due to an inefficient database query. This was not an issue during the first weeks of the system but became a problem as the size of the database tables increased. We were able to fix the issue by creating an index on a field of one of our database tables.

## Maintenance

Through the use of the Docker Swarm orchestrator, we were able to establish a system in which we could do maintenance to our containers without turning off the application. Because of our usage of 2 container replicas, we can at any time take one down and push new code to it while the other is serving clients. 

While pushing application updates was a pretty seamless experience, the database migration was a bit more clunky as we could not find a way to keep the service online while we swapped out the database engine.

## DevOps style of work

Compared to previous development projects at ITU, our workflow has shifted significantly in this course. The workflow of DevOps has been much more structured and automated. The repository has functioned as the central component of our work by both hosting our codebase and serving as a communication hub through issues. Branching was used extensively to separate the main development codebase from new feature implementations, thus allowing the development branch to always be in a functioning state. 













