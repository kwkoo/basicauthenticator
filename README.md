# HTTP Basic Authentication Identity Provider

This project demonstrates how you can write an identity provider to authenticate users logging into OpenShift.

The user database (userdb.tsv) is a text file with a user entry on each line. Each line has 4 fields separated by tabs - userid, password, name, email.

Follow the instructions [here](https://docs.okd.io/latest/install_config/configuring_authentication.html#BasicAuthPasswordIdentityProvider) to configure OpenShift for this identity provider.

Here is a sample of the changes made to `master-config.yaml`:

````
...
  identityProviders:
  - challenge: true
    login: true
    mappingMethod: claim
    name: custom_basic_auth_provider
    provider:
      apiVersion: v1
      kind: BasicAuthPasswordIdentityProvider
      url: http://192.168.1.239:8080
  masterCA: ca-bundle.crt
...
````

Note that if you are using `oc cluster up` to start OpenShift, it writes `master-config.yaml` in 3 places - `kube-apiserver`, `openshift-apiserver`, `openshift-controller-manager`. You will need to modify `master-config.yaml` in all 3 directories.
