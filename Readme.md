
# fs

freeswitch command line tool.

## dev init

```shell
#init golang fs
cd 
mkdir fs && cd fs
go mod init github.com/bob1118/fs

#init cobra-cli
go install github.com/spf13/cobra-cli@latest
cobra-cli init --author "bob" --viper

cobra-cli add config
cobra-cli add fsconfig -p configCmd
cobra-cli add gateway
cobra-cli add server

git init
git add --all

```

# next

next
