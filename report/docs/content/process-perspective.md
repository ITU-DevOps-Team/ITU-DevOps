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
