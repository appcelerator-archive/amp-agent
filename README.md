###amp-agent

amp-agent is an infrastucture service installed on each node sending local node Docker events and local containers logs to kafka.

### version 1.0.0-0

Get Docker events related to container start, stop , kill, die, destroy and use it to maintain internally a list of running containers with their id, node_id, service_id, service_name
For each, open a log stream and then them to Kafka.

### api v1.0.0.0-0

api doc:


    * /api/v1/health: return code 200 or 400 concidering if amp-agent is ready, maining log-agent has open an event stream with docker and this stream is still active

    * /api/v1/containers: return a json containing the active containers list with service_id, service_name, status, and health status