# configmap-generator
Creates kubernetes configmaps.

This is a CLI application.

It reads ansible group_vars and converts them to configmaps. 
It can decrypt ansible vault files.

Usage examples:
```
go run cmd/cmapgen/main.go -H  generate -n all -e vmp \
    -g ~/utvikling/vimond-ansible/vmp/group_vars/ \
    -p ~/utvikling/vimond-ansible/vmp/secrets/vault_password.txt  | kubectl apply -n dev   -f -

#test
go run cmd/cmapgen/main.go -H  generate -n all -e myenv3 \
    -c testdata/test-app-config.yml \
    -g ./testdata/ansible1/vmp/group_vars \
    -p ./testdata/ansible1/vmp/test-secret.txt  


```