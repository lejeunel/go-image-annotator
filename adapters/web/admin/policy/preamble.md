*Policies determine the correspondance between roles and method, i.e. which role can perform what actions.*

### Instructions

Access policies are configured using a [yaml](https://yaml.org/) text input following 
the form:

    version: 1
    rules:
        <ROLE-NAME>:
            - <METHOD-0>
            - <METHOD-1>
            - ...

See below for an exhaustive list of valid methods. You may also use "*" as wildcard.
The default policies can be downloaded [here]({{ .DefaultPoliciesURL }})

