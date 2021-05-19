# Process' perspective

## A description and illustration of:


- How do you interact as developers?

We interacted with each other mainly through Teams voice and instant chat where we planned meetings and aligned tasks. After our meetings we would use the built-in Github features of the kanban board and issues section.

- How is the team organized?

Our team is organized as a flat hierarchy where we would debate on issues and taking joint decisions. There was no clear manager or leader during this process but due to the relatively small scope of the project this did not feel like an issue.

- A complete description of stages and tools included in the CI/CD chains.
  -  That is, including deployment and release of your systems.

The CICD pipelines are implemented with Github Actions. Initially, there was just a single workflow which was building the docker images, pushing them to the Docker Hub repository and then running a docker compose file on the host machine while also forcing the new images to be pulled.

As the project evolved, we ended up with 5 workflows:

- `deploy-minitwit` - consists of two jobs `build-minitwit` and `deploy-minitwit`. `build-minitwit` builds the docker images for the minitwit service then pushes them to docker hub. `deploy-minitwit` uses a marketplace action [appleboy/ssh-action](https://github.com/marketplace/actions/ssh-remote-commands) to ssh into the Swarm Manager node, pull the latest image and trigger a service update. The workflow gets triggered on push to the development (default) branch or can be triggered manually on any branch.
- `deploy-minitwit-api` - consists of two jobs `build-minitwit-api` and `deploy-minitwit-api`. The steps for these jobs are exactly like the ones from above but for the minitwit-api service. The workflow gets triggered on push to the development (default) branch or can be triggered manually on any branch.
- `ci formatting` - consists of two jobs `format-minitwit` and `format-minitwit-api`. Both of these jobs use a marketplace action (Jerome1337/gofmt-action@v1.0.4)[https://github.com/marketplace/actions/check-code-formatting-using-gofmt] to perform formatting (go fmt) on the root directories and fails on code not meeting the formatting standards. The workflow gets triggered on push and pull request to the development (default) branch.
- `format` - consists of `generate-report`, a job that generates a pdf report from the markdown files, uploads the pdf as an artefact and then commits and pushes the pdf report to 'report/build'. The workflow gets triggered on push to development (default) branch or branches matching the 'docs/*' pattern.
- `sonarcloud` - enabled by installing the `SonarCloud` Github App. Performs linting and checks for bugs, vulnerabilities, code smells and security hotspots. The cicd chain get triggered on pull requests to any branch.

- Organization of your repositor(ies).
  - That is, either the structure of of mono-repository or organization of artifacts across repositories.
  - In essence, it has to be be clear what is stored where and why.
- Applied branching strategy.

A feature-based branching strategy has been used for this project, that is anytime a team member desires to add a new feature, a new branch with the name of that feature is created. Once the changes have been implemented a pull request is made and reviewed by another team member before it is merged with the development branch, which functions as the main branch.

- Applied development process and tools supporting it
  - For example, how did you use issues, Kanban boards, etc. to organize open tasks

Distribution of tasks among the team members has been done through GitHub's issues and a Kanban board. When new tasks came up they were added as issues on GitHub and assigned a team member. These issues were tracked on the Kanban board using a standard layout containing sections for TODO, In Progress, Under Review and Done. By using both Github Projects and Github Issues we have seamless integration between issues and kanban tasks which gives us a better overview of the project status.

- How do you monitor your systems and what precisely do you monitor?

Monitoring is done with Prometheus, where various metrics is defined in the application. Prometheus scapes the application for these metrics once every 5 seconds. These metrics are pulled by Grafana which has a built-in customizable dashboard for visualizing these metrics. Specificaly we monitor the following targets:
 - Frontend application:
   - asdfasdf
 - `minitwit_ui_usertimeline_requests`- asdfadsf
 - `minitwit_ui_personaltimeline_requests`- asdfadsf
 - `minitwit_ui_unfollow_requests`- asdfadsf
 - `minitwit_ui_follow_requests`- asdfadsf
 - `minitwit_ui_addmessage_requests`- asdfadsf
 - `minitwit_ui_homepage_requests`- asdfadsf
 - `minitwit_ui_register_requests`- asdfadsf
 - `minitwit_ui_login_requests`- asdfadsf
 - `minitwit_ui_logout_requests`- asdfadsf
 - `minitwit_ui_total_requests`- asdfadsf
 - `minitwit_ui_register_requests`- asdfadsf
 - `minitwit_api_register_requests`- asdfadsf
 - `minitwit_api_messages_requests`- asdfadsf
 - `minitwit_api_messages_per_user_requests`- asdfadsf
 - `minitwit_api_follow_requests`- asdfadsf
 - `minitwit_api_total_requests`- asdfadsf
 - `minitwit_api_latest_execution_time_in_ns`- asdfadsf
 - `minitwit_api_register_execution_time_in_ns`- asdfadsf
 - `minitwit_api_messages_execution_time_in_ns`- asdfadsf
 - `minitwit_api_messages_per_user_execution_time_in_ns`- asdfadsf
 - `minitwit_api_follow_execution_time_in_ns`- asdfadsf
 - `minitwit_api_authentication_middleware_execution_time_in_ns`- asdfadsf
 - `minitwit_api_latest_middleware_execution_time_in_ns`- asdfadsf

- What do you log in your systems and how do you aggregate logs?

The goal for the project was to implement and utilize the ELK stack for analysing and aggregating logs on Kibana, but there were manu challenges making the stack work with Docker Swarm. For the final release of this project the ELK stack has not been fully implemented, hence there are no logs collected and available through Kibana nor Elasticsearch. Contrarily the internal logging library of Golang is being used to some extend. HTTP responses and errors are being logged locally, but not collected by ELK stack, due to our challenges with Docker Swarm. The ELK stack is implemented in the application using Filebeat to collect and ship logfiles to Elasticsearch without the L in ELK. That is, Logstash has not been included in the logging stack for this application. Kibana is supposed to fetch logging data from Elasticsearch, but when the application switched to Docker Swarm, Elasticsearch was not receiving any logs. Although the stack is not fully functional in the final application, it was working properly before transforming to a Docker Swarm cluster.

- Brief results of the security assessment.
- Applied strategy for scaling and load balancing.

In order to secure the minitwit application for large amounts of users and operations and ensure a high level of availability the system has been set up using Docker in Swarm mode. The system is operating with a single swarm manager connected to two worker nodes forming a cluster. With Docker Swarm you can simply add more replicas of already running containers and let the manager node handle the distribution of the containers across the swarm. In case of failure within one of the worker nodes, Docker is capable of detecting this failure and spinning up new containers on the failing node. Docker Swarm also comes with internal load balancing, which is used for this project. That is, the manager node is capable of routing incoming requests to the worker nodes in order to maintain the best performance possible. 

Another benefit of Docker Swarm is overlay network and service discovery features that are enabled when using it as an orchestration tool. All containers launched by the manager are added with its own unique DNS name such that we can access and investigate separate containers with ease.
