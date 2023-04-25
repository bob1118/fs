
# fs tool

my freeswitch command line tool.

## howto

### freeswitch

install freeswitch first.
<https://freeswitch.org/confluence/display/FREESWITCH/Debian>

### postgres

install postgresql and set password.

```shell
apt-get install -y postgresql

su postgres
cd
postgres=# psql
postgres=# \password
postgres=# Enter new password for user "postgres":
postgres=# Enter it again:
postgres=# \quit
exit
```

### fs

install fs, init conf, set var=value , run gateway, run server.

```shell
#install fs.
go install github.com/bob1118/fs
# set postgres password
fs config --set postgres.password='yourpassword'
# set switch ipv4 address.
fs config --set switch.vars.ipv4='yourswitchipv4'
# set switch public ip(external_sip_ip and external_rtp_ip)
# for lan deploy
fs config --set switch.vars.external_sip_ip='$${local_ip_v4}'
fs config --set switch.vars.external_rtp_ip='$${local_ip_v4}'
# for internet deploy
fs config --set switch.vars.external_sip_ip='stun:stun.freeswitch.org'
fs config --set switch.vars.external_rtp_ip='stun:stun.freeswitch.org'
# init switch bootable conf.
fs config fsconfig --reset
fs config fsconfig --init
# run gateway
fs gateway --run
# run server
fs server --run
# restart freeswitch
systemctl restart freeswitch
```

## feature
