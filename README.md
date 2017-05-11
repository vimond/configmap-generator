# configmap-generator
Creates kubernetes configmaps.

This is a CLI application.

It reads ansible group_vars and converts them to configmaps. 
It can decrypt ansible vault files.


Usage examples:
```
go run cmd/cmapgen/main.go -H  generate -n all -e env \
    -g ~/src/ansible-proj/group_vars/ \
    -p ~/.ansible/vault_password.txt  | kubectl apply -n dev   -f -

#test
go run cmd/cmapgen/main.go -H  -c testdata/test-app-config.yml \
    generate -n all -e myenv3 \
    -g ./testdata/ansible1/vmp/group_vars \
    -p ./testdata/ansible1/vmp/test-secret.txt  
```

Might be some rough edges still, but should mainly work.

Plans are:
* be able to create several types of config items
    * docker `pull secrets`
    * plain `secrets`
* automatically split out variables derived from secrets into `secrets`
