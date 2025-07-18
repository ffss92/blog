---
title: "From mocked to real: Postgres integration tests done right in Go."
subtitle: |
  Mocks can only take you so far — discover how real Postgres integration tests in Go help catch real bugs, reflect true system behavior, and make your code production-ready.
author: "@ffss"
draft: true
date: "2025-07-18"
tags:
  - Go
  - PostgreSQL
  - Docker
---

## Introduction

In this article we will explore how to use the excellent
[ory/dockertest](https://github.com/ory/dockertest) package to
properly run Postgres tests using Docker in Go.
