preemptible-address-controller
==============================

[![CircleCI](https://circleci.com/gh/dulltz/preemptible-address-controller.svg?style=svg)](https://circleci.com/gh/dulltz/preemptible-address-controller)
[![Docker Repository on Quay](https://quay.io/repository/tsuruda/preemptible-address-controller/status "Docker Repository on Quay")](https://quay.io/repository/tsuruda/preemptible-address-controller)

If you have run a hobby GKE cluster,
you would think so that - "I want a static IP, and want to use preemptible instances, but I do not want to pay $18/month for ingress controller..."

preemptible-address-controller controls Node resources of preemptible instances, and keeps their IP at your specified static address without Ingress resources. 

How to use
----------

TBD