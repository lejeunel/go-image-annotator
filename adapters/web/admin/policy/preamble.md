*Policies specify which method are allowed for a given role.*

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
For reference, the default policies can be downloaded [here]({{ .DefaultPoliciesURL }})
