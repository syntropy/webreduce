# WebReduce

Scriptable resource manipulation tool.

Define behaviour as [Purely functional](http://en.wikipedia.org/wiki/Purely_functional) handler implemented in [Javascript](http://en.wikipedia.org/wiki/JavaScript). Control the flow of entities in, through and out of the system with the help of Destinations & Sources.

## API

### Behaviour

Create a behaviour, its destinations and its sources:

    POST http://lolcathost/behaviour HTTP/1.1
    content-type: text/javascript
    x-wr-destination: http://lolcathost/meow
    x-wr-source: http://lolcathost/ninjas

    function main(data) { return data; }
    ---
    302
    Location: http://lolcathost/abc123

Update the metainformation of a created behaviour:

    PUT http://lolcathost/abc123 HTTP/1.1
    x-wr-destination: http://sudo.com/make/sandwich
    x-wr-source: http://lolcathost/nyan
    ---
    204

Request the metainformation for a behaviour:

    HEAD http://lolcathost/abc123 HTTP/1.1
    ---
    200
    x-wr-destination: http://sudo.com/make/me/sandwich
    x-wr-source: http://lolcathost/nyan

Delete a behaviour:

    DELETE http://lolcathost/abc123 HTTP/1.1
    ---
    204


### Destination & Source

Destinations and sources MUST be a valid absolute [URI](http://en.wikipedia.org/wiki/Uniform_resource_identifier) and can be anything from HTTP url to a database location, but MUST be of a supported protocol. References are passed via the `x-wr-destination` & `x-wr-source` header.

Create multiple destinations or sources with the help of a comma-sperated list:

    x-wr-destination: http://lolcathost/sandwich,http://lolcathost/nyan

### Languages

At the time, only LUA is supported. The design allows adding new interpreters very easily, though.

The supplied code is treated as a function and has to comply to a small set of assumptions:

* The first argument is the data which is to be processed
* The second argument is the current persistent state object (see below) 
* The returned value is the new persistent state object

Every behaviour has the ability so save data which is kept across invocations. The engine is agnostic about what format the persistent data has.

The engine will execute multiple instance of the behaviour in parallel on different datasets to improve performance. By simple collision detection on the state object a behaviour may be re-run.
A machine learning algorithm will reduce (or increase) the number of instances to avoid collisions.

A behaviour can emit (multiple) data during execution which is added to the output queue.

#### Lua

Example script:

```Lua
-- Get parameters
local params = {...}; 
local data = params[1];
local state = params[2];

-- Emit some data
emit(data);
-- Emit some more data
emit(data)

-- Return the new state
return state;
```

## Development

### Setup

### Branching

This repository is organized and maintained with the help of [gitflow](https://github.com/nvie/gitflow). Developers are encouraged to use it when working with this repository.

We use the following naming convention for branches:

* `develop` (during development)
* `master` (will be or has been released)
* `feature/<name>` (feature branches)
* ` ` (empty version prefix)

During development, you should work in feature branches instead of committing to `master` directly. Tell gitflow that you want to start working on a feature and it will do the work for you (like creating a branch prefixed with `feature/`):

    git flow feature start <FEATURE_NAME>

The work in a feature branch should be kept close to the original problem. Tell gitflow that a feature is finished and it will merge it into `master` and push it to the upstream repository:

    git flow feature finish <FEATURE_NAME>

Even before a feature is finished, you might want to make your branch available to other developers. You can do that by publishing it, which will push it to the upstream repository:

    git flow feature publish <FEATURE_NAME>

To track a feature that is located in the upstream repository and not yet present locally, invoke the following command:

    git flow feature track <FEATURE_NAME>

Changes that should go into production should come from the up-to-date master branch. Enter the "release to production" phase by running:

    git flow release start <VERSION_NUMBER>

In this phase, only meta information should be touched, like bumping the version and update the history. Finish the release phase with:

    git flow release finish <VERSION_NUMBER>

### Versioning

This project is versioned with the help of the [Semantic Versioning Specification](http://semver.org/) using `0.0.0` as the initial version. Please make sure you have read the guidelines before increasing a version number either for a release or a hotfix.

## Contributing

**normal**

1. [Fork](http://help.github.com/forking/) WebReduce
2. Create a topic branch `git checkout -b my_branch`
3. Push to your branch `git push origin my_branch`
4. Create a [Pull Request](http://help.github.com/pull-requests/) from your branch

**git-flow**

1. [Fork](http://help.github.com/forking/) WebReduce
2. Create a feature `git flow feature start my-feature`
3. Publish your featue `git flow feature publish my-feature`
4. Create a [Pull Request](http://help.github.com/pull-requests/) from your branch

## Meta

* Code: `git clone git@github.com:syntropy/webreduce.git`
* Home: https://github.com/syntropy/webreduce
* Docs: https://github.com/syntropy/webreduce
* Bugs: https://github.com/syntropy/webreduce/issues

This project uses [Semantic Versioning](http://semver.org) & [git flow](http://nvie.com/posts/a-successful-git-branching-model/) with the help of [git-flow](https://github.com/nvie/gitflow)
