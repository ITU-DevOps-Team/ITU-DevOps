This is the file containing our decision and argumentations on those during our minitwit project

## REWRITING MINITWIT IN GO

The decision to go with Go (pun intended :)) was based on the following reasons:

- Minimalistic language with an extensive standard library
- Fast, lightweight and scalable
- Strongly and statically typed (an advantage over python)
- Concurrency is an integral part supported by goroutines and channels
- Programs are constructed from packages that offer clear code separation
- Google APIs are often written in Go, it would be nice to follow those implementations

Choosing GORM as our database abstraction layer was based on the following reasons:

- Full-featured ORM
- Many different association types
- Aims to be developer friendly
- One of the most popular ORMs for Go

Going with GitHub Actions as our CI/CD platform was based on the following reasons:
- Easy integration with GitHub repository
- Easy setup compared to other CI/CD solutions
- Free to use up to 2000 minutes per month

Choosing PostgreSQL as database system because:

- Supports concurrency while adhering to the ACID principles
- One of the most popular DBMSs in the world. Makes it easy to find help online
- Free to use

Choosing Prometheus/Grafana was based on the following reasons:
- The combination of Prometheus and Grafana is an industry standard for monitoring Go applications. Additionally, it is very easy to connect Prometheus to Grafana.
- Prometheus delivers metrics without creating time lag on performance
- Prometheus is relatively easy to configure in Docker
- Grafana offers very customizable dashboards for visualizing application performance

Choosing Sonarqube for static code analysis was based on the following reasons:

- Integrates well with popular IDEs such as Eclipse and Visual Studio
- Is very feature-packed. Has support for identification of duplicated code, unit testing, code complexity, code smells and much more
- One of the most popular static analysis solutions on the market

Choosing ?? for system logging was based on the following reasons:

Choosing Docker Swarm for cluster management was based on the following reasons:

- Integrates well with our existing Docker setup. Uses the same files as used by Docker Compose locally
- Good for smaller setups, but also scales very well up to around 1000 nodes
- Easy to setup