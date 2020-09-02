你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# Lovi

[![Build Status](https://travis-ci.org/pakerfeldt/lovi.svg?branch=master)](https://travis-ci.org/pakerfeldt/lovi)

Lovi is a lightweight pager application for message distribution. It offers a web api for triggering events and uses a policy configuration to determine how and to whom the events gets distributed.

This application currently supports the following transport mechanisms:

- SMS / Phone calls (through [46elks.se](https://46elks.se))
- stdout (for logging)

## How does it work?

Lovi is configured using a yaml file. Here you activate the transports you want to use and define the set of policies for which you plan to send events. A policy describes how events will be sent and to whom. See [Configuration](https://github.com/pakerfeldt/lovi/wiki/Configuration) wiki page for examples.

You can have multiple policies for different needs, each with its own set of configuration.

An event is triggered by calling `http://[your-ip]:8080/event/trigger/{policy}?message=Your%20message`.

## Running

The recommended way of running lovi is through Docker.
`docker run -p 8080:8080 -i -t pakerfeldt/lovi:1.0.0__linux_amd64`.
Lovi will try to read /config.yaml and listen to port 8080. You may want to change this by setting the `CONFIG` and `PORT` environment variables respectively.

## Contribute

See [Transports](https://github.com/pakerfeldt/lovi/wiki/Transports) wiki page for a guide on how to add new transports.
