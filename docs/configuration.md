## Configuration

### Enter the URL for your Nagios instance.

1. In Mattermost, go to **System Console > Plugins > Nagios**.
2. Set the **Nagios URL**. Remember to add `http://` or `https://` at the beginning (e.g.: `https://nagios.fedoraproject.org`).
3. Select **Save**.
4. Go to**Plugins > Nagios**, then select **Enable**.
5. The plugin is now ready to use.

### Configuring the configuration files watcher

*This step is optional, although highly recommended.*

Regenerate the token for the configuration files watcher:
1. In Mattermost, go to **System Console > Plugins > Nagios**.
2. Select **Regenerate** to regenerate the token.
3. Copy the token - you're going to use it later.
4. Select **Save**.
5. Switch to the machine where Nagios is running:
  - Download the latest stable version of the watcher from the [releases page](https://github.com/mattermost/mattermost-plugin-nagios/releases)
  - Move the watcher: `chmod +x watcher1.1.0.linux-amd64 && sudo mv watcher1.1.0.linux-amd64 /usr/local/bin/watcher`.

#### Running the watcher as a systemd service

##### Preparing the systemd service unit file

Adjust `dir` (default if not set: `/usr/local/nagios/etc/`), `url`, and `token` flags to your setup.

```shell script
sudo bash -c 'cat << EOF > /etc/systemd/system/mattermost-plugin-nagios-watcher.service
[Unit]
Description=Nagios configuration files monitoring service
After=network.target

[Service]
Restart=on-failure
ExecStart=/usr/local/bin/watcher -dir /nagios/configuration/files/directory -url https://mattermost.server.address/plugins/nagios -token TheTokenFromStep1

[Install]
WantedBy=multi-user.target
EOF'
```

##### Starting the watcher

```shell script
systemctl daemon-reload
systemctl enable mattermost-plugin-nagios-watcher.service
systemctl start  mattermost-plugin-nagios-watcher.service
```
