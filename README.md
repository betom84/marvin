# Install marvin on raspberry pi

- Create Folder `/home/pi/marvin` on raspberry pi
- Run `deploy.sh [TARGET]` to deploy binary and resources
- Copy marvin.sh to `/etc/init.d/marvin`
- Run `sudo update-rc.d marvin defaults` to install rc-scripts for system startup
