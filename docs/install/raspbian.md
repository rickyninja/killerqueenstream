# install raspbian lite operating system

* [raspbian lite](https://www.raspberrypi.org/downloads/raspbian/)

**Be careful not to overwite your system hard drive by accident!**


Run this command to see the device name of the SD card while you insert it.
`dmesg -w`
Press `CTRL+C` to interrupt the dmesg command after you identify the device name.

Install os image to micro sd card
`dd if=2017-11-29-raspbian-stretch-lite.img of=/dev/sd? bs=4096 conv=fsync`

# raspbian post install

Once the pi boots, you probably won't know it's IP address yet, and sshd may not be running yet.
I typically console the pi by connecting it to my TV, and plugging in a spare keyboard.

## squash the keyboard rebellion
It seems like the raspbian image comes with a default locale of Great Britain, so I immediately need to
change the default locale and keyboard config, or I'm unable to use pipe in the shell.

```
sudo dpkg-reconfigure locales
sudo dpkg-reconfigure keyboard-configuration
```

## enable remote access
I prefer to work remotely rather than consoled on my TV.
```
sudo systemctl enable ssh
sudo systemctl start ssh.service
sudo reboot
```

After the reboot is complete, you can find it's IP address via `ifconfig` or `ip address show` on newer Debian versions.  From this point you should be able to connect via ssh if desired.

## personal prefs

I prefer vim to nano.

```
sudo apt purge nano
sudo apt install vim-nox
```

## login setup

Since sshd will be running, definitely set up your own account and disable the default login that came with the OS install.

### add your own login
`sudo useradd jeremys`

### give yourself sudo
`sudo visudo`
Add a line like this: `jeremys ALL=(ALL:ALL) NOPASSWD:ALL`

### disable default login
`passwd -l pi`
