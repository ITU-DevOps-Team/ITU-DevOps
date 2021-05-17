# Process' perspective

##A description and illustration of:


  - How do you interact as developers?
  - How is the team organized?
  - A complete description of stages and tools included in the CI/CD chains.
    -  That is, including deployment and release of your systems.
  - Organization of your repositor(ies).
    - That is, either the structure of of mono-repository or organization of artifacts across repositories.
    - In essence, it has to be be clear what is stored where and why.
  - Applied branching strategy.

A feature based branching strategy has been used for this project, that is anytime a team member desires to add a new feature, a new branch with the name of that feature is created. Once the changes has been implemented a pull-request is made and reviewed by another team member before it is merged with the development branch, which functions as the main branch. 

  - Applied development process and tools supporting it
    - For example, how did you use issues, Kanban boards, etc. to organize open tasks

Distribution of tasks among the team members has been done through GitHub's issues. When new tasks came up they were added as issues on GitHub and assigned a team memnber. 

  - How do you monitor your systems and what precisely do you monitor?
  - What do you log in your systems and how do you aggregate logs?

The goal for the project was to implement and utilize the ELK stack for analysing and aggregating logs on Kibana, but there were tons of trouble making the stack work with Docker Swarm. For the final release of this project the ELK stack has not been fully implemented, hence there are no logs collected and available through Kibana nor Elasticsearch. Contrarily the internal logging library of Golang is being used to some extend. ****HTTP responses and errors are being logged internally, but not collected by ELK stack****

  - Brief results of the security assessment.
  - Applied strategy for scaling and load balancing.

In order to secure the minitwit application for large amounts of users and operations and ensure a high level of availability the system has been set up using Docker in Swarm mode. The system is operating with a single swarm manger connected to two worker nodes forming a cluster. With Docker Swarm you can simply add more replicas of already running containers and let the manager node handle the distribution of the containers across the swarm. In case of failure within one of the worker nodes, Docker is capable of detecting this failure and spinning up new containers on the failing node. Docker Swarm also comes with internal load balancing, which is used for this project. That is, the manager node is capable of routing incoming requests to the worker nodes in order to maintain the best performance possible. 
