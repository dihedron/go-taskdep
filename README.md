# Task Dependency Solver

A dependency solver for deployment-time delayed task execution.

## Motivation

While working on CI/CD pipeline to produce a VM-based immutable infrastructure, I wanted to achieve the following:

1. deploy all software using native package formats, such as DEBs on Ubuntu and RPM on RHEL/CentOS;
2. be able to create an immutable VM image reusable across environments, so the creation of bindings should be deferred to the VM instance creation (```cloud-init```).

These things are quite difficult to attain when  you're dealing with web applications deployed unto an application server (such as JBoss/Wildfly) which is itself a "software package manager" handling the registration and configuration of JDBC drivers, EARs, WARs, datasources, JMS queue definitions, system properties directly.

The assumprtion that a self-contained DEB wrapping an EAR file would copy the EAR onto the filesystem to a temporary directory (e.g. /opt/my-company/my-app), then in the ```postinst``` script simply launch the Wildfly CLI to deploy the application and be done with it does not take into account the fact that a web application expects all dependencies to be satisfied __before__ it is deployed, otherwise the deployment fails.

If there were a way to __deploy an application without starting it__, and the application server did not attempt to perform (and check) the bindings at deployment time, deferring it to start time, this would still be attainable. As far as I know and understand, it is not so.

A web application is almost never self-contained: it depends on systems properties, JDBC drivers, datasource connections, JMS queues an so on, whose configurations are only realy known at runtime. This is in strong contrast to what you want to do with an immutable infrastructure: move as much configuration management stuff to the bake phase, so that the VM instance startup stays as simple as can be.

Thus, we can still use DEBs and RPMs to bring software packages (EARs, WARs, JARs, settings) to the target VM template, but we __must__ defer their installation to the actual VM instance, which knows how to weave them together using environment-specific bindings.

Still, having a custom ```cloud-init``` that orchestrates the deployment on a per application basis means that we cannot build the DEBs and RPMs irrespective of the other components they will be wworking with: the native package is not self-contained and relies on the whole deployment sequence being written in a single ```cloud-init``` script that looks as follows:

```bash
jboss-cli deploy jdbc-driver
jboss-cli add app1-datasource
jboss-cli deploy app1
jboss-cli add app2-datasource
jboss-cli deploy app2 (depending on app1)
...
```

The coupling among all components in the software stack is complete.

That's ugly.

## Solution

The solution I came up with is a custom application that orchestrates the variuos tasks based on the metadata deposited on the VM by the native packages.

Is is based on the following assumptions:

1. each native package comprises:
- some metadata, in YAML or JSON format, sopied into a well known location;
- the binaries (EARs, WARs...), copied into a staging directory (e.g. ```/opt/my-company/app1```);
- one or more "delayed install script(s)", copied somewhere on the filesystem (e.g. the same staging directory where the binaries are);
2. each native package declares in its metadata:
- an ID;
- the IDs of all packages it depends upon;
- the path to the "delayed install script(s)" to be executed to configure the package;
3. the resulting dependency graph is acyclic.

```cloud-init``` won't have to know anything about the complex inter-dependencies of application components and pre-requisites: this application will pick the metadata, build a graph, then execute the delayed install scripts one at a time in due order.
