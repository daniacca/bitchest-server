# [1.1.0](https://github.com/daniacca/bitchest-server/compare/v1.0.0...v1.1.0) (2025-08-04)


### Bug Fixes

* fixed command input tokenizer with a custom algorithm to manage qouted data with space ([c835af1](https://github.com/daniacca/bitchest-server/commit/c835af1ed2207ce660cf80c825d361eb786dd217))
* fixed command input tokenizer with a custom algorithm to manage qouted data with space ([#4](https://github.com/daniacca/bitchest-server/issues/4)) ([6cb35e1](https://github.com/daniacca/bitchest-server/commit/6cb35e1db83b3e1ab612d6d8c0bacb96162d1b70))


### Features

* **commands:** added MEMORY STATS command ([#3](https://github.com/daniacca/bitchest-server/issues/3)) ([9c2d297](https://github.com/daniacca/bitchest-server/commit/9c2d297cf705d10a5a87e0ede0c59b1b20184006))

# 1.0.0 (2025-07-31)


### Bug Fixes

* **docker:** fixed docker build and docker run ([5e2a14f](https://github.com/daniacca/bitchest-server/commit/5e2a14f296635430874f3f57130941b7e1322fc9))
* **set-command:** changed nil response with proper NullBulk protocol return value ([3eb8688](https://github.com/daniacca/bitchest-server/commit/3eb8688cd57be3f66ce6630f112c09aaf7f095ad))


### Features

* **cache-server:** basic memory cache server, with SET-GET-KEYS-PING-DEL-EXISTS_FLUSHALL commands ([b3b1762](https://github.com/daniacca/bitchest-server/commit/b3b17628abfb88be51a05364141ac87835cbd08d))
* **cli:** added server CLI ([bf64f4a](https://github.com/daniacca/bitchest-server/commit/bf64f4a35948c7cfab962bf64b0208e702c4c79b))
* **expire:** adding EXPIRE, TTL commands, changed SET command with optional EX options ([ffe9ffe](https://github.com/daniacca/bitchest-server/commit/ffe9ffeea5dfcbbeceff21b4bb850b2ddb0cffbd))
* **server-conf:** added host and port configuration for server ([9d090c8](https://github.com/daniacca/bitchest-server/commit/9d090c899354f480ad1df7a58fdb3ccf7eb465db))
* **set-options:** added NX and XX optional options to SET key command ([af9c760](https://github.com/daniacca/bitchest-server/commit/af9c76000a4e7054f1629e97adadfb10533fa5b1))
