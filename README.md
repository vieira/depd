depd: a deployment daemon
=========================

depd is a deployment daemon written in Go that listens for github and
bitbucket HTTP POST requests, triggers a deployment recipe according
to a `.depd.json` file present in the root directory of the repository
and notifies the owners once the deployment is completed.

adaptors
--------------------

adaptors are used to translate an incoming json payload to a common
internal structure used by `depd` to reason about a push. adaptors for
github and bitbucket are included. others can easily be added.

deployers
--------------------

after pulling the changes from the remote repository a specific deployer
can be run if a `.depd.json` file is present in the root directory. a
wordpress deployer is included, as an example, capable of bringing wordpress
to a target version and invalidating the apc opcode cache.

notifiers
---------------------

after deployment is complete a log of the various actions executed can be sent
as a notification. a mail notifier is included but it is easy to write other
notifiers for hipchat, irc, etc.

batteries included
----------------------

also included are an example systemd unit file, a upstart script, a nginx
vhost reverse proxy config file for using `depd` behind nginx (recommended)
and an example configuration file `config.json`
