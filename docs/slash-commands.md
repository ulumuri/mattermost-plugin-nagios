### Slash commands overview

- `nagios`
    - `get-logs <alerts|notifications>`
        - `[host|service <host name|service description>]`
    - `set-logs-limit <count>`
    - `set-logs-start-time <seconds>`
    - `subscribe <report|configuration-changes>`
    - `unsubscribe <report|configuration-changes>`
    - `set-report-frequency <minutes>`

### Slash commands documentation

#### nagios

`nagios`

This is the root command.

##### get-logs

`get-logs <alerts|notifications> [<host|service> <host name|service description>]`

This action allows you to get alerts or notifications.

Example: `/nagios get-logs alerts`

###### host

`get-logs <alerts|notifications> host <host name>`

This optional parameter allows you to get alerts or notifications from a specific host.

Example: `/nagios get-logs alerts host bvmhost-p09-02.iad2.fedoraproject.org`

###### service

`get-logs <alerts|notifications> service <service description>`

This optional parameter allows you to get alerts or notifications from a specific service.

Example: `/nagios get-logs alerts service Swap-Is-Low`

##### set-logs-limit

`set-logs-limit <count>`

This action allows you to limit the number of logs `get-logs` fetches.

Example: `/nagios set-logs-limit 10`

##### set-logs-start-time

`set-logs-start-time <seconds>`

This action allows you to specify the age of the oldest log `get-logs` fetches.

Example: `/nagios set-logs-start-time 3600`

##### subscribe

`subscribe <report|configuration-changes>`

This action allows you to subscribe to system monitoring reports or configuration changes on the current channel.

Example: `/nagios subscribe report`

##### unsubscribe

`unsubscribe <report|configuration-changes>`

This action allows you to unsubscribe from system monitoring reports or configuration changes on the current channel.

Example: `/nagios unsubscribe configuration-changes`

##### set-report-frequency

`set-report-frequency <minutes>`

This action allows you to set the frequency of system monitoring reports.

Example: `/nagios set-report-frequency 60`
